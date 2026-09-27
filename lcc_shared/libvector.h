#pragma once

typedef struct {
    size_t size;
    int elements;
    int next;
    void* data;
} vector;

void vec_grow(vector* v, int num_elements);
vector* vec_init(size_t size, int init_elements);
