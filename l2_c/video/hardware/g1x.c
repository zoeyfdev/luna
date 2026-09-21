#include <stdlib.h>
#include <stdint.h>

#include "../nrgba/nrgba.h"
#include "font8x8/font.h"
#include "../../util/util.c"

#define SCREEN_WIDTH  320
#define SCREEN_HEIGHT 200
#define VRAM          0xFFFF

unsigned char* VIDEO_MEMORY = NULL;
nrgba_pixel Palette[256];

int cursor_x = 0;
int cursor_y = 0;

void initialize_component() {
    VIDEO_MEMORY = malloc(VRAM);
    for (int i = 0; i < 256; i++) {
        // TODO: make real BIOS color palette
        int r = ((i >> 5) & 0x07) * 255 / 7;
        int g = ((i >> 2) & 0x07) * 255 / 7;
        int b = (i & 0x03) * 255 / 3;

        Palette[i].r = r;
        Palette[i].g = g;
        Palette[i].b = b;
        Palette[i].a = 255;
    }
}

void write_video_memory(uint32_t addr, unsigned char content) {
    if (addr >= VRAM)
        VIDEO_MEMORY[addr] = content;
}

unsigned char read_video_memory(uint32_t addr) {
    if (addr >= VRAM)
        return VIDEO_MEMORY[addr];
    return 0x00000000;
}

void set_cursor(int x, int y) {
    cursor_x = x;
    cursor_y = y;
}

void get_cursor(int* x, int* y) {
    *x = cursor_x;
    *y = cursor_y;
}

nrgba_image* return_framebuffer() {
    nrgba_image* img = malloc((SCREEN_WIDTH * SCREEN_HEIGHT * sizeof(nrgba_pixel)) + sizeof(nrgba_image));
    nrgba_pixel* pix_start = (nrgba_pixel*) ((unsigned char*) img + sizeof(nrgba_image));
    
    int i = 0;
    for (int y = 0; y < 200; y++) {
        for (int x = 0; x < 320; x++) {
            pix_start[i] = Palette[VIDEO_MEMORY[i]];
            i++;
        }
    }

    img->img = (unsigned char*) pix_start;
    img->stride = img->width * 4;

    return img;
}

void clear_video_memory() {
    for (int i = 0; i < SCREEN_WIDTH * SCREEN_HEIGHT; i++)
        VIDEO_MEMORY[i] = 0;
}

void scroll() {
    int line_size = SCREEN_WIDTH * 8;
    int visible_lines = SCREEN_HEIGHT / 8;

    copy(VIDEO_MEMORY, 0, VRAM, VIDEO_MEMORY, line_size, VRAM);

    int bottom_start = (visible_lines - 1) * line_size;
    for (int i = bottom_start; i < VRAM; i++) {
        VIDEO_MEMORY[i] = 0;
    }

    cursor_y = visible_lines - 1;
}

void print_char(unsigned char c, unsigned char fg, unsigned char bg) {
    switch (c) {
    case 0x0a:
        cursor_y++;
    case 0x0d:
        cursor_x = 0;
    case 0x00:
        return;
    }
    
    if (cursor_y >= 200 / 8)
        scroll();

    if (c == 0x0a)
        return;

    int x = cursor_x * 8;
    int y = cursor_y * 8;

    int idx = (int) c;
    char* glyph = font8x8[0];

    if (idx >= 0 && idx < sizeof(font8x8))
        glyph = font8x8[idx];

    for (int row = 0; row < 8; row++) {
        char line = glyph[row];
        for (int col = 0; col < 8; col++) {
            unsigned char mask = 1 << col;
            unsigned char color;

            if (line & mask != 0)
                color = fg;
            else
                color = bg;

            int px = (y + row) * 320 + (x + col);
            VIDEO_MEMORY[px] = color;
        }
    }
}
