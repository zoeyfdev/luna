#include "../../lcc_shared/libvector.h"

#define TYPE_REGISTER 1
#define TYPE_TEXT 2

typedef struct {
    int type;
    char* value;
} token;

#define NUM_REGISTERS 35
char* registers[] = {
    "r0",
    "r1",
    "r2",
    "r3",
    "r4",
    "r5",
    "r6",
    "r7",
    "r8",
    "r9",
    "r10",
    "r11",
    "r12",
    "e0",
    "e1",
    "e2",
    "e3",
    "e4",
    "e5",
    "e6",
    "e7",
    "e8",
    "e9",
    "e10",
    "e11",
    "e12",
    "e13",
    "e14",
    "sp",
    "pc",
    "irv",
    "ir",
    "b",
    "fp",
    "s"
};

bool is_register(char* str) {
    for (int i = 0; i < NUM_REGISTERS; i++)
        if (!strcmp(str, registers[i]))
            return true;
    return false;
}

void lex(char* buffer) {
    vector* tokens = vec_init(sizeof(token**), 0);
    vector* current = vec_init(sizeof(char*), 0);
    bool comment = false;

    for (; *buffer; buffer++) {
        char c = *buffer;

        if ((c == '/' && *(buffer + 1) == '/') || (c == ';') || (c == '#' && *(buffer + 1) == ' ')) {
            if (c != ';')
                buffer++;
            comment = true;
        } else if (c == '\n') {
            if (comment) {
                comment = false;
                continue;
            }
        }

        if (comment)
            continue;

        switch (c) {
        case ' ', '\n': {
                token* t = malloc(sizeof(token));
                t->type = is_register((char*) current->data) ? 1 : 0;
                t->value = (char*) current->data;

                vec_grow(tokens, 1);
                ((token**) tokens->data)[tokens->next] = t;
                current = vec_init(sizeof(char*), 0); // leave vector allocated so we don't segfault
                break;
            }
        case '\r':
            break; // don't process carriage returns
        default:
            vec_grow(current, 1);
            ((char*) current->data)[current->next] = c;
            break;
        }
    }
}

