#include <sys/types.h>
#include <dirent.h>
#include <errno.h>
#include <stdio.h>
#include <stdint.h>

#include "types.h"
#include "../lcc_shared/libvector.h"
#include "../lcc_shared/libfile.h"

typedef struct {
    char* label;
    char* file;
} autolink_pair;

extern int nunresolved_bindings;
extern unresolved_binding** unresolved_bindings;

char* directory = "/usr/local/lib/l2ld/";

vector* autolink_file(char* filename, vector* al_groups) {
    FILE* file = fopen(filename, "rb");
    if (file == NULL) {
        fprintf(stderr, "l2ld: could not open autolink file '%s': no such file or directory\n", filename);
        return al_groups;
    }

    fseek(file, 0, SEEK_END);
    uint64_t end = ftell(file);
    fseek(file, 0, SEEK_SET);

    if (end == 0)
        return al_groups;
 
    char* file_data = malloc(end);
    fread(file_data, end, end, file);

    int cur_word = 0;

    while (true) {
        autolink_pair* pair = malloc(sizeof(autolink_pair));

        vector* word1 = get_file_word(file_data, end, cur_word);

        if (strlen((char*) word1->data) == 0)
            break;

        vector* word2 = get_file_word(file_data, end, cur_word + 1);

        pair->label = word1->data;
        pair->file = word2->data;

        vec_grow(al_groups, 1);
        ((autolink_pair**) al_groups->data)[al_groups->next] = pair;

        cur_word += 2;
    }

    fclose(file);
    return al_groups;
}

vector* autolink() {
    vector* files = vec_init(sizeof(char*), 0); 
    vector* al_groups = vec_init(sizeof(autolink_pair*), 0);

    DIR* al_dir = opendir(directory);

    if (al_dir == NULL) {
        fprintf(stderr, "l2ld: cannot open autolink directory");
        return al_groups;
    } 

    struct dirent *entry;
    while ((entry = readdir(al_dir)) != NULL) {
        if (!strcmp(entry->d_name, "..") || !strcmp(entry->d_name, ".")) 
            continue;

        if (!strcmp(lfn_get_ext(entry->d_name), "lib")) {
            vec_grow(files, 1);
            ((char**) files->data)[files->next] = entry->d_name;
        }
    }

    for (int i = 0; i < files->elements; i++) {
        char* name = malloc(4096);
        strcpy(name, directory); 
        strcpy(name + strlen(directory), ((char**) files->data)[i]);

        autolink_file(name, al_groups); 
    }

    return al_groups;
}
