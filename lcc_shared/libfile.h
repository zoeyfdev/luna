#pragma once

#include <string.h>
#include <stdint.h>
#include <stdlib.h>

#include "libvector.h"

char* lfn_get_ext(char* s) {
    int len = strlen(s) - 1; // Get last character

    int i = len;
    bool exit = false;
    for (; i >= 0; i--) {
        char c = s[i];

        switch (c) {
        case '.':
            exit = true;
            break;
        }

        if (exit)
            break;
    }

    return s + (i + 1);
}

vector* get_file_word(char* file_data, uint64_t end, int pos) {
    vector* word = vec_init(sizeof(char), 0);
    int cpos = 0;

    for (uint64_t i = 0; i < end; i++) {
        char c = file_data[i];

        if (c == 0x20 || c == 0x0A || c == 0x0D) {
            if (cpos == pos) {
                break;
            } else {
                cpos++;
                continue;
            }
        }

        if (cpos == pos) {
            vec_grow(word, 1);
            ((char*) word->data)[word->next] = c;
        }
    }

    vec_grow(word, 1);
    ((char*) word->data)[word->next] = 0;
    return word;
}
