#pragma once

#include <stdint.h>

typedef struct {
    char* name;
    uint64_t size;
    unsigned char* data;
} file;

typedef struct { 
    char* file;
    char* name;
    bool global;
    bool is_32;
    uint64_t location;
} binding;

typedef struct {
    char* name;
    char* file;
    bool solved;
    uint64_t location;
} unresolved_binding;
