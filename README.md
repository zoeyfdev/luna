# Luna<br>
A simple, lightweight RISC CPU architecture.<br><br>

# Requirements<br>
- Windows, or Unix-like system (includes MacOS, Linux, FreeBSD, etc.)<br>
- SDL2<br>
- GCC (GNU Compiler Collection)<br>
- GNU Make<br>

# Installation (Unix-like)<br>
- Clone the repository using `git clone`<br>
- Run `make; sudo make install` to install the Luna L2 emulator and toolchain<br>

# Installation (Windows)<br>
- Download and install MSYS2 [here](https://www.msys2.org/)<br>
- Open the MSYS2 MINGW64 shell **as Administrator** and run `pacman -S base-devel git mingw-w64-x86_64-gcc mingw-w64-x86_64-SDL2 mingw-w64-x86_64-pkg-config`<br>
- Clone the repository using `git clone`<br>
- Run `make BUILD_MODE=windows; make install`<br>

# Running a disk image<br>
- Run `luna-l2 <disk image>` to execute a disk image.<br>
- Note: there are 3 hardware slots for a disk image:<br>
HDD: Hard disk drive; use `-hdd <file>` or just `<file>` to insert a file into the slot.<br>
SD: Secure Digital card/USB; use `-sd <file>` to insert a file into the slot.<br>
DVD: Optical disc slot; use `-dvd <file>` to insert a file into the slot.<br>
Note: You can use F1/F2/F3 to eject/insert media into these slots.<br><br>

# Emulator parameters<br>
- `-ram <RAM amount in bytes>` - adjust the amount of RAM available to the emulator<br>
- `-gpu <gpu device>` - change the graphics processor device (defaults to `g1x`); see the list below for valid equipment.<br>
- `-apu <apu device>` - change the audio processor devices (defaults to `s1`); see the list below for valid equipment.<br><br>

# Equipment list<br>
## Graphics<br>
### Luna G1<br>
- Emulator equipment ID: `g1`<br>
- Valid resolutions: 320x200@8bpp<br>
- Video memory: 64 KiB (65,536 bytes)<br>
- Extra features: None<br><br>
### Luna G1X<br>
- Emulator equipment ID: `g1x`<br>
- Valid resolutions: 320x200@8bpp<br>
- Video memory: 64 KiB (65,536 bytes)<br>
- Extra features: Simple transparency via color code 0xE3<br><br>
## Audio<br>
### Luna S1<br>
- Emulator equipment ID: `s1`<br>
- Extra features: None<br><br>
