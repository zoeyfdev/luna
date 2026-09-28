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

bool hl_error;
bool as_error;
bool ld_error;

void compile_init() {
    hl_files = vec_init(sizeof(char**), 0);
    as_files = vec_init(sizeof(char**), 0);
    ld_files = vec_init(sizeof(char**), 0);
}

void compile() {
    for (int i = 0; i < input_files->elements; i++) { 
        char* file = ((char**) input_files->data)[i];
        char* ext = lfn_get_ext(file);
        int base_length = strlen(file) - (strlen(ext) + 1);
        char* base = malloc(base_length); 
        
        memcpy(base, file, base_length);

        if (!strcmp(ext, "c") || !strcmp(ext, "h") || !strcmp(ext, "cc") || !strcmp(ext, "hh") || !strcmp(ext, "cpp") 
                || !strcmp(ext, "hpp")) {
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
            ((char**) ld_files->data)[ld_files->next] = ld_ver;

            free(base); // should be safe since we did everything
        } else if (!strcmp(ext, "s") || !strcmp(ext, "asm")) {
            char* ld_ver = malloc(strlen(base) + 3);
            strcat(ld_ver, base);
            strcat(ld_ver, ".o");

            vec_grow(as_files, 1);
            ((char**) as_files->data)[as_files->next] = file;

            vec_grow(ld_files, 1);
            ((char**) ld_files->data)[ld_files->next] = ld_ver;

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


    for (int i = 0; i < hl_files->elements; i++) {
        char* file = ((char**) hl_files->data)[i];
        char* command_base = "lcc1 -S ";
        int len = strlen(command_base) + strlen(file) + 1;
        char* command = malloc(len);
        memset(command, 0x00, len);

        strcat(command, command_base);
        strcat(command, file);

        bool success = execute_command(command, false);
        if (!success) hl_error = true;
    }

    if (hl_error)
        exit(1);

    if (no_assemble)
        exit(0);

    for (int i = 0; i < as_files->elements; i++) {
        char* file = ((char**) as_files->data)[i];
        char* command_base = "las -c ";
        int len = strlen(command_base) + strlen(file) + 1;
        char* command = malloc(len);
        memset(command, 0x00, len);

        strcat(command, command_base);
        strcat(command, file);

        bool success = execute_command(command, false);
        if (!success) as_error = true;
    }

    if (as_error)
        exit(1);

    if (no_link)
        exit(0);

    char* ld_command_base = "l2ld -o ";

    char* ld_command = malloc(strlen(ld_command_base) + 1);
    strcpy(ld_command, ld_command_base);

    ld_command = realloc(ld_command, strlen(ld_command) + strlen(output_file) + 1);
    strcat(ld_command, output_file);
    strcat(ld_command, " ");

    for (int i = 0; i < ld_files->elements; i++) {
        char* file = ((char**) ld_files->data)[i];
        int len = strlen(file) + 1;
        
        ld_command = realloc(ld_command, strlen(ld_command) + len + 1);
        strcat(ld_command, file);
        strcat(ld_command, " "); 
    }

    ld_error = !execute_command(ld_command, false);

    if (ld_error)
        exit(1);
}
