#include <stdio.h>
#include <pthread.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>

#include "video/video.h"
#include "cpu/init.h"
#include "cpu/registers.h"
#include "bios/bios.h"
#include "io/keyboard.h"
#include "hwinit/audio.h"
#include "hwinit/pit.h"
#include "hwinit/keyboard.h"

int main(int argc, char* argv[]) {
    for (int i = 1; i < argc; i++) {
        char* arg = argv[i];
        HDD_FILE = arg;
    }

    // Initialize IO devices
    // Initialize audio
    initialize_audio();
    initialize_pit();
    initialize_keyboard();

    // Initialize keyboard

    // Execute CPU
    pthread_t cpu_thread;
    pthread_create(&cpu_thread, NULL, cpu_power_on, NULL);

    initialize_window();
    return 0;
}
