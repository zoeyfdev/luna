#include <stdio.h>
#include <stdint.h>
#include <string.h>
#include <stdlib.h>

#include "../lcc_shared/libvector.h"
#include "../lcc_shared/libfile.h"
#include "../lcc_shared/shared.h"
#include "../lcc_shared/liberror.h"
#include "lexer/lex.h"
#include "error/error.h"
#include "parse/parse.h"

int main(int argc, char* argv[]) {
    vector* files = vec_init(sizeof(char**), 0);

    for (int i = 1; i < argc; i++) {
        char* arg = argv[i];
        if (!strcmp(arg, "-v")) {
            _DISPLAY_VERSION_INFO();
        } else {
            vec_grow(files, 1);
            ((char**) files->data)[files->next] = arg;
        }
    }

    if (files->elements < 1) {
        lcc_error(NULL, 0, "no input files", NULL);
        exit(1);
    }

    for (int i = 0; i < files->elements; i++) {
        char* file = ((char**) files->data)[i];
        FILE* f = fopen(file, "rb");
        if (f == NULL) {
            lcc_error(NULL, 0, "cannot open output file '", f, "': no such file or directory", NULL);
            continue;
        }
 
        char* base = lfn_get_base(file);
        char* out_fn = calloc(1, strlen(base) + 4);
        
        strcat(out_fn, base);
        strcat(out_fn, ".o");

        fseek(f, 0, SEEK_END);
        uint64_t size = ftell(f);
        fseek(f, 0, SEEK_SET);

        char* buffer = malloc(size);
        fread(buffer, size, size, f);

        fclose(f);

        parse_init();
        vector* buf = parse(lex(buffer, size));

        if (num_errors > 0)
            exit(1);

        FILE* of = fopen(out_fn, "w+b");
        if (of == NULL) {
            lcc_error(NULL, 0, "cannot open output file '", out_fn, "'", NULL);
            continue;
        }

        fwrite(buf->data, sizeof(unsigned char), buf->elements, of);
        fclose(of);

        free(buffer);
        free(buf->data);
        free(buf);
        free(out_fn);
    }
}
