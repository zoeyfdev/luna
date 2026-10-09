#include <stdio.h>
#include <string.h>
#include <stdint.h>
#include <stdlib.h>

#include "../lexer/lex.h"
#include "../../lcc_shared/libvector.h"
#include "../../lcc_shared/liberror.h"
#include "../../lcc_shared/libfile.h"
#include "../error/error.h"
#include "../util/instructions.h"
#include "../util/registers.h"
#include "../../lcc_shared/libstoi.h"
#include "../util/push.h"

char* current_file = "lcc";
bool bits_32 = false;

vector* buffer;
char* filename = NULL;

#define EXPECT(T) do { expect(tokens, ((token**) tokens->data)[cursor], (T)); cursor++; } while(0)
#define LAST_TOKEN (((token**) tokens->data)[cursor - 1])
#define THIS_TOKEN (((token**) tokens->data)[cursor])
#define NEXT_TOKEN (((token**) tokens->data)[cursor + 1])

void expect(vector* tokens, token* t, int type) {
    if (t->type != type) {
        num_errors++;
        lcc_error(filename, t->line, "unexpected token '", t->value, "'", NULL);
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
    bool worked = false;
    int64_t num = stoi(value, &worked);
    if (worked) { // numerical value
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
            lcc_error(filename, t->line, "unclosed string", NULL);
            stargaze(tokens, t);
        }

        int bits = 0;

        uint64_t num = 0;
        for (size_t i = 1; i < strlen(value) - 1; i++, bits += 8) {
            char c = value[i];
            if (bits == 0) {
                num = c & 0xFF;
            } else {
                num <<= 8;
                num |= c;
            }
        } 

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
        write_str("LR_");

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
                write_str("L_32BIT");
            } else if (!strcmp(THIS_TOKEN->value, "16")) {
                bits_32 = false;
                write_str("L_16BIT");
            } else {
                num_errors++;
                lcc_error(filename, THIS_TOKEN->line, "invalid value to '.bits', expected '16' or '32'", NULL);
            }
            cursor++;
        } else if (!strcmp(val, "set")) {
            write(get_opcode(LAST_TOKEN));

            if (!strcmp(THIS_TOKEN->value, "32")) {
                write(1);
            } else if (!strcmp(THIS_TOKEN->value, "16")) {
                write(0);
            } else {
                num_errors++;
                lcc_error(filename, THIS_TOKEN->line, "invalid value to 'set', expected '16' or '32'", NULL);
            }

            cursor++;
        } else if (!strcmp(val, ".org")) {
            write_str("LO_");
            
            bool worked = false;
            int64_t num = stoi(THIS_TOKEN->value, &worked);

            if (!worked) {
                num_errors++;
                lcc_error(filename, THIS_TOKEN->line, "invalid number to .org, got '", THIS_TOKEN->value, "'", NULL);
            }

            write(num >> 8);
            write(num & 0xFF);

            cursor++;
        } else if (!strcmp(val, "mov")) {
            write(get_opcode(LAST_TOKEN));

            EXPECT(TYPE_REGISTER);

            if (THIS_TOKEN->type != TYPE_REGISTER) {
                write(1);
                write(get_reg(LAST_TOKEN));
                insert_any(tokens, THIS_TOKEN);
            } else {
                if (strcmp(NEXT_TOKEN->value, "+") && strcmp(NEXT_TOKEN->value, "-")) {
                    write(2);
                    write(get_reg(LAST_TOKEN));
                    write(get_reg(THIS_TOKEN));
                } else {
                    write(3);
                    write(get_reg(LAST_TOKEN));
                    write(get_reg(THIS_TOKEN));
                    cursor += 2;
                    write(!strcmp(LAST_TOKEN->value, "+") ? 1 : 2);
                    insert_any(tokens, THIS_TOKEN);
                }
            }

            cursor++;
        } else if (!strcmp(val, "hlt") || !strcmp(val, "nop") || !strcmp(val, "sti")) {
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
                || !strcmp(val, "or") || !strcmp(val, "mod") || !strcmp(val, "shl")  || !strcmp(val, "shr")
                || !strcmp(val, "xor")) {
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
                || !strcmp(val, "lod32") || !strcmp(val, "str32") || !strcmp(val, "lod_ptr")
                || !strcmp(val, "str_ptr") || !strcmp(val, "not")) {

            if (strcmp(val, "lod_ptr") && strcmp(val, "str_ptr"))
                write(get_opcode(LAST_TOKEN));
            else {
                if (!bits_32) {
                    if (!strcmp(val, "lod_ptr")) {
                        write(0x19);
                    } else {
                        write(0x18);
                    }
                } else {
                    if (!strcmp(val, "lod_ptr")) {
                        write(0x1e);
                    } else {
                        write(0x1f);
                    }
                }
            }

            for (int i = 0; i < 2; i++) {
                EXPECT(TYPE_REGISTER);
                write(get_reg(LAST_TOKEN)); 
            }
        } else if (!strcmp(val, ".ascii") || !strcmp(val, ".asciz")) {
            char* value = THIS_TOKEN->value;

            if (value[0] != '"' || value[strlen(value) - 1] != '"') {
                num_errors++;
                lcc_error(filename, THIS_TOKEN->line, "invalid/unclosed string", NULL);
            }

            for (size_t i = 1; i < strlen(value) - 1; i++) {
                if (strlen(value) > i + 1) {
                    if (value[i] == '\\' && value[i + 1] == 'n') {
                        write('\n');
                        i++;
                        continue;
                    }
                }

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
        } else if (!strcmp(val, "pusha")) {
            parse(lex(pusha_str, strlen(pusha_str)));
        } else if (!strcmp(val, "popa")) {
            parse(lex(popa_str, strlen(popa_str)));
        } else if (!strcmp(val, ".byte") || !strcmp(val, ".word") || !strcmp(val, ".dword")) {
            int line = THIS_TOKEN->line;
            while (cursor < tokens->elements) {
                if (THIS_TOKEN->line != line)
                    break;

                bool worked = false;
                int64_t n = stoi(THIS_TOKEN->value, &worked);

                if (!worked) {
                    num_errors++;
                    lcc_error(filename, t->line, "invalid number '", THIS_TOKEN->value, "'", NULL);
                }

                if (!strcmp(val, ".byte")) {
                    write(n & 0xFF);
                } else if (!strcmp(val, ".word")) {
                    write(n >> 8);
                    write(n & 0xFF);
                } else if (!strcmp(val, ".dword")) {
                    write(n >> 24);
                    write(n >> 16);
                    write(n >> 8);
                    write(n & 0xFF);
                }

                cursor++; 
            }
        } else if (!strcmp(val, ".ptr")) {
            insert_any(tokens, THIS_TOKEN);
            cursor++;
        } else if (!strcmp(val, ".embed")) {
            char* fn = THIS_TOKEN->value;
            char* lead = get_file_lead(filename);

            char* full_fn = calloc(1, strlen(lead) + strlen(fn) + 1); 

            if (strlen(fn) == 0) {
                num_errors++;
                lcc_error(filename, t->line, "invalid filename '", fn, "'", NULL);
                goto done;
            }

            if (fn[0] != '"' || fn[strlen(fn) - 1] != '"') {
                num_errors++;
                lcc_error(filename, t->line, "invalid/unclosed string", NULL);
                goto done;
            }

            char* _fn = calloc(1, strlen(fn) + 1);
            for (size_t i = 1; i < strlen(fn) - 1; i++)
                _fn[i - 1] = fn[i];

            strcat(full_fn, lead);
            strcat(full_fn, _fn);

            FILE* f = fopen(full_fn, "rb");
            if (f == NULL) {
                num_errors++;
                lcc_error(filename, t->line, "could not open '", _fn, "': no such file or directory", NULL);
                free(_fn);
                goto done;
            }

            fseek(f, 0, SEEK_END);
            uint64_t size = ftell(f);
            fseek(f, 0, SEEK_SET);
            unsigned char* buffer = malloc(size);

            fread(buffer, sizeof(unsigned char), size, f);
            
            for (uint64_t i = 0; i < size; i++)
                write(buffer[i]);

            free(_fn);
            free(buffer);
            fclose(f);
            free(full_fn);
            done:
            cursor++;
        } else if (!strcmp(val, ".global")) {
            char* label = THIS_TOKEN->value;

            write_str("L_GLOBL_");
            write_str(label);
            write(0);

            cursor++;
        } else if (!strcmp(val, ".pad")) {
            bool worked = false;
            int64_t num = stoi(THIS_TOKEN->value, &worked);

            if (!worked) {
                num_errors++;
                lcc_error(filename, t->line, "invalid number '", THIS_TOKEN->value, "'", NULL);
            }

            for (int64_t i = 0; i < num; i++) {
                write(0);
            }

            cursor++;
        } else if (!strcmp(val, ".fill")) {
            bool worked = false;
            int64_t num = stoi(THIS_TOKEN->value, &worked);

            if (!worked) {
                num_errors++;
                lcc_error(filename, t->line, "invalid number '", THIS_TOKEN->value, "'", NULL);
            }

            write_str("LP_");
            write(num >> 8);
            write(num & 0xFF);

            cursor++;
        } else {
            num_errors++;
            lcc_error(filename, t->line, "unknown instruction/directive '", t->value, "'", NULL);
        }
    }

    return buffer;
}
