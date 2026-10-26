#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <stdint.h>

#include "libvector.h"
#include "liberror.h"

#define INSERT_TOKEN(VALUE) do { \
    token* t = malloc(sizeof(token)); \
    t->value = (VALUE); \
    t->filename = file; \
    t->line = line; \
    vec_grow(tokens, 1); \
    ((token**) tokens->data)[tokens->next] = t; \
} while (0)

#ifndef _WIN32
    char* include_path = "/usr/local/include/lcc/";
#else
    char* include_path = "C:\\Program Files\\Luna L2\\include\\";
#endif

typedef struct {
    char* value;
    char* filename;
    int line;
} token;

typedef struct {
    char* original;
    vector* replacement_list;
} define;

bool is_space(char c) {
    if (c == '\t' || c == '\n' || c == '\v' || c == '\f'
        || c == '\r' || c == ' ') {
        return true;
    }
    return false;
}

bool contains(char* match, char c) {
    for (size_t i = 0; i < strlen(match); i++)
        if (match[i] == c)
            return true;
    return false;
}

vector* tokenize(char* buffer, char* file, char* match_str) {
    vector* tokens = vec_init(sizeof(token**), 0);

    size_t i = 0;
    int line = 1;

    while (i < strlen(buffer)) {
        char c = buffer[i];

        if (c == '/' && buffer[i + 1] == '/') {
            i += 2;
            while (i < strlen(buffer) && buffer[i] != '\n')
                i++;
            continue;
        }

        if (c == '/' && buffer[i + 1] == '*') {
            i += 2;

            while (i < strlen(buffer)) {
                if (buffer[i] == '*' && buffer[i + 1] == '/') {
                    i += 2;
                    break;
                }
                if (buffer[i] == '\n')
                    line++;
                i++;
            }
            continue;
        }

        if (c == '\n') {
            line++;
            i++;
            continue;
        }

        if (is_space(c)) {
            i++;
            continue;
        }

        if (c == '"' || (c == '<' && !strcmp(((token**) tokens->data)[tokens->elements - 1]->value, "#include"))) {
            bool term_on_r_angle = (c == '<' && !strcmp(((token**) tokens->data)[tokens->elements - 1]->value, "#include"));
            vector* buf = vec_init(sizeof(char*), 0);

            vec_grow(buf, 1);
            ((char*) buf->data)[buf->next] = c;
            i++;

            while (i < strlen(buffer)) {
                char c = buffer[i];

                vec_grow(buf, 1);
                ((char*) buf->data)[buf->next] = c;

                if (c == '"' || (c == '>' && term_on_r_angle))
                    break;
                else if (c == '\n')
                    line++;

                i++;
            }
            
            i++;

            vec_grow(buf, 1);
            ((char*) buf->data)[buf->next] = 0;

            token* t = malloc(sizeof(token));

            t->value = (char*) buf->data;
            t->filename = file;
            t->line = line;

            vec_grow(tokens, 1);
            ((token**) tokens->data)[tokens->next] = t;

            free(buf);

            continue;
        }

        if (c == '=' && strlen(buffer) > i + 1 && buffer[i + 1] == '=') {
            INSERT_TOKEN("==");

            i += 2;
            continue;
        }

        if (c == '>' && strlen(buffer) > i + 1 && buffer[i + 1] == '=') {
            INSERT_TOKEN(">=");

            i += 2;
            continue;
        }

        if (c == '<' && strlen(buffer) > i + 1 && buffer[i + 1] == '=') {
            INSERT_TOKEN("<=");

            i += 2;
            continue;
        }

        if (c == '!' && strlen(buffer) > i + 1 && buffer[i + 1] == '=') {
            INSERT_TOKEN("!=");

            i += 2;
            continue;
        }

        if (c == '+' && strlen(buffer) > i + 1 && buffer[i + 1] == '+') {
            INSERT_TOKEN("++");

            i += 2;
            continue;
        }

        if (c == '-' && strlen(buffer) > i + 1 && buffer[i + 1] == '-') {
            INSERT_TOKEN("--");

            i += 2;
            continue;
        }

        if (c == '<' && strlen(buffer) > i + 1 && buffer[i + 1] == '<') {
            INSERT_TOKEN("<<");

            i += 2;
            continue;
        }

        if (c == '>' && strlen(buffer) > i + 1 && buffer[i + 1] == '>') {
            INSERT_TOKEN(">>");

            i += 2;
            continue;
        }

        if (c == '-' && strlen(buffer) > i + 1 && buffer[i + 1] == '>') {
            INSERT_TOKEN("->");

            i += 2;
            continue;
        }

        if (c == '.' && strlen(buffer) > i + 2 && buffer[i + 1] == '.'
                && buffer[i + 2] == '.') {
            INSERT_TOKEN("...");

            i += 3;
            continue;
        }

        if (contains(match_str, c)) {
            char* c_loc = calloc(1, 2);
            c_loc[0] = c;
            INSERT_TOKEN(c_loc);
            
            i++;
            continue;
        }

        vector* buf = vec_init(sizeof(char*), 0);
        while (i < strlen(buffer)) {
            char c = buffer[i];

            if (is_space(c) || contains(match_str, c) || contains("\"", c))
                break;
            
            vec_grow(buf, 1);
            ((char*) buf->data)[buf->next] = c;
            i++;
        }

        if (buf->elements > 0) {
            vec_grow(buf, 1);
            INSERT_TOKEN(buf->data);
        }
    }

    return tokens;
}

