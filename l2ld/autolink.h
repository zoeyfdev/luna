#pragma once

#include "../lcc_shared/libvector.h"

typedef struct {
    char* label;
    char* file;
} autolink_pair;

vector* autolink();
