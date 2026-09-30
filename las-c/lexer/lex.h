#pragma once

#include "../../lcc_shared/libvector.h"

typedef struct {
    int type;
    char* value;
} token;

vector* lex(char* buffer);


