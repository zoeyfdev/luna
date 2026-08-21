#ifdef PORTABLE
    #pragma bits 32
#else
    #pragma bits 16
#endif

extern void putc(char c, short short int color);
extern void sleep(int seconds);
extern void* flags_start;

void render_flags() {
    short short int* fptr = (short short int*) flags_start;
    while (*fptr != 0xFE) {
        for (int i = 0; i < 40; i++) {
            putc(0x20, *fptr);
        }
        fptr++;
        if (*fptr == 0x00) {
            sleep(1);
        }
    }
}

void _cstart() __attribute__((noreturn)) {
    #ifndef PORTABLE
        // Load the next sectors on the disk
        asm ("int 0x10");
        asm ("mov r2, r1");
        asm ("mov r1, 1");
        asm ("mov r3, r1");
        asm ("int 11");
    #endif

    // Set up PIT
   
    #ifdef PORTABLE
        asm ("mov r1, 0x6FFF0008");
    #else
        asm ("mov r1, 0xFA41");
    #endif
    
    asm ("mov r2, pit_nxt");

    #ifdef PORTABLE
        asm ("str32 r1, r2");
    #else
        asm ("str16 r1, r2");
    #endif

    while (1) {
        render_flags();
    }
}
