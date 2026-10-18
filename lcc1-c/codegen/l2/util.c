#include <stdlib.h>
#include <stdio.h>

#include "types.h"

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
