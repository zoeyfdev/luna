#include <stdio.h>
#include <string.h>
#include <stdint.h>
#include <stdlib.h>

#include "../lexer/lex.h"
#include "../../lcc_shared/libvector.h"
#include "../error/error.h"
#include "../util/instructions.h"
#include "../util/registers.h"

char* current_file = "lcc";
bool bits_32 = false;

vector* buffer;

#define EXPECT(T) do { expect(tokens, ((token**) tokens->data)[cursor], (T)); cursor++; } while(0)
#define LAST_TOKEN (((token**) tokens->data)[cursor - 1])
#define THIS_TOKEN (((token**) tokens->data)[cursor])
#define ERROR_OPEN(T) do { num_errors++; fprintf(stderr, "\033[1;31merror:\033[0m %d: ", (T)->line); } while(0)

void expect(vector* tokens, token* t, int type) {
    if (t->type != type) {
        ERROR_OPEN(t);
        fprintf(stderr, "unexpected token '%s'", t->value);
        stargaze(tokens, t);
    }
}

void write(unsigned char b) {
    vec_grow(buffer, 1);
    ((unsigned char*) buffer->data)[buffer->next] = b;
}

unsigned char get_reg(token* t) {
    for (int i = 0; i < NUM_REGISTERS; i++)
        if (!strcmp(t->value, registers[i]))
            return i;
    return 0xff;
}

unsigned char get_opcode(token* t) {
    for (int i = 0; i < NUM_INSTRUCTIONS; i++)
        if (!strcmp(t->value, instructions[i]))
            return i + 1;
    return 0xff;
}

void insert_any(vector* tokens, token* t) {
    char* value = t->value;

    
    // text token
    if ((atoi(value) != 0 && strcmp(value, "0")) || (atoi(value) == 0 && !strcmp(value, "0"))) { // numerical value
        uint32_t num = atoi(value);
        if (!bits_32) {
            write(num >> 8);
            write(num & 0xFF);
        } else {
            write(num >> 24);
            write(num >> 16);
            write(num >> 8);
            write(num & 0xFF);
        }
    } else if (value[0] == '"') { // array of chars
        if (value[strlen(value) - 1] != '"') {
            ERROR_OPEN(t);
            fprintf(stderr, "unclosed string\n", t->line);
            stargaze(tokens, t);
        }

        uint32_t num = 0;
        if (!bits_32)
            num = *(uint16_t*) strlen(value) - 3;
        else
            num = *(uint32_t*) strlen(value) - 5;

        if (!bits_32) {
            write(num >> 8);
            write(num & 0xFF);
        } else {
            write(num >> 24);
            write(num >> 16);
            write(num >> 8);
            write(num & 0xFF);
        }
    } else { // label reference
        write('L');
        write('R');
        write('_');

        for (int i = 0; i < strlen(value); i++) {
            write(value[i]);
        }
        write(0);
    }
}

vector* parse(vector* tokens) {
    buffer = vec_init(sizeof(unsigned char*), 0);

    int cursor = 0;
    while (true) {
        if (cursor >= tokens->elements)
            break;

        token* t = ((token**) tokens->data)[cursor];
        char* val = t->value;

        cursor++;
        if (!strcmp(val, "mov")) {
            write(get_opcode(LAST_TOKEN));

            EXPECT(TYPE_REGISTER);

            if (THIS_TOKEN->type != TYPE_REGISTER) {
                write(1);
                write(get_reg(LAST_TOKEN));
                insert_any(tokens, THIS_TOKEN);
            } else {
                write(2);
                write(get_reg(LAST_TOKEN));
                write(get_reg(THIS_TOKEN));
            }

            cursor++;
        } else if (!strcmp(val, "hlt")) {
            write(get_opcode(LAST_TOKEN));
        } else if (!strcmp(val, "jmp")) {
            write(get_opcode(LAST_TOKEN));
            if (THIS_TOKEN->type != TYPE_REGISTER) {
                write(1);
                insert_any(tokens, THIS_TOKEN);
            } else {
                write(2);
                write(get_reg(THIS_TOKEN));
            }

            cursor++;
        } else if (!strcmp(val, "int")) {
            write(get_opcode(LAST_TOKEN));
            insert_any(tokens, THIS_TOKEN);

            cursor++;
        } else if (!strcmp(val, "jnz")) {
            write(get_opcode(LAST_TOKEN));

            EXPECT(TYPE_REGISTER);
            write(get_reg(LAST_TOKEN));



            insert_any(tokens, T)
        } else {
            ERROR_OPEN(t);
            fprintf(stderr, "unrecognized opcode/directive '%s'\n", t->value);
        }
    }

    return buffer;
}
