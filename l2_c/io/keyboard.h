#pragma once

extern unsigned char* KEYBOARD_MEMORY;

int keyboard_upper(int code);
int keyboard_lower(int code);
unsigned char keyboard_read_memory(uint32_t address);
void keyboard_write_memory(uint32_t address, unsigned char value);
