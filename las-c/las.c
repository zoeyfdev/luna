#include <stdio.h>
#include <stdint.h>
#include <string.h>
#include <stdlib.h>

#include "../lcc_shared/libvector.h"
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
        lcc_error("lcc", "no input files");
        exit(1);
    }

    for (int i = 0; i < files->elements; i++) {
        char* file = ((char**) files->data)[i];
        FILE* f = fopen(file, "rb");

        fseek(f, 0, SEEK_END);
        uint64_t size = ftell(f);
        fseek(f, 0, SEEK_SET);

        char* buffer = malloc(size);
        fread(buffer, size, size, f);

        vector* tokens = lex(buffer, size);

        for (int i = 0; i < tokens->elements; i++)
            printf("Token value: %s | register: %s | line: %d | instruction: %s\n", ((token**) tokens->data)[i]->value, 
                    is_register(((token**) tokens->data)[i]->value) ? "yes" : "no", ((token**) tokens->data)[i]->line,
                     ((token**) tokens->data)[i]->type == TYPE_INSTRUCTION ? "yes" : "no");

        vector* buf = parse(tokens);

        if (num_errors < 1) { 
            for (int i = 0; i < buf->elements; i++)
                printf("0x%02x ", ((unsigned char*) buf->data)[i]);
            printf("\n");
        }


        free(buffer);
    }
}
