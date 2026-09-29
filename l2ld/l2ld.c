#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>

#include "info.h"
#include "types.h"
#include "util.h"
#include "autolink.h"
#include "../lcc_shared/libvector.h"

#define WRITE_VALUE_16(x) do { write((x) >> 8); write((x) & 0xFF); } while (0)
#define WRITE_VALUE_32(x) do { write((x) >> 24); write((x) >> 16); write((x) >> 8); write((x) & 0xFF); } while(0)

bool is_32 = false;
char* output_file;

uint16_t current_org = 0;
uint16_t padding_dir = 0;

vector* bindings;
vector* unresolved_bindings;
vector* preset_globals;
vector* buffer;

bool do_not_compile = false;

void ld_link(file* f) {
    uint64_t size = f->size;
    unsigned char* data = f->data;
    
    preset_globals = vec_init(sizeof(char*), 0);

    for (uint64_t i = 0; i < size; i++) {
        unsigned char* current = data + i;
        if (!memcmp(current, "LD16_", 5) || !memcmp(current, "LD32_", 5)) {
            binding* decl = calloc(1, sizeof(binding));

            decl->is_32 = !memcmp(current, "LD32_", 5);
            decl->location = buffer->elements + current_org;
            decl->file = f->name;
            decl->global = false;

            uint64_t j = i + 5;

            vector* name = vec_init(sizeof(char), 0);
            while (j < size && data[j]) {
                vec_grow(name, 1);
                ((char*) name->data)[name->next] = data[j];
                j++;
            }
            j++;

            decl->name = (char*) name->data;

            for (int i = 0; i < preset_globals->elements; i++) {
                if (((char**) preset_globals->data)[i] == NULL)
                    continue;
                if (!strcmp(decl->name, ((char**) preset_globals->data)[i])) {
                    decl->global = true;
                    break;
                }
            }

            if (find_binding(decl->name, f->name) != NULL) {
                fprintf(stderr, "%s:(0x%08lx): redefinition of `%s'\n", f->name, i, decl->name);
                do_not_compile = true;
            }

            vec_grow(bindings, 1);
            ((binding**) bindings->data)[bindings->next] = decl;

            cleanup_unresolved();
            i = j - 1;
        } else if (!memcmp(current, "LR_", 3)) {
            uint64_t j = i + 3;

            vector* name = vec_init(sizeof(char), 0);
            while (j < size && data[j]) {
                vec_grow(name, 1);
                ((char*) name->data)[name->next] = data[j];
                j++;
            }
            j++;

            binding* b = find_binding((char*) name->data, f->name);
            if (b == NULL) {
                unresolved_binding* ub = calloc(1, sizeof(unresolved_binding));
                ub->name = (char*) name->data;
                ub->location = buffer->elements;
                ub->file = f->name;
                ub->solved = false;
                
                vec_grow(unresolved_bindings, 1);
                ((unresolved_binding**) unresolved_bindings->data)[unresolved_bindings->next] = ub;

                if (!is_32)
                    WRITE_VALUE_16(0x00);
                else
                    WRITE_VALUE_32(0x00);
            } else {
                if (b->is_32 == true && is_32 == false) {
                    printf("%s:(0x%08lx): warning: referencing 32-bit label from 16-bit code\n", f->name, i);
                } else if (b->is_32 == false && is_32 == true) {
                    printf("%s:(0x%08lx): warning: referencing 16-bit label from 32-bit code\n", f->name, i);
                }

                if (!b->is_32)
                    WRITE_VALUE_16(b->location);
                else
                    WRITE_VALUE_32(b->location);
            }

            cleanup_unresolved();
            i = j - 1;
        } else if (!memcmp(current, "L_16BIT", 7) || !memcmp(current, "L_32BIT", 7)) {
            i += 6;
            is_32 = !memcmp(current, "L_32BIT", 7);
        } else if (!memcmp(current, "L_GLOBL_", 8)) {
            uint64_t j = i + 8;

            vector* name = vec_init(sizeof(char), 0);
            while (j < size && data[j]) {
                vec_grow(name, 1);
                ((char*) name->data)[name->next] = data[j];
                j++;
            }
            j++;

            binding* b = find_binding((char*) name->data, f->name);
            if (b != NULL) {
                b->global = true;
                cleanup_unresolved();
            } else {
                vec_grow(preset_globals, 1);
                ((char**) preset_globals->data)[preset_globals->next] = (char*) name->data;
            }

            i = j - 1;
        } else if (!memcmp(current, "LO_", 3)) {
            i += 3;
            current_org = data[i] << 8 | data[i + 1];
            i++;
        } else if (!memcmp(current, "LP_", 3)) {
            i += 3;
            padding_dir = data[i] << 8 | data[i + 1];
            i++;
        } else
            write(data[i]);
    }

    cleanup_unresolved();
}

