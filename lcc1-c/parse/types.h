#pragma once

#include <stdint.h>

#include "../lex/types.h"

typedef enum {
    TYPE_NONE,
    TYPE_I8,
    TYPE_I16,
    TYPE_I32,
    TYPE_I64,
    TYPE_STRUCT,
    TYPE_VOID,
} c_type_base;

typedef struct {
    c_type_base type;

    char* name;

    int pointer_length;
    int function_length;
    uint32_t offset;
    uint32_t size;

    bool constant;
    bool external;
    bool _static;

    struct c_type** children;
} c_type;

typedef struct {
    int scope;
    int parent;
} scope;
