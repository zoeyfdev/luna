package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"time"
	"luna_l2/shared"
	"math/rand"
	"fmt"
	"os"
)

var MemoryAudio [10]byte

// MEMORY layout:
	// Byte 1: Play flag
	// Byte 2, 3, 4, 5: 32-bit pointer of audio
	// Byte 6, 7, 8, 9: 32-bit size of audio
	// Byte 10: done flag

func AudioController() {
	Device, err := sdl.OpenAudioDevice("", false, &sdl.AudioSpec {
		Freq: 48000,
		Format: sdl.AUDIO_S8,
		Channels: 1,
		Samples: 4096, // Just in case of long audios
	}, nil, 0)
	if err != nil {
		fmt.Println("luna-l2: failed opening audio:", err)
		os.Exit(1)
	}

	for {
		if MemoryAudio[0] != 0 {
			MemoryAudio[0] = 0	

			Length := uint32(MemoryAudio[5]) << 24 | uint32(MemoryAudio[6]) << 16 | uint32(MemoryAudio[7]) << 8 | uint32(MemoryAudio[8])
			Cursor := uint32(MemoryAudio[1]) << 24 | uint32(MemoryAudio[2]) << 16 | uint32(MemoryAudio[3]) << 8 | uint32(MemoryAudio[4])
			
			Buffer := make([]byte, Length)
			for i := uint32(0); i < Length; i++ {
				Buffer[i] = shared.Mapper(Cursor)
				Cursor++
			}

			sdl.PauseAudioDevice(Device, false)
			sdl.QueueAudio(Device, Buffer)

			for sdl.GetQueuedAudioSize(Device) > 0 {
				time.Sleep(15 * time.Millisecond)
			}

			MemoryAudio[9] = 1
		}
		time.Sleep(time.Duration(15) * time.Millisecond)
	}
}

func WriteAudioMemory(addr uint32, content byte) {
	if addr >= uint32(len(MemoryAudio)) { return }
	MemoryAudio[addr] = content
}

func ReadAudioMemory(addr uint32) byte {
	if addr >= uint32(len(MemoryAudio)) { return byte(rand.Intn(0xFF)) }
	return MemoryAudio[addr]
}
