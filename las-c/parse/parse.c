#include <stdio.h>
#include <string.h>
#include <stdint.h>
#include <stdlib.h>

#include "../lexer/lex.h"
#include "../../lcc_shared/libvector.h"
#include "../../lcc_shared/liberror.h"
#include "../error/error.h"
#include "../util/instructions.h"
#include "../util/registers.h"

char* current_file = "lcc";
bool bits_32 = false;

vector* buffer;
char* filename;

#define EXPECT(T) do { expect(tokens, ((token**) tokens->data)[cursor], (T)); cursor++; } while(0)
#define LAST_TOKEN (((token**) tokens->data)[cursor - 1])
#define THIS_TOKEN (((token**) tokens->data)[cursor])

void expect(vector* tokens, token* t, int type) {
    if (t->type != type) {
        num_errors++;
        lcc_error(NULL, t->line, "unexpected token '", t->value, "'", NULL);
        stargaze(tokens, t);
    }
}

void write(unsigned char b) {
    vec_grow(buffer, 1);
    ((unsigned char*) buffer->data)[buffer->next] = b;
}

void write_str(char* str) {
    for (size_t i = 0; i < strlen(str); i++)
        write(str[i]);
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
            num_errors++;
            lcc_error(NULL, t->line, "unclosed string", NULL);
            stargaze(tokens, t);
        }

        uint32_t num = 0;
        if (!bits_32)
            num = *(uint16_t*) value + (strlen(value) - 3);
        else
            num = *(uint32_t*) value + (strlen(value) - 5);

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

        for (size_t i = 0; i < strlen(value); i++) {
            write(value[i]);
        }
        write(0);
    }
}

void parse_init() {
    buffer = vec_init(sizeof(unsigned char*), 0);
}

vector* parse(vector* tokens);
vector* parse(vector* tokens) {
    int cursor = 0;
    while (true) {
        if (cursor >= tokens->elements)
            break;

        token* t = ((token**) tokens->data)[cursor];
        char* val = t->value;

        cursor++;

        if (val[strlen(val) - 1] == ':') { // label
            if (!bits_32)
                write_str("LD16_");
            else
                write_str("LD32_");

            for (size_t i = 0; i < strlen(val) - 1; i++)
                write(val[i]);

            write(0);

            continue;
        }

        if (!strcmp(val, ".bits")) {
            if (!strcmp(THIS_TOKEN->value, "32")) {
                bits_32 = true;
            } else if (!strcmp(THIS_TOKEN->value, "16")) {
                bits_32 = false;
            } else {
                num_errors++;
                lcc_error(NULL, THIS_TOKEN->line, "invalid value to '.bits', expected '16' or '32'", NULL);
            }
            cursor++;
        } else if (!strcmp(val, "mov")) {
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
        } else if (!strcmp(val, "hlt") || !strcmp(val, "nop")) {
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
        } else if (!strcmp(val, "jnz") || !strcmp(val, "jz")) {
            write(get_opcode(LAST_TOKEN));

            EXPECT(TYPE_REGISTER);

            if (THIS_TOKEN->type != TYPE_REGISTER)
                write(1);
            else
                write(2);

            write(get_reg(LAST_TOKEN));

            if (THIS_TOKEN->type != TYPE_REGISTER)
                insert_any(tokens, THIS_TOKEN);
            else
                write(get_reg(THIS_TOKEN));

            cursor++;
        } else if (!strcmp(val, "cmp") || !strcmp(val, "add") || !strcmp(val, "sub") || !strcmp(val, "mul")
                || !strcmp(val, "div") || !strcmp(val, "igt") || !strcmp(val, "ilt") || !strcmp(val, "and")
                || !strcmp(val, "or") || !strcmp(val, "xor") || !strcmp(val, "mod") || !strcmp(val, "shl") 
                || !strcmp(val, "shr")) {
            write(get_opcode(LAST_TOKEN));
            
            for (int i = 0; i < 3; i++) {
                EXPECT(TYPE_REGISTER);
                write(get_reg(LAST_TOKEN));
            }
        } else if (!strcmp(val, "inc") || !strcmp(val, "dec") || !strcmp(val, "pop")) {
            write(get_opcode(LAST_TOKEN));
            
            EXPECT(TYPE_REGISTER);
            write(get_reg(LAST_TOKEN));
        } else if (!strcmp(val, "push")) {
            write(get_opcode(LAST_TOKEN));

            if (THIS_TOKEN->type != TYPE_REGISTER) {
                write(1);
                insert_any(tokens, THIS_TOKEN);
            } else {
                write(2);
                write(get_reg(THIS_TOKEN));
            }

            cursor++;
        } else if (!strcmp(val, "lod") || !strcmp(val, "str") || !strcmp(val, "lod16") || !strcmp(val, "str16")
                || !strcmp(val, "lod32") || !strcmp(val, "str32") || !strcmp(val, "not")) {
            write(get_opcode(LAST_TOKEN));

            for (int i = 0; i < 2; i++) {
                EXPECT(TYPE_REGISTER);
                write(get_reg(LAST_TOKEN)); 
            }
        } else if (!strcmp(val, ".ascii") || !strcmp(val, ".asciz")) {
            char* value = THIS_TOKEN->value;
            if (value[0] != '"' || value[strlen(value) - 1] != '"') {
                num_errors++;
                lcc_error(NULL, THIS_TOKEN->line, "invalid/unclosed string", NULL);
            }

            for (size_t i = 1; i < strlen(value) - 1; i++) {
                write(value[i]);
            }

            if (!strcmp(val, ".asciz"))
                write(0);

            cursor++;
        } else if (!strcmp(val, "call")) {
            char* name = THIS_TOKEN->value;

            if (!bits_32) {
                char* base = "mov e11, pc\nmov r0, 20\nadd e11, e11, r0\npush e11\njmp ";
                char* stream = calloc(1, strlen(base) + strlen(name) + 2);
                strcat(stream, base);
                strcat(stream, name);
                parse(lex(stream, strlen(stream)));
                free(stream);
            } else {
                char* base = "mov e11, pc\nmov r0, 24\nadd e11, e11, r0\npush e11\njmp ";
                char* stream = calloc(1, strlen(base) + strlen(name) + 2);
                strcat(stream, base);
                strcat(stream, name);
                parse(lex(stream, strlen(stream)));
                free(stream);
            }

            cursor++;
        } else if (!strcmp(val, "ret")) {
            char* stream = "jmp e11";
            parse(lex(stream, strlen(stream)));
        } else {
            num_errors++;
            lcc_error(NULL, t->line, "unrecognized opcode/directive '", t->value, "'", NULL);
        }
    }

    return buffer;
}
