#pragma bits 32

#include "stdlib.h"
#include "shell.h"
#include "util.h"
#include "lufs.h"
#include "images.h"
#include "stdbool.h"

void _cstart() __attribute__((noreturn)) {
    if (fopen((char*) fntf("NOTEPAD.SYS"), false)->Address == NULL)
        fcreate((char*) fntf("NOTEPAD.SYS"), 256);

    puts32("Welcome to ", COLOR_WHITE, COLOR_BLACK);
    puts32("Luna", COLOR_LCYAN, COLOR_BLACK);
    puts32("OS!\n", COLOR_WHITE, COLOR_BLACK);

    while (1)
        shell();
}

