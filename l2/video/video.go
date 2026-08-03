package video

import (
	"fmt"
	"image"
	"luna_l2/component"
	"luna_l2/keyboard"
	"luna_l2/proxy"
	"luna_l2/shared"
	"os"
	"runtime"
	"unsafe"

	"github.com/ncruces/zenity"
	"github.com/veandco/go-sdl2/sdl"
)

var CommonComponentPathPrefix string = "/usr/local/lib/l2/video/"
var WindowsComponentPathPrefix string = "C:\\Program Files (x86)\\Luna L2\\lib\\l2\\video\\"
var VideoComponent component.Component

// Function definitions
var ReturnFramebuffer func() *image.RGBA

// Frontend code
var Ready bool
var FS bool
var Grab bool

func ResetAspectRatio(renderer *sdl.Renderer) {
	WO, HO, _ := renderer.GetOutputSize()
	aspect := float32(960) / float32(600)
	actual := float32(WO) / float32(HO)

	var H int
	var W int
	var X int
	var Y int

	if actual > aspect {
		H = int(HO)
		W = int(float32(HO) * aspect)
		X = (int(WO) - W) / 2
		Y = 0
	} else {
		W = int(WO)
		H = int(float32(WO) / aspect)
		X = 0
		Y = (int(HO) - H) / 2
	}

	renderer.SetViewport(&sdl.Rect{X: int32(X), Y: int32(Y), W: int32(W), H: int32(H)})
}

func FileOpenDialogue(title string, drive int) {
    ZOpen := func(title string) {
        _path, err := zenity.SelectFile(
            zenity.Title(title),
        )
        if err != nil {
            return
        }
        switch drive {
        case 0:
            shared.Filename = _path
        case 1:
            shared.SDFilename = _path
        case 2:
            shared.OpticalFilename = _path
        }
    }

    switch runtime.GOOS {
    case "darwin":
        ZOpen(title)
    default:
        go ZOpen(title)
    }
}

