#pragma bits 32

#include "stdlib.h"
#include "audio.h"
#include "stdbool.h"
#include "fzip.h"

void tohex(long int number, char capitalized) {
    puts32("0x", 255, 0);
    puts32((char*) itoa(number, capitalized, (char*) malloc(11)), 255, 0);
    puts32("\n", 255, 0);
    free(11);
}

char pause() {
    puts32("Press any key to continue...\n\n", COLOR_WHITE, COLOR_BLACK);
    char code = wait_for_key();
    return code;
}

void kernel_panic() __attribute__((noreturn)) {
    asm ("push r2");
    asm ("push r1");

    play_sound((void*) fzip_decode(CRASH_SOUND), 164352, false);
    screen_fill(0x80808080);
    puts32("System error\n\nYour PC ran into an error and needs to\nbe restarted.\n\nPress any key to reboot.\n\n\n", COLOR_WHITE, COLOR_RED);
    
    puts32("Instruction: 0x", COLOR_WHITE, COLOR_RED);
    asm ("pop e9");
    puts32((char*) itoa(_e9, 1, (char*) malloc(11)), COLOR_WHITE, COLOR_RED);
    puts32("\n", COLOR_WHITE, COLOR_RED);

    puts32("Location: 0x", COLOR_WHITE, COLOR_RED);
    asm ("pop e9");
    puts32((char*) itoa(_e9, 1, (char*) malloc(11)), COLOR_WHITE, COLOR_RED);
    puts32("\n", COLOR_WHITE, COLOR_RED);

    wait_for_key();
    asm ("int 0x10");
    asm ("int 0xf");
}

void video_set_cursor(int x, int y) {
    // Arguments in e0, e1
    asm ("mov r1, e0");
    asm ("mov r2, e1");
    asm ("int 0x0c");
}

int get_cursor_x() {
    asm ("int 0xe");
    asm ("mov e6, r1");
}
int get_cursor_y() {
    asm ("int 0xe");
    asm ("mov e6, r2");
}

int cursor_x = 0;
int cursor_y = 0;
void video_save_cursor() {
    cursor_x = get_cursor_x();
    cursor_y = get_cursor_y();
}
void video_load_cursor() {
    video_set_cursor(cursor_x, cursor_y);
}

int query_drive_inserted(char drive) {
    asm ("mov r1, e0"); // Move drive number to r1
    asm ("int 0x3"); // Query drive inserted, return in r1
    asm ("mov e12, r1");
    asm ("mov e6, e12");
}

void reboot() {
    asm ("int 0x10");
    asm ("int 0xf");
}

void load_sector(char drive, long int* dest_sector, long int real_sector) {
    asm ("mov r2, e0");
    asm ("mov r1, e1");
    asm ("mov r3, e2");
    asm ("int 0x0b");
}

void load_executable() {
    if (query_drive_inserted(2) == 0) {
        puts32("Error! ", COLOR_RED, COLOR_BLACK);
        puts32("Please insert a disc into the DVD\ndrive and try again.\n", COLOR_WHITE, COLOR_BLACK);
        return;
    }

    long int* address = (long int*) ASLR_generate_address();
    address++;
   
    load_sector(2, address / 512, 0);
    load_sector(2, address / 512 + 1, 1);
    load_sector(2, address / 512 + 2, 2);

    if (*address != 0x4C325049) {   
        puts32("Error! ", COLOR_RED, COLOR_BLACK);
        puts32("Invalid executable file format.\n", COLOR_WHITE, COLOR_BLACK);
        return;
    }
    lexec_core(address);
}

void app_error() __attribute__((noreturn)) {
    puts32("Error! ", COLOR_RED, COLOR_BLACK);
    puts32("Executable automatically\nterminated due to instruction fault.\n", COLOR_WHITE, COLOR_BLACK);
    goto lexec_done; 
}

char* get_word(char* string, int pos) {
    char* buffer = (char*) malloc(1024);
    char* ogbuf = buffer;
    int cpos = 1;

    while (*string != 0x00) {
        if (*string == 0x20) {
            if (cpos == pos) {
                break;
            } else {
                cpos++;
                string++;
            }
        }

        if (cpos == pos) {
            *buffer = *string;
            buffer++;
        }

        string++;
    }
    free(1024);
    return ogbuf;
}

char* atoi(long int n) {
    char* buf = (char*) malloc(32);
    char* buf2 = (char*) malloc(32);
    char* ogbufptr = buf2;
    
    long int i = 0;
    
    if (n == 0) {
        *buf2 = 0x30;
        return ogbufptr;
    }

    while (n > 0) {
        long int res = n % 10;
        *buf = 0x30 + res;
        n = n / 10;
        i++;
        buf++;
    }

    while (i > 0) { 
        i--;
        buf--;
        *buf2 = *buf;
        buf2++; 
    }

    *buf2 = 0;
    free(64);
    return ogbufptr;
}

void toint(long int n) {
    puts32((char*) atoi(n), COLOR_WHITE, COLOR_BLACK);
    puts32("\n", COLOR_WHITE, COLOR_BLACK);
}

void printf(char* str) {
    puts32(str, COLOR_WHITE, COLOR_BLACK);
}
