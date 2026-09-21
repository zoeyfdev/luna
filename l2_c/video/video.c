#include <stdio.h>
#include <stdlib.h>
#include <SDL2/SDL.h>
#include <stdint.h>

#include "../component/component.h"
#include "nrgba/nrgba.h"

#define SCREEN_WIDTH    960
#define SCREEN_HEIGHT   600

#ifdef SDL_HINT_RENDER_SCALE_QUALITY
    #undef SDL_HINT_RENDER_SCALE_QUALITY
#endif

#define SDL_HINT_RENDER_SCALE_QUALITY "0"

unsigned char (*v_read_video_memory)(uint32_t);
void (*v_write_video_memory)(uint32_t, unsigned char);
void (*v_print_char)(unsigned char, unsigned char, unsigned char);
bool VIDEO_READY = false;

void reset_aspect_ratio(SDL_Renderer* renderer) {
    int wo;
    int ho;

    SDL_GetRendererOutputSize(renderer, &wo, &ho);

    float aspect = (float) SCREEN_WIDTH / (float) SCREEN_HEIGHT;
    float actual = (float) wo / (float) ho;

    int w;
    int h;
    int x;
    int y;

    if (actual > aspect) {
        w = (int) (((float) ho) * aspect);
        h = (int) ho;
        x = (int) (wo - w) / 2;
        y = 0;
    } else {
        h = (int) (((float) wo) / aspect);
        w = (int) wo;
        y = (int) (ho - h) / 2;
        x = 0;
    }

    SDL_Rect rect;
    rect.x = (int32_t) x;
    rect.y = (int32_t) y;
    rect.w = (int32_t) w;
    rect.h = (int32_t) h;

    SDL_RenderSetViewport(renderer, &rect);
}

int initialize_window() {
    // for now
    void* video_component = initialize_component("/usr/local/lib/l2/video/g1x.so"); 
    
    if (SDL_Init(SDL_INIT_EVERYTHING) < 0) {
        printf("luna-l2: could not initialize window: %s\n", SDL_GetError());
        exit(1);
    }

    #if defined linux && SDL_VERSION_ATLEAST(2, 0, 8)
        if (!SDL_SetHint(SDL_HINT_VIDEO_X11_NET_WM_BYPASS_COMPOSITOR, "0")) {
            printf("luna-l2: could not disable compositor bypass: %s\n", SDL_GetError());
            exit(1);
        }
    #endif


    SDL_Window* window = SDL_CreateWindow("Luna L2", SDL_WINDOWPOS_UNDEFINED, SDL_WINDOWPOS_UNDEFINED, SCREEN_WIDTH, SCREEN_HEIGHT, SDL_WINDOW_SHOWN | SDL_WINDOW_RESIZABLE);
    if (!window) {
        printf("luna-l2: could not initialize window: %s\n", SDL_GetError());
        exit(1);
    }

    SDL_Renderer* renderer = SDL_CreateRenderer(window, -1, SDL_RENDERER_ACCELERATED);
    if (!renderer) {
        printf("luna-l2: could not initialize renderer: %s", SDL_GetError());
        exit(1);
    }


    // Define function pointers
    nrgba_image* (*return_framebuffer)();
    void (*initialize_component)();

    return_framebuffer = return_component_function(video_component, "return_framebuffer");
    initialize_component = return_component_function(video_component, "initialize_component");

    v_read_video_memory = return_component_function(video_component, "read_video_memory");
    v_write_video_memory = return_component_function(video_component, "write_video_memory");
    v_print_char = return_component_function(video_component, "print_char");

    (*initialize_component)();

    nrgba_image* img = (*return_framebuffer)();
    SDL_Texture* texture = SDL_CreateTexture(renderer, SDL_PIXELFORMAT_ABGR8888, SDL_TEXTUREACCESS_STREAMING, img->width, img->height); 

    SDL_UpdateTexture(texture, NULL, img->img, img->stride);

    bool quit = false;

    // Event loop
    while(!quit)
    {
        SDL_Event e;    
        while (SDL_PollEvent(&e)) {
            switch (e.type) {
            case SDL_QUIT:
                quit = true;
                exit(0);
            }
        }
 
        img = (*return_framebuffer)();
        SDL_UpdateTexture(texture, NULL, img->img, img->stride);
        SDL_RenderClear(renderer);
        SDL_RenderCopy(renderer, texture, NULL, NULL);
        SDL_RenderPresent(renderer);

        reset_aspect_ratio(renderer);

        SDL_Delay(10);

        VIDEO_READY = true;
    }

    SDL_DestroyWindow(window);
    SDL_Quit();

    return 0;
}
