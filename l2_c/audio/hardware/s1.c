#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <SDL2/SDL.h>
#include <pthread.h>

#include "../../util/psleep.c" // i want it to be standalone

// memory layout
#define AUDIO_RAM_SIZE 10

unsigned char (*get_memory)(uint32_t);

unsigned char* AUDIO_MEMORY = NULL;

SDL_AudioDeviceID device;

bool audio_initialized = false;

void* audio_poll(void* VOID) {
    for (;;) {
        if (AUDIO_MEMORY[0] != 0) {
            AUDIO_MEMORY[0] = 0;

            uint32_t cursor = ((uint32_t) AUDIO_MEMORY[1]) << 24
                | ((uint32_t) AUDIO_MEMORY[2]) << 16
                | ((uint32_t) AUDIO_MEMORY[3]) << 8
                | ((uint32_t) AUDIO_MEMORY[4]);

            uint32_t length = ((uint32_t) AUDIO_MEMORY[5]) << 24
                | ((uint32_t) AUDIO_MEMORY[6]) << 16
                | ((uint32_t) AUDIO_MEMORY[7]) << 8
                | ((uint32_t) AUDIO_MEMORY[8]);

            unsigned char* buffer = malloc(length);
            for (uint32_t i = 0; i < length; i++, cursor++)
                buffer[i] = get_memory(cursor);

            SDL_PauseAudioDevice(device, false);
            SDL_QueueAudio(device, buffer, length);

            while (SDL_GetQueuedAudioSize(device) > 0)
                psleep(15);

            AUDIO_MEMORY[9] = 1;
            free(buffer);
        }
        psleep(15);
    }
    return NULL;
}

void audio_init(unsigned char (*_get_memory)(uint32_t)) {
    if (audio_initialized == true) return;
    audio_initialized = true;
    AUDIO_MEMORY = malloc(AUDIO_RAM_SIZE);
    get_memory = _get_memory;

    if (SDL_Init(SDL_INIT_AUDIO) < 0) {
        printf("luna-l2: could not initialize audio: %s\n", SDL_GetError());
        return;
    }

    SDL_AudioSpec spec;
    spec.freq = 48000;
    spec.format = AUDIO_S8;
    spec.channels = 1;
    spec.samples = 4096;
    SDL_AudioSpec obtained;

    device = SDL_OpenAudioDevice(NULL, false, &spec, &obtained, SDL_AUDIO_ALLOW_FORMAT_CHANGE);

    printf("Obtained data:\nFreq: %d\nFormat: %d\nChannels: %d\nSamples: %d\n", obtained.freq, obtained.format, obtained.channels, obtained.samples);

    if (device == 0) {
        printf("luna-l2: could not initialize audio: %s\n", SDL_GetError());
        return;
    }

    pthread_t poll_thread;
    pthread_create(&poll_thread, NULL, audio_poll, NULL);
}

unsigned char audio_read_memory(uint32_t address) {
    if (address < AUDIO_RAM_SIZE)
        return AUDIO_MEMORY[address];
    return (unsigned char) rand() & 0xFF; 
}

void audio_write_memory(uint32_t address, unsigned char value) {
    if (address < AUDIO_RAM_SIZE)
        AUDIO_MEMORY[address] = value;
}
