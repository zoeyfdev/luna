#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "ast.h"
#include "types.h"
#include "../../lcc_shared/libvector.h"
#include "../lex/types.h"
#include "util.h"

#define PEEK(by) (peek(tokens, &i, (by)))
#define EXPECT(t) (expect(tokens, &i, (t)))

void add_decl(ast* tu, ast_node* node) {
    vec_grow(tu->declarations, 1);
    ((ast_node**) tu->declarations->data)[tu->declarations->next] = node;
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
            ast_node* _asm = new_node(AST_ASSEMBLY);
            _asm->content.assembly.string = PEEK(-1)->value;
            add_decl(tu, _asm);
            EXPECT(TOK_R_PAREN);
        }
    }

    return tu;
}
