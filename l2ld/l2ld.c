#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdcountof.h>
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

size_t nbuffer;
unsigned char* buffer;

void ld_link(file* f) {
    uint64_t size = f->size;
    unsigned char* data = f->data;
    for (uint64_t i = 0; i < size; i++) {
        unsigned char* current = data + i;
        if (!memcmp(current, "LD16_", 5) || !memcmp(current, "LD32_", 5)) {
            binding* decl = malloc(sizeof(binding));

            decl->is_32 = !memcmp(data, "LD32_", 5);
            decl->location = nbuffer;

            uint64_t j = i + 5;

            size_t nsize = 0;
            decl->name = calloc(0, sizeof(char));
            while (j < size && data[j]) {
                decl->name = bump_arr(decl->name, nsize, sizeof(char));
                decl->name[nsize++] = data[j];
                j++;
            }
            j++;

            printf("Added binding %s\n", decl->name);

            bindings = bump_arr(bindings, nbindings, sizeof(bindings));
            bindings[nbindings++] = decl;

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


            binding* b = find_binding(sym_name);
            if (b == NULL) {
                unresolved_binding* ub = malloc(sizeof(unresolved_binding));
                ub->name = sym_name;
            } else {
                WRITE_VALUE_16(b->location);
                free(sym_name);
            }

            i = j - 1;
        } else if (!memcmp(current, "L_16BIT", 7) || !memcmp(current, "L_32BIT", 7)) {
            i += 6;
            is_32 = !memcmp(current, "L_32BIT", 7);
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
            files = reallocarray(files, current_arr_size + 1, sizeof(file*));
            array_size = current_arr_size + 1;
            files[current_arr_size] = f;

            fclose(f_real);
        }
    }

    if (array_size < 1) {
        fprintf(stderr, "l2ld: no input files\n");
        exit(1);
    }

    for (int i = 0; i < array_size; i++) {
        printf("doing file\n");
        file* file = files[i];
        ld_link(file);
    }

    if (output_file == NULL)
        output_file = "a.o";

    FILE* out_file = fopen(output_file, "w+b");
    if (out_file == NULL) {
        fprintf(stderr, "l2ld: could not write '%s'\n", output_file);
        exit(1);
    }
    fwrite(buffer, sizeof(unsigned char), nbuffer, out_file);
    fclose(out_file);

    printf("%s", output_file);
}
