#include "ast.h"
#include "types.h"
#include "../lex/types.h"
#include "../../lcc_shared/libvector.h"
#include "util.h"

#define PEEK(by) (peek(tokens, &i, (by)))
#define EXPECT(t) (expect(tokens, &i, (t)))

void parse_statement(vector* tokens, int i, ast_node* decl) {
    token* t = ((token**) tokens->data)[i];

    if (t->type == TOK_ASM) {
        EXPECT(TOK_ASM);
        EXPECT(TOK_L_PAREN);
        token* s = EXPECT(TOK_STRING);

        ast_node* _asm_node = new_node(AST_ASM_STATEMENT);
        _asm_node->content.assembly.string = s->value;
        add_node(decl, _asm_node);

        EXPECT(TOK_R_PAREN);
        EXPECT(TOK_SEMI);
    } else {
        // assume expression
    }

    return i;
}

int parse_local(vector* tokens, int start, ast_node* decl) {
    int i = start;
    for (; i < tokens->elements; i++) {
        i = parse_statement(tokens, i, decl);
    }

    return i;
}
