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
extern vector* cleanup_files;

vector* input_files;
char* output_file = "a.bin"; // default output file name

void cleanup(vector* cleanup_files) {
    for (int i = 0; i < cleanup_files->elements; i++) {
        vector* command = vec_init(sizeof(char**), 0);

        vec_grow(command, 1);
        #ifndef _WIN32
        ((char**) command->data)[command->next] = "rm";
        #else
        ((char**) command->data)[command->next] = "del";
        #endif

        vec_grow(command, 1);
        #ifndef _WIN32
        ((char**) command->data)[command->next] = "-f";
        #else
        ((char**) command->data)[command->next] = "/f";
        #endif

        vec_grow(command, 1);
        ((char**) command->data)[command->next] = ((char**) cleanup_files->data)[i];

        execute_command(command, false);
    }
}

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
                lcc_error(NULL, 0, "argument to '-o' is missing (expected 1 value)", NULL);
            }
        } else if (!strcmp(arg, "-v"))
            verbose = true;
        else {
            FILE* f = fopen(arg, "rb");
            if (f != NULL) {
                vec_grow(input_files, 1);
                ((char**) input_files->data)[input_files->next] = arg;
                fclose(f);
            } else {
                lcc_error(NULL, 0, "could not stat ", arg, ": no such file or directory", NULL);
            } 
        }
    }

    if (input_files->elements < 1) {
        lcc_error(NULL, 0, "no input files", NULL);
        exit(1);
    }

    sort_files();

    for (int i = 0; i < hl_files->elements; i++) {
        vector* command = vec_init(sizeof(char**), 0);

        vec_grow(command, 1);
        ((char**) command->data)[command->next] = "lcc1";

        vec_grow(command, 1);
        ((char**) command->data)[command->next] = "-S";

        vec_grow(command, 1);
        ((char**) command->data)[command->next] = ((char**) hl_files->data)[i];

        bool success = execute_command(command, false);
        if (!success) hl_error = true;
    }

    if (hl_error) {
        cleanup(cleanup_files);
        exit(1);
    }    

    if (no_assemble)
        exit(0);

    for (int i = 0; i < as_files->elements; i++) {
        vector* command = vec_init(sizeof(char**), 0);

        vec_grow(command, 1);
        ((char**) command->data)[command->next] = "las";

        vec_grow(command, 1);
        ((char**) command->data)[command->next] = "-c";

        vec_grow(command, 1);
        ((char**) command->data)[command->next] = ((char**) as_files->data)[i];

        bool success = execute_command(command, false);
        if (!success) as_error = true;
    }

    if (as_error) {
        cleanup(cleanup_files);
        exit(1);
    }

    if (no_link)
        exit(0);



    vector* command = vec_init(sizeof(char**), 0);

    vec_grow(command, 1);
    ((char**) command->data)[command->next] = "l2ld";

    vec_grow(command, 1);
    ((char**) command->data)[command->next] = "-a";

    vec_grow(command, 1);
    ((char**) command->data)[command->next] = "-o";

    vec_grow(command, 1);
    ((char**) command->data)[command->next] = output_file;

    for (int i = 0; i < ld_files->elements; i++) {
        vec_grow(command, 1);
        ((char**) command->data)[command->next] = ((char**) ld_files->data)[i]; 
    }

    ld_error = !execute_command(command, false);

    if (ld_error) {
        cleanup(cleanup_files);
        exit(1);
    }

    cleanup(cleanup_files);
}
