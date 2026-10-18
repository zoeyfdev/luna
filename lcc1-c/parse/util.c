#include <stdlib.h>
#include <string.h>
#include <stdio.h>

#include "ast.h"
#include "../lex/types.h"
#include "../../lcc_shared/liberror.h"
#include "../../lcc_shared/libvector.h"

ast_node* new_node(ast_node_type type) {
    ast_node* n = calloc(1, sizeof(ast_node));
    n->type = type;

    if (type == AST_VARIABLE) {
        n->content.variable.parameters = vec_init(sizeof(variable**), 0);
        n->content.variable.children = vec_init(sizeof(ast_node**), 0);
    }

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
        printf("out of bounds\n");
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

bool is_qualifier(token* t) {
    if (t->type == TOK_SIGNED || t->type == TOK_UNSIGNED || t->type == TOK_SHORT || t->type == TOK_LONG ||
        t->type == TOK_STATIC || t->type == TOK_EXTERN)
        return true;
    return false;
}

void add_node(ast_node* p, ast_node* v) {
    printf("added node\n");
    if (p->type == AST_VARIABLE) {
        vec_grow(p->content.variable.children, 1);
        ((ast_node**) p->content.variable.children->data)[p->content.variable.children->next] = v;
    } else {
        lcc_ice(NULL, 0, "invalid variable type to add_node()", NULL);
    }
}
