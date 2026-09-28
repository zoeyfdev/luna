#include <stdlib.h>
#include <string.h>
#include <stdio.h>

#include "../lcc_shared/libfile.h"
#include "../lcc_shared/libvector.h"
#include "../lcc_shared/liberror.h"
#include "exec.h"

extern vector* input_files;
extern char* output_file;
extern bool no_assemble;
extern bool no_link;

vector* hl_files;
vector* as_files;
vector* ld_files;

void sort_files() {
    hl_files = vec_init(sizeof(char**), 0);
    as_files = vec_init(sizeof(char**), 0);
    ld_files = vec_init(sizeof(char**), 0);

    for (int i = 0; i < input_files->elements; i++) { 
        char* file = ((char**) input_files->data)[i];
        char* ext = lfn_get_ext(file);
        int base_length = strlen(file) - (strlen(ext) + 1);
        char* base = malloc(base_length); 
        
        memcpy(base, file, base_length);

        if (!strcmp(ext, "c") || !strcmp(ext, "h") || !strcmp(ext, "cc") || !strcmp(ext, "hh") || !strcmp(ext, "cpp") 
                || !strcmp(ext, "hpp") || !strcmp(ext, "cxx") || !strcmp(ext, "hxx")) {
            vec_grow(hl_files, 1);
            ((char**) hl_files->data)[hl_files->next] = file; 

            char* as_ver = malloc(strlen(base) + 3);
            char* ld_ver = malloc(strlen(base) + 3);

            strcpy(as_ver, base);
            strcat(as_ver, ".s");

            strcat(ld_ver, base);
            strcat(ld_ver, ".o");

            vec_grow(as_files, 1);
            ((char**) as_files->data)[as_files->next] = as_ver;

            vec_grow(ld_files, 1);
            ((char**) ld_files->data)[ld_files->next] = strcat(lfn_get_base(ld_ver), ".o");

            free(base); // should be safe since we did everything
        } else if (!strcmp(ext, "s") || !strcmp(ext, "asm")) {
            char* ld_ver = malloc(strlen(base) + 3);
            strcat(ld_ver, base);
            strcat(ld_ver, ".o");

            vec_grow(as_files, 1);
            ((char**) as_files->data)[as_files->next] = file;

            vec_grow(ld_files, 1);
            ((char**) ld_files->data)[ld_files->next] = strcat(lfn_get_base(ld_ver), ".o");

            free(base);
        } else if (!strcmp(ext, "o") || !strcmp(ext, "obj") || !strcmp(ext, "a")) {
            vec_grow(ld_files, 1);
            ((char**) ld_files->data)[ld_files->next] = file;
        } else {
            char* strbase = "unknown file type in '";
            char* str = malloc(strlen(strbase) + strlen(file) + 2);
            strcat(str, strbase);
            strcat(str, file);
            strcat(str, "'");

            lcc_error("lcc", str);

            free(str);
        }
    }    
}
