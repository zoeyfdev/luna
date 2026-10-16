#include "ast.h"
#include "types.h"
#include "../lex/types.h"
#include "../../lcc_shared/libvector.h"

void parse_statement(vector* tokens, int i) {
    token* t = ((token**) tokens->data)[i];

    if (t->type == TOK_ASM) {

    }
}

int parse_local(vector* tokens, int start, ast_node* decl) {
    int i = start;
    for (; i < tokens->elements; i++) {
        i = parse_statement(tokens, i);
    }

    return i;
}
