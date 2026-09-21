#include <stdio.h>
#include <pthread.h>
#include <stdlib.h>

#include "video/video.h"
#include "cpu/cpu.h"

int main() {
    // Execute CPU
    pthread_t cpu_thread;
    pthread_create(&cpu_thread, NULL, cpu_init, NULL);
    pthread_join(cpu_thread, NULL);

    InitializeWindow();

    return 0;
}
