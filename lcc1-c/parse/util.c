#include <stdlib.h>
#include <string.h>
#include <stdio.h>

#include "ast.h"
#include "../lex/types.h"

ast_node* new_node(ast_node_type type) {
    ast_node* n = calloc(1, sizeof(ast_node));

    n->type = type;

    return n;
}

bool is_type(token* t) {
    char* v = t->value;
    if (!strcmp(v, "int") || !strcmp(v, "char") || !strcmp(v, "void"))
        return true;
    return false;
}

token* peek(vector* tokens, int* i, int by) {
    if (*i >= tokens->elements) {
        token* t = calloc(1, sizeof(token));
        t->type = TOK_NONE;
        return t;
    }

    return ((token**) tokens->data)[*i + by];
}

token* expect(vector* tokens, int* i, token_type t) {
    token* _t = peek(tokens, i, 0);
    if (_t->type != t)
        printf("invalid token %s\n", _t->value);
    (*i)++;

    return _t;
}
