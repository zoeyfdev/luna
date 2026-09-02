#pragma bits 32

#include "stdlib.h"
#include "util.h"

void textedit_init(char* buffer) {
    save_graphics_buf();
    video_save_cursor();
    render_buf((void*) 0x40404040);
    video_set_cursor(0, 0);
    
    puts32("TextEdit v0.1 - LunaOS\n\n", COLOR_WHITE, COLOR_BLACK);
    readin((char*) buffer, 0, 0);

    render_buf(0x30303030);
    video_load_cursor();
}