file* load_file(char* filename) {
    FILE* f_real = fopen(filename, "rb");
    if (f_real == NULL) {
        fprintf(stderr, "l2ld: cannot find '%s': %s\n", filename, "no such file or directory");
        exit(1);
    }

    fseek(f_real, 0, SEEK_END);
    int64_t end = ftell(f_real);

    if (end == -1) {
        fprintf(stderr, "l2ld: cannot seek file '%s'\n", filename);
        exit(1);
    }
    fseek(f_real, 0, SEEK_SET);

    file* f = calloc(1, sizeof(file));

    f->name = filename;
    f->size = end;
    f->data = calloc(1, end);
    
    fread(f->data, end, end, f_real);

    
    fclose(f_real);

    return f;
}

void do_autolink() {
    bool unresolved = false;

    for (int i = 0; i < unresolved_bindings->elements; i++) {
        unresolved_binding* ub = ((unresolved_binding**) unresolved_bindings->data)[i];
        if (!ub->solved) {
            unresolved = true;
        }
    }

    if (!unresolved)
        return;

    vector* groups = autolink();

    for (int i = 0; i < groups->elements; i++) {
        autolink_pair* pair = ((autolink_pair**) groups->data)[i];

        for (int k = 0; k < unresolved_bindings->elements; k++) {
            unresolved_binding* ub = ((unresolved_binding**) unresolved_bindings->data)[k];
            if (!ub->solved) {
                if (!strcmp(ub->name, pair->label)) {
                    file* f = load_file(pair->file);
                    ld_link(f);
                    cleanup_unresolved();
                }
            }

        }
    }
}

int main(int argc, char* argv[]) {
    vector* files = vec_init(sizeof(file*), 0);
    bool allow_autolink = false;

    bindings = vec_init(sizeof(binding*), 0);
    unresolved_bindings = vec_init(sizeof(unresolved_binding*), 0);
    buffer = vec_init(sizeof(unsigned char), 0);

    for (int i = 1; i < argc; i++) {
        char* arg = argv[i];
        if (!strcmp(arg, "-v")) {
            printf("%s\n", "Luna ld (Luna Compiler Collection) " VERSION);
            exit(0);
        } else if (!strcmp(arg, "--help")) {
            printf("%s\n", HELP_STRING); 
            exit(0);
        } else if (!strcmp(arg, "-o")) {
            if (argc > i + 1) {
                output_file = argv[i + 1];
                i++;
            }
        } else if (!strcmp(arg, "-a")) {
            allow_autolink = true;
        } else {
            vec_grow(files, 1);
            ((file**) files->data)[files->next] = load_file(arg);
        }
    }

    if (files->elements < 1) {
        fprintf(stderr, "l2ld: no input files\n");
        exit(1);
    }

    for (int i = 0; i < files->elements; i++)
        ld_link(((file**) files->data)[i]);

    if (allow_autolink)
        do_autolink();

    for (int i = 0; i < unresolved_bindings->elements; i++) {
        unresolved_binding* ub = ((unresolved_binding**) unresolved_bindings->data)[i];
        if (!ub->solved) {
            do_not_compile = true;
            fprintf(stderr, "%s:(0x%08lx): undefined reference to `%s'\n", ub->file, ub->location, ub->name);
        }
    }

    if (padding_dir > 0) {
        if (buffer->elements <= padding_dir) {
            for (int i = buffer->elements; i < padding_dir; i++) {
                write(0x00);
            }
        } else {
            fprintf(stderr, "l2ld: binary exceeds padding directive: requested: %d, actual: %d\n", padding_dir, buffer->elements);
            do_not_compile = true;
        }
    }

    if (do_not_compile)
        exit(1);

    if (output_file == NULL)
        output_file = "a.bin";

    FILE* out_file = fopen(output_file, "w+b");
    if (out_file == NULL) {
        fprintf(stderr, "l2ld: could not write '%s'\n", output_file);
        exit(1);
    }

    fwrite(buffer->data, sizeof(unsigned char), buffer->elements, out_file);
    fclose(out_file);
}
