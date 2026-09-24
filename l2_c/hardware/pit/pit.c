#include <stdlib.h>
#include <stdint.h>
#include <pthread.h>
#include <stdio.h>
#include <string.h>

#include "../../cpu/register_defs.h"
#include "../../util/psleep.c"

#define PIT_RAM_SIZE 8

bool pit_initialized = false;
unsigned char* PIT_MEMORY = NULL;

void (*set_register)(unsigned char, uint32_t);
uint32_t (*get_register)(unsigned char);

void* pit_poll(void* VOID) {
    uint32_t current; 

    for (;;) {
        current = (uint32_t) PIT_MEMORY[4] << 24
            | (uint32_t) PIT_MEMORY[5] << 16
            | (uint32_t) PIT_MEMORY[6] << 8
            | (uint32_t) PIT_MEMORY[7];

        if (current == 0) {
            PIT_MEMORY[4] = PIT_MEMORY[0];
            PIT_MEMORY[5] = PIT_MEMORY[1];
            PIT_MEMORY[6] = PIT_MEMORY[2];
            PIT_MEMORY[7] = PIT_MEMORY[3];
            set_register(IR, get_register(IR) | (1 << 1));
        } else {
            current--;
            PIT_MEMORY[4] = current >> 24;
            PIT_MEMORY[5] = current >> 16;
            PIT_MEMORY[6] = current >> 8;
            PIT_MEMORY[7] = current & 0xFF;
        }
        psleep(1);
    }
    return NULL;
}

void pit_init(void (*_set_register)(unsigned char, uint32_t), uint32_t (*_get_register)(unsigned char)) {
    if (pit_initialized == true) return;
    pit_initialized = true;

    set_register = _set_register;
    get_register = _get_register;

    PIT_MEMORY = malloc(PIT_RAM_SIZE);
    memset(PIT_MEMORY, 0x00, PIT_RAM_SIZE);

    PIT_MEMORY[2] = 0x03;
    PIT_MEMORY[3] = 0xE8;
    PIT_MEMORY[6] = 0x03;
    PIT_MEMORY[7] = 0xE8;

    pthread_t pit_thread;
    pthread_create(&pit_thread, NULL, pit_poll, NULL);
}

void pit_write_memory(uint32_t address, unsigned char value) {
    if (address < PIT_RAM_SIZE)
        PIT_MEMORY[address] = value;
}

unsigned char pit_read_memory(uint32_t address) {
    if (address < PIT_RAM_SIZE)
        return PIT_MEMORY[address];
    return (unsigned char) rand() & 0xFF;
}
