set 32
.bits 32

start:
    mov sp, 0x500
    push 0x12345678
    pop r1
    jmp pc
