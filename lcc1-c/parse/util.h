#pragma once

#include "ast.h"
#include "../lex/types.h"

ast_node* new_node(ast_node_type type);
bool is_type(token* t);
token* peek(vector* tokens, int* i, int by);
token* expect(vector* tokens, int* i, token_type t);
