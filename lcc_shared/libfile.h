#pragma once

#include <stdint.h>

#include "libvector.h"

char* lfn_get_ext(char* s);
vector* get_file_word(char* file_data, uint64_t end, int pos);
char* lfn_get_base(char* _full);
