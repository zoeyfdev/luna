#pragma once

extern unsigned char* MEMORY;

void initialize_memory();
unsigned char get_memory(uint32_t address);
void set_memory(uint32_t address, unsigned char value);
