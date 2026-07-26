#pragma bits 32

#include "stdlib.h"
#include "shell.h"
#include "util.h"
#include "lufs.h"
#include "images.h"

#ifndef __LCC__
    #error "LunaOS must be compiled with LCC (other compilers are not supported.)"
#endif

long int* rand = 0x90909090;
int recursed = 0;

void myfunc() {
    long int a = *rand;
    tohex(a, 1);
    if (recursed == 0) {
        recursed = 1;
        myfunc();
    } else {
        return;
    }
    tohex(a, 1);
}

void _cstart() __attribute__((noreturn)) {
    asm (".byte 0x22");
    myfunc();
    if (fopen("NOTEPAD     SYS", 0)->Address == 0x00000000) {
        fcreate("NOTEPAD     SYS", 256);
    }
 
    puts32("Welcome to ", COLOR_WHITE, COLOR_BLACK);
    puts32("Luna", COLOR_LCYAN, COLOR_BLACK);
    puts32("OS!\n", COLOR_WHITE, COLOR_BLACK);

    while (1) {
        shell();
    }
}

