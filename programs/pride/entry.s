.bits 32
jmp _entry

_entry:
    mov sp, 0x60000000
    mov fp, 0x6000F000

    mov r1, 0x6FFF001A
    mov r2, key_click
    str32 r1, r2

    mov r1, 0x6FFF0019
    mov r2, 1
    str r1, r2 // ENABLE KEYBOARD INTERRUPT

    jmp _start
