#pragma once

#include "ast.h"
#include "../lex/types.h"
#include "../../lcc_shared/libvector.h"

ast_node* new_node(ast_node_type type);
bool is_type(token* t);
token* peek(vector* tokens, int* i, int by);
token* expect(vector* tokens, int* i, token_type t);
bool is_qualifier(token* t);
void add_node(ast_node* p, ast_node* v);
