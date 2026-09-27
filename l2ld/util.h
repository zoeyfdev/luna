#pragma once

#include <stddef.h>

#include "types.h"

binding* find_binding(char* name, char* filename);
void write(unsigned char b);
void cleanup_unresolved();
