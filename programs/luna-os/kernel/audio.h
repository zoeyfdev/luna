#include "stdbool.h"

extern void play_sound(void* buffer, long int size, bool block);
extern void* CRASH_SOUND;
extern void* BOOT_SOUND;
extern void* play_sound_loc;
