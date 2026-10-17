#pragma once

#include "ast.h"
#include "../../lcc_shared/libvector.h"

int parse_local(vector* tokens, int start, ast_node* decl);
