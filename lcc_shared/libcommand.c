#include <stdio.h>
#include <stdarg.h>
#include <limits.h>
#include <sys/wait.h>

int lc_command_execute(char* command) {
    FILE* pipe = popen(command, "r");

    char out[33];
    while (!feof(pipe)) { 
        if (fgets(out, 33, pipe) != NULL)
            printf("%s", out);
    }

    int code = WEXITSTATUS(pclose(pipe));

    return code;
}
