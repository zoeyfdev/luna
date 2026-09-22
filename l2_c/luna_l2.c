#include <stdio.h>
#include <pthread.h>
#include <stdlib.h>
#include <string.h>

#include "video/video.h"
#include "cpu/init.h"
#include "bios/bios.h"
#include "io/keyboard.h"

int main(int argc, char* argv[]) {
    for (int i = 0; i < argc; i++) {
        HDD_FILE = argv[i];
    }

    // Initialize IO devices
    KEYBOARD_MEMORY = malloc(1);

    // Execute CPU
    pthread_t cpu_thread;
    pthread_create(&cpu_thread, NULL, cpu_poweron, NULL);

    initialize_window(); 
    return 0;
}
