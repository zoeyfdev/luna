#pragma once

#include <stdint.h>

#include "../../lcc_shared/libvector.h"

#define TYPE_REGISTER 1
#define TYPE_INSTRUCTION 2
#define TYPE_TEXT 3

typedef struct {
    int type;
    int line;
    char* value;
} token;

vector* lex(char* buffer, uint64_t size);
bool is_register(char* str);
bool is_instruction(char* str);
