#include <stdio.h>

#include "version.h"

char* files[];

int main(int argc, char* argv[]) {
    for (int i = 1; i < argc; i++) {
        char* arg = argv[i];
        if (strcmp(arg, "-v")) {
            printf("%s\n", "Luna ld " VERSION);
        } 
    }
}
