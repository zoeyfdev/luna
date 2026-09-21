#include <stdio.h>
#include <stdlib.h>
#include <SDL2/SDL.h>

#include "../component/component.h"
#include "nrgba/nrgba.h"

// Define MAX and MIN macros
#define MAX(X, Y) (((X) > (Y)) ? (X) : (Y))
#define MIN(X, Y) (((X) < (Y)) ? (X) : (Y))

// Define screen dimensions
#define SCREEN_WIDTH    800
#define SCREEN_HEIGHT   600
#define SDL_HINT_RENDER_SCALE_QUALITY       "0"

int initialize_window() {
    // for now
    printf("ready to init");
    void* video_component = initialize_component("/usr/local/lib/l2/video/g1x.so");
    
    
    if (SDL_Init(SDL_INIT_EVERYTHING) < 0) {
        printf("luna-l2: could not initialize window: %s\n", SDL_GetError());
        exit(1);
    }

    #if defined linux && SDL_VERSION_ATLEAST(2, 0, 8)
        if (!SDL_SetHint(SDL_HINT_VIDEO_X11_NET_WM_BYPASS_COMPOSITOR, "0")) {
            printf("luna-l2: could not disable compositor bypass: %s\n");
            exit(1);
        }
    #endif

    printf("init window");

    SDL_Window* window = SDL_CreateWindow("Luna L2", SDL_WINDOWPOS_UNDEFINED, SDL_WINDOWPOS_UNDEFINED, SCREEN_WIDTH, SCREEN_HEIGHT, SDL_WINDOW_SHOWN);
    if (!window) {
        printf("luna-l2: could not initialize window: %s\n", SDL_GetError());
        exit(1);
    }

    SDL_Renderer* renderer = SDL_CreateRenderer(window, -1, SDL_RENDERER_ACCELERATED);
    if (!renderer) {
        printf("luna-l2: could not initialize renderer: %s", SDL_GetError());
        exit(1);
    }

    printf("getting fb");
    nrgba_image* (*return_framebuffer)();
    return_framebuffer = return_component_function(video_component, "return_framebuffer");
    nrgba_image* img = (*return_framebuffer)();
    SDL_Texture* texture = SDL_CreateTexture(renderer, SDL_PIXELFORMAT_ABGR8888, SDL_TEXTUREACCESS_STREAMING, img->width, img->height); 

    printf("updating tex");
    SDL_UpdateTexture(texture, NULL, img->img, img->stride);

    bool quit = false;

    // Event loop
    printf("here we go");
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

        SDL_Delay(10);
    }

    SDL_DestroyWindow(window);
    SDL_Quit();

    return 0;
}
