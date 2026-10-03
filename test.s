mov sp, 512
push string
call print
jmp pc

print:
    pop e11
    pop r4

    mov r2, 255

    mov e10, pc

    lod r4, r1
    int 1
    inc r4

    jnz r1, e10
    ret

string:
    .asciz "Hello world!"
