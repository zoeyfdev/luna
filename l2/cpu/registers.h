#pragma once

#include <stdint.h>

#include "register_defs.h"

void initialize_registers();
uint32_t get_register(unsigned char address);
void set_register(unsigned char address, uint32_t value);
char* get_register_name(unsigned char address);
void reg_dump(); 

#define CPU_FLAG_XEN 0b00000000000000000000000000000001
#define CPU_FLAG_IIF 0b00000000000000000000000000000010

#define IS_XEN ((get_register(0x22) & (1 << (CPU_FLAG_XEN - 1))))
#define IS_IIF ((get_register(0x22) & (1 << (CPU_FLAG_IIF - 1))))
