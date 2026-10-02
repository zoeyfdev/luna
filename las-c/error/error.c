#include <stdio.h>
#include <string.h>

#include "../lexer/lex.h"
#include "../../lcc_shared/libvector.h"

int num_errors;

void stargaze(vector* tokens, token* t) {
    printf("    %d | ", t->line);
    for (int i = 0; i < tokens->elements; i++) {
        if (((token**) tokens->data)[i]->line != t->line)
            break;

        if (!memcmp(t, ((token**) tokens->data)[i], sizeof(token)))
            printf("\033[1;31m%s \033[1;39m", ((token**) tokens->data)[i]->value); 
        else
            printf("%s ", ((token**) tokens->data)[i]->value);
    }
    printf("\n\033[0m");
}
