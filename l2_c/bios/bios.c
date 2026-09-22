#include <stdint.h>
#include <stdlib.h>

#include "../video/videocard.h"
#include "../cpu/registers.h"
#include "disk.h"

char* HDD_FILE;
char* SD_FILE;
char* DVD_FILE;
int bios_boot_drive;

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

void bios_handle_interrupt(uint32_t code) {
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
        // Keyboard reserved
        break;
    case 0x06:
        // Power interrupt
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
        load_sector(get_register(R1), get_register(R2), get_register(R3));
        break;
    case 0x0C:
        v_set_cursor(get_register(R1), get_register(R2));
        break;
    case 0x0D:
        write_sector(get_register(R1), get_register(R2), get_register(R3));
        break;
    case 0x0E: {
            int x;
            int y;
            v_get_cursor(&x, &y);

            set_register(R1, x);
            set_register(R2, y);
            break;
        }
    case 0x0F:
        // Reboot
        break;
    case 0x10:
        set_register(R1, bios_boot_drive);
        break;
    case 0x11:
        // Shut down machine
        break;
    }
}
