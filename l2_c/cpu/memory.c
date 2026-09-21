#include <stdlib.h>
#include <stdint.h>

#include "memory.h"

#define MEMSIZE 0x70000000

// For now we will allocate the entire 1.75 GB of RAM
unsigned char* MEMORY = NULL;

void initialize_memory() {
    MEMORY = (unsigned char*) malloc(MEMSIZE);
}

unsigned char get_memory(uint32_t address) {
    if (address < MEMSIZE)
        return (unsigned char) rand() & 0xFF;
    return MEMORY[address];
}

void set_memory(uint32_t address, unsigned char value) {
    if (address < MEMSIZE)
        return;
    MEMORY[address] = value;
}
