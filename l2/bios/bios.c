#include <stdint.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#include "../hardware/video_common.h"
#include "../cpu/registers.h"
#include "../cpu/memory.h"
#include "../cpu/memdefs.h"
#include "disk.h"

char* HDD_FILE;
char* SD_FILE;
char* DVD_FILE;
int bios_boot_drive;
unsigned char* BIOS_RAM = NULL;

#define BIOS_RAM_SIZE 192

void bios_write_string(char* s) {
    while (*s) {
        v_print_char(*s, 255, 0);
        s++;
    }
}

void bios_write_line(char* s) {
    bios_write_string(s);
    v_print_char(0x0a, 255, 0);
}

void bios_splash() {
    bios_write_line("Luna L2");
    bios_write_line("BIOS: Integrated BIOS");
    bios_write_line("Copyright (c) 2026 zoeyfdev\n");
}

void bios_init() {
    BIOS_RAM = malloc(BIOS_RAM_SIZE); // 6 bytes per interrupt * 32 interrupts
}

unsigned char bios_read_memory(uint32_t address) {
    if (address < BIOS_RAM_SIZE)
        return BIOS_RAM[address];
    return (unsigned char) rand() & 0xFF;
}

void bios_write_memory(uint32_t address, unsigned char value) {
    if (address < BIOS_RAM_SIZE)
        BIOS_RAM[address] = value;
}

void bios_handle_interrupt(uint32_t code) {
    uint32_t handler_addr = (code - 1) * 6;

    // 0: BIOS
    // 1: software
    // else: disabled
    if (BIOS_RAM[handler_addr] == 1) {
        uint32_t pc_addr = (BIOS_RAM[handler_addr + 1] << 24)
            | (BIOS_RAM[handler_addr + 2] << 16)
            | (BIOS_RAM[handler_addr + 3] << 8)
            | (BIOS_RAM[handler_addr + 4]);
        set_register(PC, pc_addr);
        return;
    }

    switch (code) {
    case 0x01:
        v_print_char((unsigned int) get_register(R1) & 0xFF, (unsigned int) get_register(R2) & 0xFF, (unsigned int) get_register(R3) & 0xFF);
        break;
    case 0x02:
        // PIT reserved
        break;
    case 0x03: {
            char* file;
            switch (get_register(R1)) {
            case 0:
                file = HDD_FILE;
                break;
            case 1:
                file = SD_FILE;
                break;
            case 2:
                file = DVD_FILE;
                break;
            }

            if (file != NULL)
                set_register(R1, 1);
            else
                set_register(R1, 0);
            break;
        }
    case 0x04:
        // Syscall reserved
        break;
    case 0x05:
        // Power reserved
        break;
    case 0x06:
        // Keyboard interrupt
        break;
    case 0x07:
        // Illegal instruction trap
        break;
    case 0x08:
        // Unmapped
        break;
    case 0x09:
        // Unmapped
        break;
    case 0x0A:
        // Memory query
        break;
    case 0x0B:
        // Load sector from disk
        load_sector(get_register(R2), get_register(R1), get_register(R3));
        break;
    case 0x0C:
        v_set_cursor(get_register(R1), get_register(R2));
        break;
    case 0x0D:
        write_sector(get_register(R2), get_register(R1), get_register(R3));
        break;
    case 0x0E: {
            int x;
            int y;
            v_get_cursor(&x, &y);

            set_register(R1, x);
            set_register(R2, y);
            break;
        }
    case 0x0F: {
            // Reboot
            memset(MEMORY, 0x00, MEMSIZE);
            memset(BIOS_RAM, 0x00, BIOS_RAM_SIZE);
            for (int i = 0; i < 35; i++)
                set_register(i, 0);
            v_gpu_reset();
            break;
        }
    case 0x10:
        set_register(R1, bios_boot_drive);
        break;
    case 0x11:
        exit(0);
        // Shut down machine
        break;
    }
}
