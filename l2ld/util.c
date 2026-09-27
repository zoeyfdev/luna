#include <stdlib.h>
#include <string.h>
#include <stdio.h>

#include "types.h"
#include "../lcc_shared/libvector.h"

extern vector* bindings;
extern vector* unresolved_bindings;
extern vector* buffer;

binding* find_binding(char* name, char* filename) {
    for (int i = 0; i < bindings->elements; i++) {
        if (!strcmp(((binding**) bindings->data)[i]->name, name) && (!strcmp(((binding**) bindings->data)[i]->file, filename) || ((binding**) bindings->data)[i]->global))
            return ((binding**) bindings->data)[i];
    }
    return NULL;
}

void write(unsigned char b) {
    vec_grow(buffer, 1);
    ((unsigned char*) buffer->data)[buffer->next] = b;
}

void cleanup_unresolved() {
    for (int i = 0; i < unresolved_bindings->elements; i++) {
        unresolved_binding* ub = ((unresolved_binding**) unresolved_bindings->data)[i];

        binding* decl = find_binding(ub->name, ub->file);
        if (decl != NULL) {
            ub->solved = true;
            if (!decl->is_32) {
                ((unsigned char*) buffer->data)[ub->location] = decl->location >> 8;
                ((unsigned char*) buffer->data)[ub->location + 1] = decl->location & 0xFF;
            } else {
                ((unsigned char*) buffer->data)[ub->location] = decl->location >> 24;
                ((unsigned char*) buffer->data)[ub->location + 1] = decl->location >> 16;
                ((unsigned char*) buffer->data)[ub->location + 2] = decl->location >> 8;
                ((unsigned char*) buffer->data)[ub->location + 3] = decl->location & 0xFF;
            } 
        }
    }
}
