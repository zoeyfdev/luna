#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <stdint.h>

#include "../../lcc_shared/libvector.h"
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

vector* add_token(vector* current, vector* tokens, bool no_ins) {
    token* t = malloc(sizeof(token));

    if (strlen((char*) current->data) < 1)
        return current;

    if (is_register((char*) current->data))
        t->type = TYPE_REGISTER;
    else if ((is_instruction((char*) current->data) || ((char*) current->data)[0] == '.') && !no_ins)
        t->type = TYPE_INSTRUCTION;
    else
        t->type = TYPE_TEXT;

    t->value = (char*) current->data;
    t->line = line;

    vec_grow(tokens, 1);
    ((token**) tokens->data)[tokens->next] = t;
    free(current); // Free header without freeing data
    current = vec_init(sizeof(char*), 0);

    return current;
}

vector* lex(char* buffer, uint64_t size) {
    vector* tokens = vec_init(sizeof(token**), 0);
    vector* current = vec_init(sizeof(char*), 0);
    bool comment = false;
    bool in_string = false;

    for (uint64_t i = 0; i < size; i++, buffer++) {
        if (*buffer == 0) 
            break;

        char c = *buffer;

        if (((c == '/' && *(buffer + 1) == '/') || (c == ';') || (c == '#' && *(buffer + 1) == ' ')) && !in_string) {
            if (c != ';')
                buffer++;
            comment = true;
        } else if (c == '\n') {
            if (comment) {
                comment = false;
                continue;
            }
        }

        if (comment)
            continue;

        switch (c) {
        case '\n': 
        case ' ':
            if (!in_string) {
                current = add_token(current, tokens, false);
                if (c == '\n')
                    line++;
            } else {
                vec_grow(current, 1);
                ((char*) current->data)[current->next] = c;
            }
            break;
        case '\r':
            break; // don't process carriage returns
        case '"':
            in_string = !in_string;

            vec_grow(current, 1);
            ((char*) current->data)[current->next] = c;

            if (!in_string) {
                current = add_token(current, tokens, true);
                buffer++;
            } 

            break;
        case ',':
            if (is_register((char*) current->data) && !in_string) {
                current = add_token(current, tokens, false);
                buffer++;
                break;
            } 
        default:
            vec_grow(current, 1);
            ((char*) current->data)[current->next] = c;
            break;
        }
    }

    if (current->elements > 0) {
        current = add_token(current, tokens, false);
    }

    return tokens;
}

