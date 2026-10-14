#include <stdio.h>
#include <string.h>

#include "../lcc_shared/shared.h"
#include "../lcc_shared/libcpp.h"
#include "../lcc_shared/libfile.h"
#include "lex/lex.h"
#include "parse/parse.h"
#include "parse/ast.h"
#include "codegen/codegen.h"

int main(int argc, char* argv[]) {
    for (int i = 1; i < argc; i++) {
        char* arg = argv[i];
        if (!strcmp(arg, "-v")) {
            _DISPLAY_VERSION_INFO();
        } else {
            char* f_contents = open_file(arg);
            ast* tu = parse(lex(cpp(arg, f_contents, "+-*/%&|^~<>=!?:;.,()[]{}")));
            char* str = codegen_l2(tu, false);
            printf("%s\n", str); 
        }
    }
    
    return 0;
} 
