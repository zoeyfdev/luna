#include <stdlib.h>
#include <string.h>

#include "../../lcc_shared/libvector.h"
#include "../../lcc_shared/libcpp.h"
#include "types.h"

#define INSERT_CHAR() do { \
    vec_grow(current, 1); \
    ((char*) current->data)[current->next] = c; \
} while (0)

token_type return_tok_type(char* value) {
    if (!strcmp(value, "+"))
        return TOK_PLUS;
    else if (!strcmp(value, "-"))
        return TOK_MINUS;
    else if (!strcmp(value, "*"))
        return TOK_STAR;
    else if (!strcmp(value, "/"))
        return TOK_SLASH;
    else if (!strcmp(value, "("))
        return TOK_L_PAREN;
    else if (!strcmp(value, ")"))
        return TOK_R_PAREN;
    else
        return TOK_NUMBER;

    exit(1);
    return TOK_NUMBER;
}

vector* lex(vector* buffer) {
    vector* output = vec_init(sizeof(token**), 0);

    for (int i = 0; i < buffer->elements; i++) {
        cpp_token* t = ((cpp_token**) buffer->data)[i];

        token* _t = malloc(sizeof(token));
        _t->line = t->line;
        _t->filename = t->filename;
        _t->value = t->value;
        _t->displayed_value = t->value;
        _t->type = return_tok_type(_t->value);

        vec_grow(output, 1);
        ((token**) output->data)[output->next] = _t;
    }
}
