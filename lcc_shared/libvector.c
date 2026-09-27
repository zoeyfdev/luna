#include "stdlib.h"

typedef struct {
    size_t size;
    int elements;
    int next;
    void* data;
} vector;

void vec_grow(vector* v, int num_elements) {
    v->elements += num_elements;
    v->next = v->elements - 1;
    v->data = realloc(v->data, v->size * v->elements); 
}

vector* vec_init(size_t size, int init_elements) {
    vector* v = malloc(sizeof(vector));

    v->elements = init_elements;
    v->next = v->elements - 1;
    v->size = size;
    v->data = calloc(init_elements, size);

    return v;
}

