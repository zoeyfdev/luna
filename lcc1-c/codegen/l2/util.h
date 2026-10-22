#pragma once

#include "types.h"
#include "../../parse/types.h"

codegen_return* new_cgr();
char* generate_generic();
void emit_store(c_type* type_info, char* dst, char* src);
void emit_load(c_type* type_info, char* dst, char* src);
