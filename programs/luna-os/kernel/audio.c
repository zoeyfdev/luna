#pragma bits 32

// #include "stdbool.h"
#include "stdlib.h"

asm (".global play_sound_loc");
asm ("play_sound_loc:");

#define true 1

void play_sound(void* buffer, long int size, char block) {
    puts32("NEW!", COLOR_WHITE, COLOR_BLACK);
    char* done_flag = (char*) 0x80000009;
    *done_flag = 0;

    *(long int*) 0x80000001 = buffer;
    *(long int*) 0x80000005 = size;
    *(char*) 0x80000000 = 1; 

    if (block == true) {
        while (*done_flag == 0) {}
    }

    return;
}
