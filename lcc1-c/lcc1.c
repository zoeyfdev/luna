#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "../lcc_shared/shared.h"
#include "../lcc_shared/libcpp.h"
#include "../lcc_shared/libfile.h"
#include "../lcc_shared/liberror.h"
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
            char* base = lfn_get_base(arg);
            char* out_fn = calloc(1, strlen(base) + 4);

            if (f_contents == NULL) {
                lcc_fatal(NULL, 0, "could not open file '", arg, "': no such file or directory", NULL);
                exit(1);
            }

            parse_init();
            ast* tu = parse(lex(cpp(arg, f_contents, "+-*/%&|^~<>=!?:;.,()[]{}")));
            char* buf = codegen_l2(tu, false);

            strcat(out_fn, base);
            strcat(out_fn, ".s");
            write_file(out_fn, (unsigned char*) buf, strlen(buf));
        }
    }
    
    return 0;
} 
