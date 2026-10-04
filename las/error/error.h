#pragma once

#include "../lexer/lex.h"
#include "../../lcc_shared/libvector.h"

extern int num_errors;

void stargaze(vector* tokens, token* t);
