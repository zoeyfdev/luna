#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "ast.h"
#include "types.h"
#include "../../lcc_shared/libvector.h"
#include "../lex/types.h"
#include "util.h"
#include "../env/env.h"

#define PEEK(by) (peek(tokens, &i, (by)))
#define EXPECT(t) (expect(tokens, &i, (t)))

int add_decl(ast* tu, ast_node* node) {
    vec_grow(tu->declarations, 1);
    ((ast_node**) tu->declarations->data)[tu->declarations->next] = node;
    return tu->declarations->next;
}

c_type* parse_type(vector* tokens, int* _i) {
    int i = *_i;

    int pointer_length = 0;

    bool external = false;
    bool _static = false;
    bool constant = false;

    c_type_base base = INT_BASE_TYPE;
    c_type* type = malloc(sizeof(c_type));

    while (!is_type(PEEK(0))) {
        token* t = PEEK(0);

        if (t->type == TOK_STATIC)
            _static = true;
        else if (t->type == TOK_EXTERN)
            external = true;
        else if (t->type == TOK_CONST)
            constant = false;
        i++;
    }

    token* type_tok = PEEK(0);

    if (!strcmp(type_tok->value, "int")); // no action needed
    else if (!strcmp(type_tok->value, "char"))
        base = TYPE_I8;
    else if (!strcmp(type_tok->value, "void"))
        base = TYPE_VOID;

    while (PEEK(0)->type == TOK_STAR) {
        pointer_length++;
        i++;
    }

    type->type = base;
    type->pointer_length = pointer_length;
    type->external = external;
    type->constant = constant;
    type->_static = _static;

    *_i = i;

    return type;
}

ast* parse(vector* tokens) {
    ast* tu = malloc(sizeof(ast));
    tu->declarations = vec_init(sizeof(ast_node**), 0);

    for (int i = 0; i < tokens->elements; i++) {
        token* t = ((token**) tokens->data)[i];

        if (!strcmp(t->value, "asm") || !strcmp(t->value, "__asm__")) {
            i++;
            EXPECT(TOK_L_PAREN);
            EXPECT(TOK_STRING);
            ast_node* _asm = new_node(AST_ASSEMBLY_DECL);
            _asm->content.assembly.string = PEEK(-1)->value;
            add_decl(tu, _asm);
            EXPECT(TOK_R_PAREN);
        } else if (is_type(t) || is_qualifier(t)) { // add qualifiers
            c_type* t = parse_type(tokens, &i);
            char* name = PEEK(0)->value;

            ast_node* v = new_node(AST_VARIABLE);
            int pos = add_decl(tu, v);
            v->content.variable.name = name;
            v->content.variable.type_info = *t;
           
            if (PEEK(0)->type == TOK_EQUAL) {
                EXPECT(TOK_EQUAL);
            } else if (PEEK(0)->type == TOK_L_PAREN) {
                EXPECT(TOK_L_PAREN);
                EXPECT(TOK_R_PAREN);
                EXPECT(TOK_L_CURLY);
                EXPECT(TOK_R_CURLY);
            }

            add_decl(tu, v);
        } else {
            printf("unexpected tok %s\n", t->value);
        }
    }

    return tu;
}
