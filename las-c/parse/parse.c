#include <stdio.h>

#include "../lexer/lex.h"
#include "../../lcc_shared/libvector.h"
#include "../error/error.h"

char* current_file = "lcc";
bool bits_32 = false;

#define EXPECT(T) do { expect(tokens, t, (T)); cursor++; } while(0)

void expect(vector* tokens, token* t, int type) {
    if (t->type != type) {
        fprintf(stderr, "\033[1;31merror:\033[0m%d: unexpected token '%s'", t->line, t->value);
        stargaze(tokens, t);
    }
}

void parse(vector* tokens) {
    vector* buffer = vec_init(sizeof(unsigned char*), 0);

    int cursor = 0;
    while (true) {
        if (cursor >= tokens->elements)
            break;

        token* t = ((token**) tokens->data)[cursor];

        if (!strcmp("mov")) {
            EXPECT(TYPE_INSTRUCTION);
            
        }
    }
}
