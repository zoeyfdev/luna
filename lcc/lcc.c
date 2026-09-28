#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#include "../lcc_shared/shared.h"
#include "../lcc_shared/libvector.h"
#include "sort.h"
#include "exec.h"
#include "../lcc_shared/liberror.h"

bool no_assemble;
bool no_link;
bool verbose;
bool hl_error;
bool as_error;
bool ld_error;

extern vector* hl_files;
extern vector* as_files;
extern vector* ld_files;

vector* input_files;
char* output_file = "a.bin"; // default output file name

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
        } else if (!strcmp(arg, "-v"))
            verbose = true;
        else {
            FILE* f = fopen(arg, "rb");
            if (f != NULL) {
                vec_grow(input_files, 1);
                ((char**) input_files->data)[input_files->next] = arg;
            } else {
                char* str = malloc(strlen(arg) + 8192);
                sprintf(str, "could not stat '%s': no such file or directory", arg);
                lcc_error("lcc", str);
                free(str);
            }
            fclose(f);
        }
    }

    if (input_files->elements < 1) {
        lcc_error("lcc", "no input files");
        exit(1);
    }

    sort_files();

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
