#pragma bits 32

#include "images.h"
#include "stdlib.h"
#include "lufs.h"
#include "util.h"
#include "audio.h"
#include "stdbool.h"
#include "textedit.h"

bool tried = false;
void teststack() {
    long int* rand = (long int*) 0x90909090;
    long int val = *rand;

    tohex(val, 1);
    
    if (tried == false) {
        tried = true;
        teststack();
    } else {
        return;
    }
    
    tohex(val, 1);
}

bool alloced = false;
void shell() {
    while (1) {
        if (alloced)
            free(256);

        long int* buf = malloc(256);
        alloced = true;

        puts32((char*) PROMPTBUF, COLOR_WHITE, COLOR_BLACK);
        readin((char*) buf, 1, 0);

        char* command = get_word((char*) buf, 1);

        if (strcmp("reboot", command)) {
            puts32("Rebooting...", COLOR_WHITE, COLOR_BLACK);
            asm ("mov r1, 0");
            asm ("int 0xf"); 
        }

        if (strcmp("about", command)) {
            puts32("LunaOS 2.0.0\nBy Zoey Flax\n", COLOR_WHITE, COLOR_BLACK);
            puts32("\n\n", COLOR_WHITE, COLOR_BLACK);
            continue;
        }

        if (strcmp("promptedit", command)) {
            puts32("Enter terminal prompt: ", COLOR_WHITE, COLOR_BLACK);
            readin((char*) PROMPTBUF, 0, 0);
            save_buffer((char*) PROMPTBUF, 0);

            puts32("\n", COLOR_WHITE, COLOR_BLACK);
            continue;
        }
        
        if (strcmp("open", command)) {
            char* file_name = get_word((char*) buf, 2);
            malloc(64); // Unfree memory from get_word, but only as much as we need (yes I know this is cancerous)
            
            if (strlen(file_name) == 0) {
                puts32("Usage: open <filename>\n", COLOR_WHITE, COLOR_BLACK);
                free(64); // refree memory from get_word
                continue;
            }

            long int size = fgetsize((char*) fntf(file_name)); 
            File* f = fopen((char*) fntf(file_name), true);

            long int* file = f->Address;
            if (file == NULL) {
                free(64); // refree memory from get_word
                continue;
            }

            long int* buf = malloc(size);
            char* final = strcpy((char*) file, (char*) buf);
            *final = 0;

            textedit_init((char*) buf);
            fwrite((char*) fntf(file_name), (char*) buf);

            free(64); // refree memory from get_word
            free(size);

            puts32("\n", COLOR_WHITE, COLOR_BLACK);
            continue;
        }

        if (strcmp("files", command)) {
            flist();
            puts32("\n", COLOR_WHITE, COLOR_BLACK);
            continue;
        }

        if (strcmp("shutdown", command)) {
            puts32("Shutting down...\n", COLOR_WHITE, COLOR_BLACK);
            asm ("int 0x11");
        } 

        if (strcmp("testfault", command)) { 
            asm ("mov r1, 4");
            asm ("mov r2, pc");
            asm ("int 0x07");
        }
        
        if (strcmp("clear", command)) {
            render_buf((void*) 0x40404040);
            video_set_cursor(0, 0);
            continue;
        } 

        if (strcmp("exec", command)) {
            load_executable();
            continue;
        }

        if (strcmp("battery", command)) {
            puts32("Battery level: ", COLOR_WHITE, COLOR_BLACK);
            puts32((char*) atoi(*(char*) 0x80000026), COLOR_WHITE, COLOR_BLACK);
            puts32("%\n", COLOR_WHITE, COLOR_BLACK);
            continue;
        } 

        if (strcmp("teststack", command)) {
            teststack();
            puts32("The first and third value should be the same.\n", COLOR_WHITE, COLOR_BLACK);
            continue;
        }

        if (strcmp("time", command)) { // HH:MM (24 hour)
            puts32("Time: ", COLOR_WHITE, COLOR_BLACK);
            puts32(atoi(*(char*) 0x80000022), COLOR_WHITE, COLOR_BLACK);
            puts32(":", COLOR_WHITE, COLOR_BLACK);
            puts32(atoi(*(char*) 0x80000021), COLOR_WHITE, COLOR_BLACK);
            puts32("\n", COLOR_WHITE, COLOR_BLACK);
            continue;
        }

        puts32("Bad command '", COLOR_LRED, COLOR_BLACK);
        puts32(command, COLOR_LRED, COLOR_BLACK);
        puts32("'\n", COLOR_LRED, COLOR_BLACK);
    }

    return;
}
