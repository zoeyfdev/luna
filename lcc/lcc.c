#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#include "../lcc_shared/shared.h"
#include "../lcc_shared/libvector.h"
#include "compile.h"
#include "../lcc_shared/liberror.h"

bool no_assemble;
bool no_link;
bool verbose;

vector* input_files;
char* output_file = "a.bin"; // default output file

int main(int argc, char* argv[]) {
    input_files = vec_init(sizeof(char*), 0);
    for (int i = 1; i < argc; i++) {
        char* arg = argv[i];

        if (!strcmp(arg, "--version")) {
            _DISPLAY_VERSION_INFO();
            exit(0);
        } else if (!strcmp(arg, "-v")) {
            _DISPLAY_VERSION_INFO();
            verbose = true;
        } else if (!strcmp(arg, "-c"))
            no_link = true;
        else if (!strcmp(arg, "-S"))
            no_assemble = true;
        else if (!strcmp(arg, "-o")) {
            if (argc > i + 1) {
                output_file = argv[i + 1];
                i++;
            } else {
                lcc_error("lcc", "argument to '-o' is missing (expected 1 value)");
            }
        } else {
            vec_grow(input_files, 1);
            ((char**) input_files->data)[input_files->next] = arg;
        }
    }

    compile_init();
    compile();
}
