#pragma bits 32

#include "stdbool.h"

asm (".global play_sound_loc");
asm ("play_sound_loc:");

void play_sound(void* buffer, long int size, bool block) {
    short short int* done_flag = (char*) 0x80000009;
    *done_flag = 0;

    *(long int*) 0x80000001 = buffer;
    *(long int*) 0x80000005 = size;
    *(char*) 0x80000000 = 1; 

    if (block == true) {
        while (*done_flag == 0) {}
    }

    return;
}
