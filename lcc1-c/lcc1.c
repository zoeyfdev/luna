#include <stdio.h>
#include <string.h>

#include "../lcc_shared/shared.h"
#include "../lcc_shared/libcpp.h"
#include "lex/lex.h"
#include "parse/parse.h"
#include "parse/ast.h"
#include "codegen/codegen.h"

int main(int argc, char* argv[]) {
    for (int i = 0; i < argc; i++) {
        char* arg = argv[i];
        if (!strcmp(arg, "-v")) {
            _DISPLAY_VERSION_INFO();
        }
    }

    ast* tu = parse(lex(cpp("<auto>", "asm (\"test\")", "()[]{}")));
    char* str = codegen_l2(tu, false);
    printf("%s\n", str);


    return 0;
} 
