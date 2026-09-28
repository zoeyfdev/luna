#include <stdio.h>

#include "../lcc_shared/libcommand.h"
#include "../lcc_shared/liberror.h"

bool execute_command(char* command, bool report_failure) {
    int code = lc_command_execute(command);

    printf("%s\n", command);

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
