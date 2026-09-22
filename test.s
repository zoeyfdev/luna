mov r4, string
mov r2, 255
mov r3, 0

mov r5, "xf"
str16 r4, r5

mov e10, pc

lod r4, r1
int 1
inc r4

jnz r1, e10
jmp pc

string:
    .asciz "Hello world!"
