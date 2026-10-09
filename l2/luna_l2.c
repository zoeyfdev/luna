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
#include "hwinit/video.h"
#include "hwinit/audio.h"
#include "hwinit/pit.h"
#include "hwinit/keyboard.h"

int main(int argc, char* argv[]) { 
    char* video_name = "g1x";
    char* audio_name = "s1";
    char* pit_name = "pit";

    for (int i = 1; i < argc; i++) {
        char* arg = argv[i];
        if (!strcmp("-hdd", arg) && i + 1 < argc)
            HDD_FILE = argv[i + 1];
        else if (!strcmp("-sd", arg) && i + 1 < argc)
            SD_FILE = argv[i + 1];
        else if (!strcmp("-dvd", arg) && i + 1 < argc)
            DVD_FILE = argv[i + 1];
        else if (!strcmp("-video", arg) && i + 1 < argc)
            video_name = argv[i + 1];
        else if (!strcmp("-audio", arg) && i + 1 < argc)
            audio_name = argv[i + 1];
        else if (!strcmp("-pit", arg) && i + 1 < argc)
            pit_name = argv[i + 1];
        else
            HDD_FILE = argv[i];
    }

    // Initialize IO devices
    initialize_video(video_name);
    initialize_audio(audio_name);
    initialize_pit(pit_name);
    initialize_keyboard();

    // Execute CPU
    pthread_t cpu_thread;
    pthread_create(&cpu_thread, NULL, cpu_power_on, NULL);

    initialize_window();
    return 0;
}
