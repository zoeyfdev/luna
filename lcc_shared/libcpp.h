#pragma once

#include <string.h>

#include "libvector.h"

typedef struct {
    char* value;
    char* filename;
    int line;
} cpp_token;

vector* cpp(char* filename, char* buffer, char* match_str);
