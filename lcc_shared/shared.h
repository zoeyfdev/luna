#include <stdio.h>
#include <unistd.h>

#define VERSION_MAJOR "9"
#define VERSION_MINOR "1"
#define TARGET "luna-l2"

#define SUPPORTED_LANGUAGES "C, asm"

void _DISPLAY_VERSION_INFO() {
    char CWD[4096];
    getcwd(CWD, 4096);

    printf("%s\n", "Luna Compiler Collection version " VERSION_MAJOR "." VERSION_MINOR);
    printf("%s\n", "Target: " TARGET);
    printf("%s %s\n", "InstalledDir:", CWD);
    printf("%s\n", "Supported languages: " SUPPORTED_LANGUAGES);
}
