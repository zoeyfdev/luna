#include <stdio.h>
#include <stdarg.h>
#include <limits.h>

#ifndef _WIN32
    #include <sys/wait.h>
#endif

int lc_command_execute(char* command) {
    FILE* pipe = popen(command, "r");

    char out[33];
    while (!feof(pipe)) { 
        if (fgets(out, 33, pipe) != NULL)
            printf("%s", out);
    }

    int status = pclose(pipe);

#ifndef _WIN32
    int code = WEXITSTATUS(status);
#else
    int code = 0;
#endif

    return code;
}
