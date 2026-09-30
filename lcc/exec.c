#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "../lcc_shared/libcommand.h"
#include "../lcc_shared/liberror.h"
#include "../lcc_shared/libvector.h"
#include "../lcc_shared/shared.h"

extern bool verbose;

bool execute_command(vector* list, bool report_failure) {
    size_t s = 0;

    for (int i = 0; i < list->elements; i++) {
        s += strlen(((char**) list->data)[i]) + 1; // +1 for the space after it
    }
    s++; // null char
    
    char* command = malloc(s);
    memset(command, 0, s);

    for (int i = 0; i < list->elements; i++) {
        strcat(command, ((char**) list->data)[i]);
        strcat(command, " ");
    }

    if (verbose)
        printf("%s\n", command);

    int code = lc_command_execute(command);

    free(command);

    switch (code) {
    case 1:
        if (report_failure) {
            char msg[8192];
            sprintf(msg, "compilation command failed with exit code %d", code);
            lcc_error("lcc", msg);
        }
        return false;
        break;
    case SIGSEGV_CODE: {
            lcc_info("lcc", ICE_MESSAGE);
            return false;
            break;
        }
    }

    return true;
}
