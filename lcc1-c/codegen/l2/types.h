#pragma once

#include "../../parse/types.h"

typedef struct {
    c_type type_info;
    char* _register;
    bool rvalue;
    bool read;
} codegen_return;
