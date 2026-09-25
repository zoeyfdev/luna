#include <stdlib.h>
#include <string.h>

#include "types.h"

extern size_t nbindings;
extern binding** bindings;

extern size_t nbuffer;
extern unsigned char* buffer;

binding* find_binding(char* name, char* filename) {
    for (int i = 0; i < nbindings; i++) {
        if (!strcmp(bindings[i]->name, name) && (!strcmp(bindings[i]->file, filename) || bindings[i]->global))
            return bindings[i];
    }
    return NULL;
}

void* bump_arr(void* arr, size_t nelem, size_t size) {
    return realloc(arr, (nelem + 1) * size);
}

void write(unsigned char b) {
    buffer = bump_arr(buffer, nbuffer, sizeof(unsigned char));
    buffer[nbuffer] = b;
    nbuffer++;
}
