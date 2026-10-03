#include <stdio.h>
#include <stdarg.h>

void lcc_error(char* label, int line, ...) {
    if (label != NULL)
        fprintf(stderr, "\033[1;39m%s:\033[0m", label);
    else
        fprintf(stderr, "\033[1;39m%s:\033[0m", "lcc");

    if (line > 0)
        fprintf(stderr, "\033[1;39m%d:\033[0m ", line);
    else
        fprintf(stderr, " ");

    fprintf(stderr, "\033[1;31merror:\033[0m ");

    va_list strings;
    va_start(strings, line);

    char* s;
    while ((s = ((char*) va_arg(strings, char*))) && (s != NULL))
        fprintf(stderr, "\033[1;39m%s\033[0m", s);

    fprintf(stderr, "\n");
}

void lcc_info(char* label, char* str) {
    printf("\033[1;39m%s: \033[1;36minfo: \033[1;39m%s\033[0m\n", label, str);
}
