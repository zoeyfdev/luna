#include <stdio.h>
#include <stdint.h>

#include "../lcc_shared/libvector.h"
#include "../lcc_shared/shared.h"
#include "../lcc_shared/liberror.h"

int main(int argc, char* argv[]) {
    vector* files = vec_init(sizeof(char**), 0);

    for (int i = 0; i < argc; i++) {
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



        free(buffer);
    }
}
