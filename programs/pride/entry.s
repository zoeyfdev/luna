.bits 16
jmp _entry

_entry:
    mov sp, 0xEFFF
    mov fp, sp
    jmp _start
