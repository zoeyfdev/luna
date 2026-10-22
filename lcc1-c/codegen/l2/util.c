#include <stdlib.h>
#include <stdio.h>

#include "types.h"
#include "codegen.h"
#include "../../parse/types.h"

codegen_return* new_cgr() {
    codegen_return* cgr = calloc(1, sizeof(codegen_return));

    return cgr;
}

int generic_ticker = 0;
char* generate_generic() {
    char* name = calloc(1, 512);
    sprintf(name, "lcc_compile_time_%d", generic_ticker++);
    return name;
}

void emit_store(c_type* type_info, char* dst, char* src) {
    if (type_info->pointer_length <= 0) {
        if (type_info->type == TYPE_I8) {
            write("str ", true); 
        } else if (type_info->type == TYPE_I16) {
            write("str16 ", true);
        } else if (type_info->type == TYPE_I32) {
            write("str32 ", true);
        } else if (type_info->type == TYPE_I64) {
            // unimpl msg
        }
    } else {
        write("str_ptr ", true);  
    }

    write(dst, false);
    write(", ", false);
    write(src, false);
    write("\n", false);
}

void emit_load(c_type* type_info, char* dst, char* src) {
    if (type_info->pointer_length <= 0) {
        if (type_info->type == TYPE_I8) {
            write("lod ", true); 
        } else if (type_info->type == TYPE_I16) {
            write("lod16 ", true);
        } else if (type_info->type == TYPE_I32) {
            write("lod32 ", true);
        } else if (type_info->type == TYPE_I64) {
            // unimpl msg
        }
    } else {
        write("lod_ptr ", true);  
    }

    write(dst, false);
    write(", ", false);
    write(src, false);
    write("\n", false);
}
