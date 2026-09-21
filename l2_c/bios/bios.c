#include "../video/videocard.h"

char* HDD_FILE;
char* SD_FILE;
char* DVD_FILE;

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
