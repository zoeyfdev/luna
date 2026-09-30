#include <stdio.h>

void lcc_error(char* label, char* str) {
    fprintf(stderr, "\033[1;39m%s: \033[1;31merror: \033[1;39m%s\033[0m\n", label, str);
}

void lcc_info(char* label, char* str) {
    printf("\033[1;39m%s: \033[1;36minfo: \033[1;39m%s\033[0m\n", label, str);
}
