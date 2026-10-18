#include "ast.h"
#include "types.h"
#include "../lex/types.h"
#include "../../lcc_shared/libvector.h"
#include "util.h"

#define PEEK(by) (peek(tokens, &i, (by)))
#define EXPECT(t) (expect(tokens, &i, (t)))

int parse_leaf(vector* tokens, int i, ast_node* decl, ast_node* r) {
    if (PEEK(0)->type == TOK_IDENTIFIER) {

    } else if (PEEK(0)->type == TOK_NUMBER) {
        token* int_token = EXPECT(TOK_NUMBER);
        ast_node* num_lit = new_node(AST_INTEGER_LITERAL);

        num_lit->content.integer_literal.value = int_token->value;

        r = num_lit;
    } else if (PEEK(0)->type == TOK_STRING) {
        token* str_token = EXPECT(TOK_STRING);
        ast_node* num_lit = new_node(AST_STRING_LITERAL);

        str_lit->content.string_literal.value = str_token->value;

        r = str_lit;
    }

    return i;
}

int parse_1(vector* tokens, int i, ast_node* decl, ast_node* r) {
    i = parse_leaf(tokens, i, decl, r);

    return i;
}

int parse_2(vector* tokens, int i, ast_node* decl, ast_node* r) {
    i = parse_1(tokens, i, decl, r);

    return i;
}

int parse_3(vector* tokens, int i, ast_node* decl, ast_node* r) { // *, /, %
    i = parse_2(tokens, i, decl, r);

    return i;
}

int parse_4(vector* tokens, int i, ast_node* decl, ast_node* r) { // + and -
    i = parse_3(tokens, i, decl, r);

    return i;
}

int parse_5(vector* tokens, int i, ast_node* decl, ast_node* r) { // << and >>
    i = parse_4(tokens, i, decl, r);

    return i;
}

int parse_6(vector* tokens, int i, ast_node* decl, ast_node* r) { // <, <=, >, >=
    i = parse_5(tokens, i, decl, r);

    return i;
}

int parse_7(vector* tokens, int i, ast_node* decl, ast_node* r) { // == and !=
    i = parse_6(tokens, i, decl, r);

    return i;
}

int parse_8(vector* tokens, int i, ast_node* decl, ast_node* r) { // bitwise AND
    i = parse_7(tokens, i, decl, r);

    return i;
}

int parse_9(vector* tokens, int i, ast_node* decl, ast_node* r) { // bitwise AND
    i = parse_8(tokens, i, decl, r);

    return i;
}

int parse_10(vector* tokens, int i, ast_node* decl, ast_node* r) { // bitwise OR
    i = parse_9(tokens, i, decl, r);

    return i;
}

int parse_11(vector* tokens, int i, ast_node* decl, ast_node* r) { // logical AND
    i = parse_10(tokens, i, decl, r);

    return i;
}

int parse_12(vector* tokens, int i, ast_node* decl, ast_node* r) { // logical OR
    i = parse_11(tokens, i, decl, r);

    return i;
}

int parse_13(vector* tokens, int i, ast_node* decl, ast_node* r) { // ternary
    i = parse_12(tokens, i, decl, r);

    return i;
}

int parse_14(vector* tokens, int i, ast_node* decl, ast_node* r) { // simple and compound assignments
    ast_node* left = r;
P14_TOP:
    i = parse_13(tokens, i, decl, left);

    if (PEEK(0)->type == TOK_EQUAL) {
        EXPECT(TOK_EQUAL);
        ast_node* right = parse_13(tokens, i, decl, right);

        ast_node* assignment = new_node(AST_ASSIGNMENT);
        assignment->content.assignment.target = left;
        assignment->content.assignment.value = right;
        left = assignment;
        goto P14_TOP;
    }

    r = left;

    return i;
}

int parse_15(vector* tokens, int i, ast_node* decl, ast_node* r) { // ,
    i = parse_14(tokens, i, decl, r);

    return i;
}

int parse_expression(vector* tokens, int i, ast_node* decl, ast_node* r) { // top level wrapper
    i = parse_15(tokens, i, decl, r);

    return i;
}

int parse_statement(vector* tokens, int i, ast_node* decl) {
    token* t = ((token**) tokens->data)[i];

    if (t->type == TOK_ASM) {
        EXPECT(TOK_ASM);
        EXPECT(TOK_L_PAREN);
        token* s = EXPECT(TOK_STRING);

        ast_node* _asm_node = new_node(AST_ASSEMBLY_STATEMENT);
        _asm_node->content.assembly.string = s->value;
        add_node(decl, _asm_node);

        EXPECT(TOK_R_PAREN);
    } else {
        // assume expression
        ast_node* e;
        i = parse_expression(tokens, i, decl, e);

        ast_node* stmt_expr = new_node(AST_STATEMENT_EXPRESSION);
        stmt_expr.expression = e;
    }

    EXPECT(TOK_SEMI);

    return i;
}

int parse_local(vector* tokens, int start, ast_node* decl) {
    int i = start;
    for (; i < tokens->elements; i++) {
        i = parse_statement(tokens, i, decl);
    }

    return i;
}
