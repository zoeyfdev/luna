#include "ast.h"
#include "types.h"
#include "../lex/types.h"
#include "../../lcc_shared/libvector.h"
#include "util.h"

#define PEEK(by) (peek(tokens, &i, (by)))
#define EXPECT(t) (expect(tokens, &i, (t)))

void parse_statement(ast* tu, vector* tokens, int i) {
    token* t = ((token**) tokens->data)[i];

    if (t->type == TOK_ASM) {
         
    }

    return i;
}

int parse_local(ast* tu, vector* tokens, int start, ast_node* decl) {
    int i = start;
    for (; i < tokens->elements; i++) {
        i = parse_statement(tokens, i);
    }

    return i;
}
