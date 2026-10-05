#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <stdint.h>

#include "libvector.h"
#include "liberror.h"

#define INSERT_TOKEN() do { token* t = malloc(sizeof(token) + 5); \
    vec_grow(current, 1); \
    ((char*) current->data)[current->next] = 0; \
    printf("%s\n", (char*) current->data); \
    t->value = (char*) current->data; \
    t->line = line; \
    vec_grow(tokens, 1); \
    ((token**) tokens->data)[tokens->next] = t; \
    free(current); \
    current = vec_init(sizeof(char*), 0); \
} while(0)

#define INSERT_CHAR() do { vec_grow(current, 1); ((char*) current->data)[current->next] = c; } while(0)

#ifndef _WIN32
    char* include_path = "/usr/local/include/lcc/";
#else
    char* include_path = "C:\\Program Files\\Luna L2\\include\\";
#endif

typedef struct {
    char* value;
    int line;
} token;

typedef struct {
    char* original;
    vector* replacement_list;
} define;

vector* tokenize(char* buffer) {
    vector* current = vec_init(sizeof(char*), 0);
    vector* tokens = vec_init(sizeof(token**), 0);
    int line = 1;

    bool in_string = false;
    bool last_was_str = false;

    for (size_t i = 0; i < strlen(buffer); i++) {
        char c = buffer[i];

        if (c == '"') {
            last_was_str = true;
        } else {
            last_was_str = false;
        }

        switch (c) {
        case '"':
            in_string = !in_string;

            INSERT_CHAR();
            if (!in_string) {
                INSERT_TOKEN();
            }
            break;
        case '\r':
            break;
        case '\n':
        case ' ':
            if (!in_string) {
                INSERT_TOKEN(); 
            } else {
                INSERT_CHAR();
                break;
            }

            if (c == 0x0a) {
                INSERT_CHAR();
                INSERT_TOKEN();
                line++;
            }

            break;
        default:
            INSERT_CHAR(); 
            break;
        } 
    }

    if (current->elements > 0) {
        INSERT_TOKEN();
    }

    return tokens;
}

char* cpp(char* filename, char* buffer) { 
    vector* defines = vec_init(sizeof(define**), 0); 
    bool again = false;
top:
    vector* tokens = tokenize(buffer);
    vector* output = vec_init(sizeof(char**), 0);
    int i = 0;

    while (true) {
        if (i >= tokens->elements)
            break;

        token* t = ((token**) tokens->data)[i];
        char* value = t->value;

        if (!strcmp(value, "#define")) {
            vector* list = vec_init(sizeof(token**), 0);

            i++;

            token* orig = ((token**) tokens->data)[i];

            i++;

            while (i < tokens->elements && ((token**) tokens->data)[i]->line == t->line) {
                token* current = ((token**) tokens->data)[i];

                vec_grow(list, 1);
                ((token**) list->data)[list->next] = current;

                i++;
            }

            i--;

            define* d = malloc(sizeof(define) + 5);
            d->original = orig->value;
            d->replacement_list = list;

            vec_grow(defines, 1);
            ((define**) defines->data)[defines->next] = d;
        } else if (!strcmp(value, "#ifdef") || !strcmp(value, "#ifndef")) {
            char* label = ((token**) tokens->data)[++i]->value;
            bool found = false;
            bool is_rev = !strcmp(value, "#ifndef");

            for (int i = 0; i < defines->elements; i++) {
                define* d = ((define**) defines->data)[i];

                if (!strcmp(d->original, label)) {
                    found = true;
                    break;
                }
            }
            i++;

            while (strcmp(((token**) tokens->data)[i]->value, "#endif") && 
                    strcmp(((token**) tokens->data)[i]->value, "#else")) {
              
                if ((found && !is_rev) || (!found && is_rev)) {
                    vec_grow(output, 1);
                    ((char**) output->data)[output->next] = ((token**) tokens->data)[i]->value;
                }
                i++;
            }

            if (!strcmp(((token**) tokens->data)[i]->value, "#else")) {
                i++;
                while (strcmp(((token**) tokens->data)[i]->value, "#endif")) {
                    if ((found && is_rev) || (!found && !is_rev)) {
                        vec_grow(output, 1);
                        ((char**) output->data)[output->next] = ((token**) tokens->data)[i]->value;
                    }
                    i++;
                }
            }
        } else if (!strcmp(value, "#include")) {
            again = true;
            token* path = ((token**) tokens->data)[++i];

            if (strlen(path->value) < 2) {
                lcc_error(filename, path->line, "invalid filename '", path->value, "'", NULL);
                continue; 
            }

            char* _value = path->value;
            bool abs = false;
            char term = '"';

            if (_value[0] == '"')
                term = '"';
            else if (_value[0] == '<') {
                term = '>';
                abs = true;
            } else {
                lcc_error(filename, path->line, "invalid string", NULL);
                continue;
            }

            if (_value[strlen(_value) - 1] != term) {
                lcc_error(filename, path->line, "invalid string", NULL);
                continue;
            }
  
    f_try_top:
            int times = abs ? 1 : 0;
            char* nostrm = calloc(1, strlen(_value) + 1);
            char* actual = calloc(1, strlen(_value) + strlen(include_path) + 1); // should be enough
          
            for (size_t i = 1; i < strlen(_value) - 1; i++) {
                nostrm[i - 1] = _value[i];
            }

            if (abs)
                strcat(actual, include_path);
            strcat(actual, nostrm);
            free(nostrm);
 
            FILE* f = fopen(actual, "rb");
            if (f == NULL) {
                if (times == 1) {
                    lcc_error(filename, path->line, "could not open file '", actual, "'", NULL);
                    fprintf(stderr, "compilation terminated.\n");
                    exit(1);
                } else {
                    // free(actual);
                    abs = true;
                    goto f_try_top;
                }
            }

            fseek(f, 0, SEEK_END);
            uint64_t size = ftell(f);
            fseek(f, 0, SEEK_SET);

            char* fbuf = calloc(1, size + 5);

            fread(fbuf, sizeof(char), size, f);
            fclose(f);

            vector* ftokens = tokenize(fbuf);

            for (int i = 0; i < ftokens->elements; i++) {
                vec_grow(output, 1);
                ((char**) output->data)[output->next] = ((token**) ftokens->data)[i]->value;
            }
        } else {
            bool found = false;
            for (int k = 0; k < defines->elements; k++) {
                define* d = ((define**) defines->data)[k];

                if (!strcmp(value, d->original)) {
                    found = true;
                    again = true;
                    
                    for (int j = 0; j < d->replacement_list->elements; j++) {
                        vec_grow(output, 1);
                        ((char**) output->data)[output->next] = ((token**) d->replacement_list->data)[j]->value; 
                    }
                }
            }
    
            if (!found) {
                vec_grow(output, 1);
                ((char**) output->data)[output->next] = value;
            }
        }

        i++;
    }

    size_t len = 0;
    for (int i = 0; i < output->elements; i++) {
        len += strlen(((char**) output->data)[i]) + 2;
    }
    len++;

    char* buf = calloc(1, len + 2);

    for (int i = 0; i < output->elements; i++) {
        char* ot = ((char**) output->data)[i];

        strcat(buf, ot);
        if (buf[strlen(buf) - 1] != '\n' && i != output->next)
            strcat(buf, " ");
    }

    buf[strlen(buf) - 1] = 0;

    if (again) {
        again = false;
        buffer = buf;
        goto top;
    } 

    return buf;
}
