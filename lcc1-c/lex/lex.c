#include <stdlib.h>
#include <string.h>

#include "../../lcc_shared/libvector.h"

#define INSERT_CHAR() do { \
    vec_grow(current, 1); \
    ((char*) current->data)[current->next] = c; \
} while (0)

typedef enum {
    TOK_NUMBER,
    TOK_PLUS,
    TOK_MINUS,
    TOK_STAR,
    TOK_SLASH
} token_type;

typedef struct {
    token_type type;
    int line;
    char* file; 
    char* value;
    char* displayed_value;
} token;

vector* current;

void insert_token(int line, char* filename) {
    char* value = (char*) current->data;
    token t = malloc(sizeof(token));

    t->line = line;
    t->file = NULL;
    t->value = value;
    t->displayed_value = value;

    free(current);
    current = vec_init(sizeof(char*), 0);
}

vector* lex(char* buf, char* filename) {
    vector* tokens = vec_init(sizeof(token**), 0);
    current = vec_init(sizeof(char*), 0);
    int i = 0;

    bool in_string = false;
    bool in_comment = false;
    bool in_long_comment = false;
    int line = 0;
    current_filename = filename;

    size_t buf_len = strlen(buf);
    for (size_t i = 0; i < buf_len; i++) {
        char c = buf[i];

        if (in_comment)
            continue;

        switch (c) {
        case '/':
            if (buf_len > i + 1) {
                if (buf[i + 1] == '/') {
                    buf++;
                    in_comment = true;
                    break;
                } else if (buf[i + 1] == '*') {
                    buf++;
                    in_comment = true;
                    in_long_comment = true;
                    break;
                }
            }

            INSERT_CHAR();
            break;
        case '\r':
            break;
        case '\n':
            if (in_comment && !in_long_comment)
                in_comment = false;

            insert_token(line, filename);

            line++;
            break;
        default:
            INSERT_CHAR();
            break;
        }
    }

    return tokens;
}
