#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#include "libvector.h"

#define INSERT_TOKEN() do { token* t = malloc(sizeof(token)); \
    t->value = (char*) current->data; \
    t->line = line; \
    vec_grow(tokens, 1); \
    ((token**) tokens->data)[tokens->next] = t; \
    free(current); \
    current = vec_init(sizeof(char*), 0); \
} while(0)

#define INSERT_CHAR() do { vec_grow(current, 1); ((char*) current->data)[current->next] = c; } while(0)

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

    for (size_t i = 0; i < strlen(buffer); i++) {
        char c = buffer[i];

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
            }

            if (c == 0x0a)
                line++;

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

void c_preprocessor(char* buffer) { 
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

            define* d = malloc(sizeof(define));
            d->original = orig->value;
            d->replacement_list = list;

            vec_grow(defines, 1);
            ((define**) defines->data)[defines->next] = d;
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

    char* buf = calloc(1, len + 1);

    for (int i = 0; i < output->elements; i++) {
        strcat(buf, ((char**) output->data)[i]);
        if (buf[strlen(buf) - 1] != '\n' && i != output->next)
            strcat(buf, " ");
    }
    buf[strlen(buf)] = 0; 

    if (again) {
        again = false;
        buffer = buf;
        goto top;
    }

    printf("Final: %s\n", buf);
}
