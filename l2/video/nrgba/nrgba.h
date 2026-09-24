#pragma once

typedef struct {
    unsigned char r;
    unsigned char g;
    unsigned char b;
    unsigned char a;
} nrgba_pixel;

typedef struct {
    uint32_t width;
    uint32_t height;
    uint32_t stride;
    unsigned char* img;
} nrgba_image;
