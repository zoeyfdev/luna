#pragma once

#include <stddef.h>

#include "types.h"

void* bump_arr(void* arr, size_t nelem, size_t size);
binding* find_binding(char* name);
void write(unsigned char b);
