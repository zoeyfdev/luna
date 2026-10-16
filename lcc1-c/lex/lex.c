#include <stdlib.h>
#include <string.h>

#include "../../lcc_shared/libvector.h"
#include "../../lcc_shared/libcpp.h"
#include "../../lcc_shared/libstoi.h"
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
    else if (!strcmp(value, "{"))
        return TOK_L_CURLY;
    else if (!strcmp(value, "}"))
        return TOK_R_CURLY;
    else if (!strcmp(value, "="))
        return TOK_EQUAL;
    else if (!strcmp(value, "asm") || !strcmp(value, "__asm__"))
        return TOK_ASM;
    else if (!strcmp(value, "short"))
        return TOK_SHORT;
    else if (!strcmp(value, "long"))
        return TOK_LONG;
    else if (!strcmp(value, "signed"))
        return TOK_SIGNED;
    else if (!strcmp(value, "unsigned"))
        return TOK_UNSIGNED;
    else if (!strcmp(value, "extern"))
        return TOK_EXTERN;
    else if (!strcmp(value, "static"))
        return TOK_STATIC;
    else if (!strcmp(value, "const"))
        return TOK_CONST;
    else { 
        bool worked = false;
        stoi(value, &worked);

        if (worked)
            return TOK_NUMBER;
        return TOK_IDENTIFIER;
    }

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

        if (_t->value[0] == '"' && _t->value[strlen(_t->value) - 1] == '"') {
            char* v = calloc(1, strlen(_t->value));
            for (size_t i = 1; i < strlen(_t->value) - 1; i++)
                v[i - 1] = _t->value[i];
            free(_t->value);
            _t->value = v;
            _t->type = TOK_STRING;
        } else
            _t->type = return_tok_type(_t->value);

        vec_grow(output, 1);
        ((token**) output->data)[output->next] = _t;
    }

    return output;
}
