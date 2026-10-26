#pragma once

#include <stdint.h>

#include "libvector.h"

char* lfn_get_ext(char* s);
vector* get_file_word(char* file_data, uint64_t end, int pos);
char* lfn_get_base(char* _full);
char* get_file_lead(char* _fn);
char* open_file(char* filename);
void write_file(char* filename, unsigned char* buffer, size_t len);
