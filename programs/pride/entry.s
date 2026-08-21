#ifdef PORTABLE
    .bits 32
#else
    .bits 16
#endif

jmp _start

_start:
    #ifdef PORTABLE
        mov sp, 0x60000000
        mov fp, 0x6000F000

        mov r1, 0x6FFF001A
        mov r2, key_click
        str32 r1, r2

        mov r1, 0x6FFF0019
        mov r2, 1
        str r1, r2 // ENABLE KEYBOARD INTERRUPT
    #else
        mov sp, 0xEFFF
        mov fp, 0xE000
    #endif

    jmp _cstart
