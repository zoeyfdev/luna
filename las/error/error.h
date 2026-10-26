#pragma once

#include "../lexer/lex.h"
#include "../../lcc_shared/libvector.h"

extern bool errors;

void stargaze(vector* tokens, token* t);
