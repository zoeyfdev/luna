#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <stdint.h>

#include "../../lcc_shared/libvector.h"
#include "../../lcc_shared/libcpp.h"
#include "../util/registers.h"
#include "../util/instructions.h"

#define TYPE_REGISTER 1
#define TYPE_INSTRUCTION 2
#define TYPE_TEXT 3

typedef struct {
    int type;
    int line;
    char* value;
} token;

bool is_register(char* str) {
    for (int i = 0; i < NUM_REGISTERS; i++)
        if (!strcmp(str, registers[i]))
            return true;
    return false;
}

bool is_instruction(char* str) {
    for (int i = 0; i < NUM_INSTRUCTIONS; i++)
        if (!strcmp(str, instructions[i]))
            return true;
    return false;
}

int line = 1;

vector* lex(vector* buffer) {
    vector* output = vec_init(sizeof(token**), 0);
    for (int i = 0; i < buffer->elements; i++) {
        cpp_token* t = ((cpp_token**) buffer->data)[i];

        if (!strcmp(t->value, ",")) {
            if (i - 1 >= 0 && is_register((char*) ((cpp_token**) buffer->data)[i - 1]->value))
                continue;
        }

        token* _t = malloc(sizeof(token));
        _t->line = t->line;
        _t->value = t->value;

        if (is_register(_t->value))
            _t->type = TYPE_REGISTER;
        else if (is_instruction(_t->value) || _t->value[0] == '.')
            _t->type = TYPE_INSTRUCTION;
        else
            _t->type = TYPE_TEXT;

        vec_grow(output, 1);
        ((token**) output->data)[output->next] = _t;
    }

    return output;
}

