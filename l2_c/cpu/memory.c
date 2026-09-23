#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>

#include "registers.h"
#include "memdefs.h"
#include "../video/videocard.h"
#include "../audio/audiocard.h"
#include "../io/keyboard.h"
#include "../bios/bios.h"

// For now we will allocate the entire 1.75 GB of RAM
unsigned char* MEMORY = NULL;

void initialize_memory() {
    MEMORY = (unsigned char*) malloc(MEMSIZE);
    memset(MEMORY, 0x00, MEMSIZE);
}

unsigned char get_memory_main(uint32_t address) {
    if (address >= MEMSIZE)
        return (unsigned char) rand() & 0xFF;
    return MEMORY[address];
}

void set_memory_main(uint32_t address, unsigned char value) {
    if (address >= MEMSIZE)
        return;
    MEMORY[address] = value;
}

unsigned char get_memory(uint32_t address) {
    if (!IS_XEN) {
        // 16-bit memory map
        if (address >= S_SYS_RAM_16 && address <= E_SYS_RAM_16)
            return get_memory_main(address - S_SYS_RAM_16);
        else if (address >= S_VIDEO_RAM_16 && address <= E_VIDEO_RAM_16)
            return v_read_video_memory((address - S_VIDEO_RAM_16) + (get_register(B) * 256));
        else if (address >= S_KEYBOARD_RAM_16 && address <= E_KEYBOARD_RAM_16)
            return keyboard_read_memory(address - S_KEYBOARD_RAM_16);
        else if (address >= S_BIOS_RAM_16 && address <= E_BIOS_RAM_16)
            return bios_read_memory(address - S_BIOS_RAM_16);
        else if (address >= S_AUDIO_RAM_16 && address <= E_AUDIO_RAM_16)
            return audio_read_memory(address - S_AUDIO_RAM_16);
    } else {
        if (address >= S_SYS_RAM_32 && address <= E_SYS_RAM_32)
            return get_memory_main(address - S_SYS_RAM_32);
        else if (address >= S_VIDEO_RAM_32 && address <= E_VIDEO_RAM_32)
            return v_read_video_memory(address - S_VIDEO_RAM_32);
        else if (address >= S_KEYBOARD_RAM_32 && address <= E_KEYBOARD_RAM_32)
            return keyboard_read_memory(address - S_KEYBOARD_RAM_32);
        else if (address >= S_BIOS_RAM_32 && address <= E_BIOS_RAM_32)
            return bios_read_memory(address - S_BIOS_RAM_32);
        else if (address >= S_AUDIO_RAM_32 && address <= E_AUDIO_RAM_32)
            return audio_read_memory(address - S_AUDIO_RAM_32);
    }
    return (unsigned char) rand() & 0xFF;
}

void set_memory(uint32_t address, unsigned char value) {
    if (!IS_XEN) {
        // 16-bit memory map
        if (address >= S_SYS_RAM_16 && address <= E_SYS_RAM_16)
            set_memory_main(address - S_SYS_RAM_16, value);
        else if (address >= S_VIDEO_RAM_16 && address <= E_VIDEO_RAM_16)
            v_write_video_memory(address - S_VIDEO_RAM_16 + (get_register(B) * 256), value);
        else if (address >= S_KEYBOARD_RAM_16 && address <= E_KEYBOARD_RAM_16)
            keyboard_write_memory(address - S_KEYBOARD_RAM_16, value);
        else if (address >= S_BIOS_RAM_16 && address <= E_BIOS_RAM_16)
            bios_write_memory(address - S_BIOS_RAM_16, value);
        else if (address >= S_AUDIO_RAM_16 && address <= E_AUDIO_RAM_16)
            audio_write_memory(address - S_AUDIO_RAM_16, value);
    } else {
        if (address >= S_SYS_RAM_32 && address <= E_SYS_RAM_32)
            set_memory_main(address - S_SYS_RAM_32, value);
        else if (address >= S_VIDEO_RAM_32 && address <= E_VIDEO_RAM_32)
            v_write_video_memory(address - S_VIDEO_RAM_32, value);
        else if (address >= S_KEYBOARD_RAM_32 && address <= E_KEYBOARD_RAM_32)
            keyboard_write_memory(address - S_KEYBOARD_RAM_32, value);
        else if (address >= S_BIOS_RAM_32 && address <= E_BIOS_RAM_32)
            bios_write_memory(address - S_BIOS_RAM_32, value);
        else if (address >= S_AUDIO_RAM_32 && address <= E_AUDIO_RAM_32)
            audio_write_memory(address - S_AUDIO_RAM_32, value);
    }
}