func InitializeWindow(ComponentName string) {
	prefix := CommonComponentPathPrefix
	ext := ".so"
	if runtime.GOOS == "windows" {
		prefix = WindowsComponentPathPrefix
		ext = ".dll"
	}
	VideoComponent = component.InitializeComponent(prefix + ComponentName + ext)
	component.ReturnComponentFunction(VideoComponent, "InitializePalette").(func())()

	ReturnFramebuffer = component.ReturnComponentFunction(VideoComponent, "ReturnFramebuffer").(func() *image.RGBA)
	proxy.VideoPrintChar = component.ReturnComponentFunction(VideoComponent, "PrintChar").(func(rune, byte, byte))
	proxy.VideoSetCursor = component.ReturnComponentFunction(VideoComponent, "SetCursor").(func(int, int))
	proxy.VideoGetCursor = component.ReturnComponentFunction(VideoComponent, "GetCursor").(func() (int, int))
	proxy.VideoClearVideoMemory = component.ReturnComponentFunction(VideoComponent, "ClearVideoMemory").(func())
	proxy.VideoReadVideoMemory = component.ReturnComponentFunction(VideoComponent, "ReadVideoMemory").(func(uint32) byte)
	proxy.VideoWriteVideoMemory = component.ReturnComponentFunction(VideoComponent, "WriteVideoMemory").(func(uint32, byte))

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println("luna-l2: failed to initialize SDL:", err)
		os.Exit(1)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Luna L2", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 960, 600, sdl.WINDOW_SHOWN | sdl.WINDOW_RESIZABLE)
	if err != nil {
		fmt.Println("luna-l2: failed to create window:", err)
		os.Exit(1)
	}
	defer window.Destroy()

	renderer, _ := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	defer renderer.Destroy()

	Frame := ReturnFramebuffer()
	Texture, _ := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(Frame.Bounds().Dx()), int32(Frame.Bounds().Dy()))
	defer Texture.Destroy()

	Texture.Update(nil, unsafe.Pointer(&Frame.Pix[0]), Frame.Stride)

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
				os.Exit(0)
			case *sdl.MouseButtonEvent:
				if t.State == sdl.PRESSED {
					if Grab == false {
						sdl.SetRelativeMouseMode(true)
						window.SetTitle("Luna L2 - Press Ctrl+Alt+G to release grab")
						Grab = true
						continue
					}
				}
			case *sdl.MouseMotionEvent:
				if Grab == false { continue }
				
				t.X /= 3 // not good because hardcoded scale factor
				t.Y /= 3

				if t.X < 0 { t.X = 0} else if t.X > 320 { t.X = 320 }
				if t.Y < 0 { t.Y = 0} else if t.Y > 200 { t.Y = 200 }

				ixh := int(t.X) >> 8
				ixl := int(t.X) & 0xFF

				iyh := int(t.Y) >> 8
				iyl := int(t.Y) & 0xFF

				keyboard.MemoryMouse[2] = byte(ixh)
				keyboard.MemoryMouse[3] = byte(ixl)
				keyboard.MemoryMouse[6] = byte(iyh)
				keyboard.MemoryMouse[7] = byte(iyl)
				
				shared.RaiseInterrupt(0x12)
			case *sdl.KeyboardEvent:
				Shift := t.Keysym.Mod & sdl.KMOD_SHIFT != 0
				Alt := t.Keysym.Mod & sdl.KMOD_ALT != 0
				Control := t.Keysym.Mod & sdl.KMOD_CTRL != 0

				KeyCode := t.Keysym.Sym	
				if t.State == sdl.PRESSED || t.Repeat > 0 {
					switch KeyCode {
					case sdl.K_F1:
						// Insert into HDD slot
						if shared.Filename == "" {
							FileOpenDialogue("Select hard disk file", 0)
						} else {
							shared.Filename = ""
						}
					case sdl.K_F2:
						// Insert into SD slot
						if shared.SDFilename == "" {
							FileOpenDialogue("Select SD/USB file", 1)
						} else {
							shared.SDFilename = ""
						}	
					case sdl.K_F3:
						// Insert into CD/DVD slot
						if shared.OpticalFilename == "" {
							FileOpenDialogue("Select CD/DVD file", 2)
						} else {
							shared.OpticalFilename = ""
						}	
					case sdl.K_F4:
						if shared.Debug == true {
							shared.LogOn = false
							shared.Debug = false
						} else {
							shared.LogOn = true
							shared.Debug = true
						}
					case sdl.K_F5:
						shared.RaiseInterrupt(0xF)
					case sdl.K_F6:
						f, _ := os.Create("memory_dump.bin")
						f.Write((*shared.Memory)[:])
						f.Close()
					case sdl.K_F11:
						if FS == false {
							FS = true
							window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
						} else {
							FS = false
							window.SetFullscreen(0)
						}
					default:
						if KeyCode >= 10000 { continue }
						char := ""
						if Shift == true {
							char = keyboard.Upper(string(KeyCode))
						} else {
							char = keyboard.Lower(string(KeyCode))
						}
	
						if len(char) > 0 {
							if Control == true && Alt == true && (char[0] == 'g' || char[0] == 'G') {
								sdl.SetRelativeMouseMode(false)
								window.SetTitle("Luna L2")
								Grab = false
							} else {
								keyboard.MemoryKeyboard[0] = byte(char[0])
								shared.RaiseInterrupt(0x5)
								shared.SetRegister(0x001b, uint32(char[0]))
							}
						}
					}	
				}
			}
        }

		sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0")
		renderer.Clear()
		Frame := ReturnFramebuffer()
		Texture, _ := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(Frame.Bounds().Dx()), int32(Frame.Bounds().Dy()))
		Texture.Update(nil, unsafe.Pointer(&Frame.Pix[0]), Frame.Stride)
		renderer.Copy(Texture, nil, nil)
		renderer.Present()
		Texture.Destroy()

		ResetAspectRatio(renderer)	

		sdl.Delay(10)
		Ready = true
	}	
}
