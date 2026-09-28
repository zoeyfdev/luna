#include <stdio.h>

#include "../lcc_shared/libcommand.h"
#include "../lcc_shared/liberror.h"

extern bool verbose;

bool execute_command(char* command, bool report_failure) {
    if (verbose)
        printf("%s\n", command);

    int code = lc_command_execute(command);

    switch (code) {
    case 1:
        return false;
        break;
    case 2: {
            if (report_failure) {
                char msg[8192];
                sprintf(msg, "compilation command failed with exit code %d", code);
                lcc_error("lcc", msg);
            }
            return false;
        }
    }

    return true;
}
