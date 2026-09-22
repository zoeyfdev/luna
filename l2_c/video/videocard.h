#pragma once

#include <stdint.h>

extern unsigned char (*v_read_video_memory)(uint32_t);
extern void (*v_write_video_memory)(uint32_t, unsigned char);
extern void (*v_print_char)(unsigned char, unsigned char, unsigned char);
extern void (*v_set_cursor)(int, int);
extern void (*v_get_cursor)(int*, int*);
