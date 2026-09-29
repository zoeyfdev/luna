#include <string.h>
#include <stdint.h>
#include <stdlib.h>
#include <ctype.h>
#include <stdio.h>

#include "libvector.h"

void rev(char* s) {
    int l = 0;
    int r = strlen(s) - 1;
    char t;

    while (l < r) { 
        t = s[l];
        s[l] = s[r];
        s[r] = t;

        l++;
        r--;
    }
}

char* lfn_get_ext(char* s) {
    int len = strlen(s) - 1; // Get last character

    int i = len;
    bool exit = false;
    for (; i >= 0; i--) {
        s[i] = tolower(s[i]); // just in case
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

char* lfn_get_base(char* full) {
    vector* name = vec_init(sizeof(char*), 0);

    for (size_t i = 0; i < strlen(full); i++) {
        char c = full[i];

        if (c != '/' && c != '\\') {
            vec_grow(name, 1);
            ((char*) name->data)[name->next] = c;
        } else {
            free(name->data);
            free(name);
            name = vec_init(sizeof(char*), 0);
        }
    }

    for (int i = name->elements; i >= 0; i--) {
        if (((char*) name->data)[i] == '.') {
            ((char*) name->data)[i] = 0;
            break;
        }
        ((char*) name->data)[i] = 0;
    }

    return ((char*) name->data);
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