vector* cpp(char* filename, char* buffer, char* match_str) { 
    vector* defines = vec_init(sizeof(define**), 0); 
    bool again = false;
    vector* tokens = tokenize(buffer, filename, match_str);
top:
    vector* output = vec_init(sizeof(token**), 0);
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
                    ((token**) output->data)[output->next] = ((token**) tokens->data)[i];
                }
                i++;
            }

            if (!strcmp(((token**) tokens->data)[i]->value, "#else")) {
                i++;
                while (strcmp(((token**) tokens->data)[i]->value, "#endif")) {
                    if ((found && is_rev) || (!found && !is_rev)) {
                        vec_grow(output, 1);
                        ((token**) output->data)[output->next] = ((token**) tokens->data)[i];
                    }
                    i++;
                }
            }
        } else if (!strcmp(value, "#include")) {
            again = true;
            token* path = ((token**) tokens->data)[++i];

            if (strlen(path->value) < 2) {
                lcc_error(filename, path->line, NULL, "invalid filename '", path->value, "'", NULL);
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
                lcc_error(filename, path->line, NULL, "invalid string", NULL);
                continue;
            }

            if (_value[strlen(_value) - 1] != term) {
                lcc_error(filename, path->line, NULL, "invalid string", NULL);
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
                    lcc_fatal(filename, path->line, "could not open file '", actual, "'", NULL);
                } else {
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

            vector* ftokens = tokenize(fbuf, actual, match_str);

            for (int i = 0; i < ftokens->elements; i++) {
                vec_grow(output, 1);
                ((token**) output->data)[output->next] = ((token**) ftokens->data)[i];
            }
        } else {
            bool found = false;
            for (int k = 0; k < defines->elements; k++) {
                define* d = ((define**) defines->data)[k];

                if (!strcmp(value, d->original)) {
                    found = true;
                    again = true;
                    
                    for (int j = 0; j < d->replacement_list->elements; j++) {
                        token* __t = ((token**) d->replacement_list->data)[j];
                        token* _t = malloc(sizeof(token));
                        memcpy(_t, __t, sizeof(token));
                        
                        _t->line = t->line;

                        vec_grow(output, 1);
                        ((token**) output->data)[output->next] = _t; 
                    }
                }
            }
    
            if (!found) {
                vec_grow(output, 1);
                ((token**) output->data)[output->next] = t;
            }
        }

        i++;
    }

    if (again) {
        again = false;
        tokens = output;
        goto top;
    }

    return output;
}
