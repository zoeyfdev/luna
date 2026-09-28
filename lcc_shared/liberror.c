#include <stdio.h>

void lcc_error(char* label, char* str) {
    fprintf(stderr, "\033[1;39m%s: \033[1;31merror: \033[1;39m%s\033[0m\n", label, str);
}
