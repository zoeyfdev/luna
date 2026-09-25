#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>

#include "info.h"
#include "types.h"
#include "util.h"

#define WRITE_VALUE_16(x) do { write((x) >> 8); write((x) & 0xFF); } while (0)
#define WRITE_VALUE_32(x) do { write((x) >> 24); write((x) >> 16); write((x) >> 8); write((x) & 0xFF); } while(0)

bool is_32 = false;
char* output_file;
uint64_t current_org = 0;

size_t nbindings;
binding** bindings;

size_t nunresolved_bindings;
unresolved_binding** unresolved_bindings;

size_t nbuffer;
unsigned char* buffer;

bool do_not_compile = false;

void ld_link(file* f) {
    uint64_t size = f->size;
    unsigned char* data = f->data;
    for (uint64_t i = 0; i < size; i++) {
        unsigned char* current = data + i;
        if (!memcmp(current, "LD16_", 5) || !memcmp(current, "LD32_", 5)) {
            binding* decl = malloc(sizeof(binding));

            decl->is_32 = !memcmp(current, "LD32_", 5);
            decl->location = nbuffer + current_org;
            decl->file = f->name;

            uint64_t j = i + 5;

            size_t nsize = 0;
            decl->name = calloc(0, sizeof(char));
            while (j < size && data[j]) {
                decl->name = bump_arr(decl->name, nsize, sizeof(char));
                decl->name[nsize++] = data[j];
                j++;
            }
            j++;

            if (find_binding(decl->name, f->name) != NULL) {
                fprintf(stderr, "%s:(0x%08x): redefinition of `%s'\n", f->name, i, decl->name);
                do_not_compile = true;
            }

            bindings = bump_arr(bindings, nbindings, sizeof(bindings));
            bindings[nbindings++] = decl;

            for (int k = 0; k < nunresolved_bindings; k++) {
                unresolved_binding* ub = unresolved_bindings[k];
                if (!strcmp(decl->name, ub->name) && !ub->solved) { 
                    ub->solved = true;
                    if (!decl->is_32) {
                        buffer[ub->location] = decl->location >> 8;
                        buffer[ub->location + 1] = decl->location & 0xFF;
                    } else {
                        buffer[ub->location] = decl->location >> 24;
                        buffer[ub->location + 1] = decl->location >> 16;
                        buffer[ub->location + 2] = decl->location >> 8;
                        buffer[ub->location + 3] = decl->location & 0xFF;
                    }
                }
            }

            i = j - 1;
        } else if (!memcmp(current, "LR_", 3)) {
            uint64_t j = i + 3;

            size_t nsize = 0;
            char* sym_name = calloc(0, sizeof(char));
            while (j < size && data[j]) {
                sym_name = bump_arr(sym_name, nsize, sizeof(char));
                sym_name[nsize++] = data[j];
                j++;
            }
            j++;

            binding* b = find_binding(sym_name, f->name);
            if (b == NULL) {
                unresolved_binding* ub = malloc(sizeof(unresolved_binding));
                ub->name = sym_name;
                ub->location = nbuffer;
                ub->file = f->name;
                unresolved_bindings = bump_arr(unresolved_bindings, nunresolved_bindings, sizeof(unresolved_binding));
                unresolved_bindings[nunresolved_bindings++] = ub;

                if (!is_32)
                    WRITE_VALUE_16(0x00);
                else
                    WRITE_VALUE_32(0x00);
            } else {
                if (b->is_32 == true && is_32 == false) {
                    printf("%s:(0x%08x): warning: referencing 32-bit label from 16-bit code\n", f->name, i);
                } else if (b->is_32 == false && is_32 == true) {
                    printf("%s:(0x%08x): warning: referencing 16-bit label from 32-bit code\n", f->name, i);
                }
                WRITE_VALUE_16(b->location);
                free(sym_name);
            }

            i = j - 1;
        } else if (!memcmp(current, "L_16BIT", 7) || !memcmp(current, "L_32BIT", 7)) {
            i += 6;
            is_32 = !memcmp(current, "L_32BIT", 7);
        } else if (!memcmp(current, "L_GLOBL_", 8)) {
            uint64_t j = i + 8;

            size_t nsize = 0;
            char* sym_name = calloc(0, sizeof(char));
            while (j < size && data[j]) {
                sym_name = bump_arr(sym_name, nsize, sizeof(char));
                sym_name[nsize++] = data[j];
                j++;
            }
            j++;

            printf("gname: %s\n", sym_name);
            binding* b = find_binding(sym_name, f->name);
            if (b != NULL) {
                b->global = true;
            }

            free(sym_name);
            i = j - 1;
        } else
            write(data[i]);
    }
}

int main(int argc, char* argv[]) {
    file** files = calloc(0, sizeof(file));
    uint64_t array_size = 0;

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
        } else {
            FILE* f_real = fopen(arg, "rb");
            if (f_real == NULL) {
                fprintf(stderr, "l2ld: cannot find '%s': %s\n", arg, "no such file or directory");
                exit(1);
            }

            fseek(f_real, 0, SEEK_END);
            int64_t end = ftell(f_real);

            if (end == -1) {
                fprintf(stderr, "l2ld: cannot seek file '%s'\n", arg);
                exit(1);
            }
            fseek(f_real, 0, SEEK_SET);

            file* f = malloc(sizeof(file));

            f->name = arg;
            f->size = end;
            f->data = malloc(end);
            
            fread(f->data, end, end, f_real);

            uint64_t current_arr_size = array_size / sizeof(file*);
            #ifndef __APPLE__
                files = reallocarray(files, array_size + 1, sizeof(file*));
            #else
                files = realloc(files, (array_size + 1) * sizeof(file*));
            #endif
            printf("current: %d\nnew: %d\n", array_size, current_arr_size);
            files[array_size++] = f;
            fclose(f_real);
        }
    }

    if (array_size < 1) {
        fprintf(stderr, "l2ld: no input files\n");
        exit(1);
    }

    printf("%d\n", array_size);
    for (int i = 0; i < array_size; i++) {
        file* file = files[i];
        ld_link(file);
    }

    bool unresolved;
    for (int i = 0; i < nunresolved_bindings; i++) {
        unresolved_binding* ub = unresolved_bindings[i];
        if (!ub->solved) {
            unresolved = true;
            fprintf(stderr, "%s:(0x%08x): undefined reference to `%s'\n", ub->file, ub->location, ub->name);
        }
    }

    if (unresolved || do_not_compile)
        exit(1);

    if (output_file == NULL)
        output_file = "a.o";

    FILE* out_file = fopen(output_file, "w+b");
    if (out_file == NULL) {
        fprintf(stderr, "l2ld: could not write '%s'\n", output_file);
        exit(1);
    }
    fwrite(buffer, sizeof(unsigned char), nbuffer, out_file);
    fclose(out_file);
}
