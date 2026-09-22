mov sp, 0x200

push 1000
pop r1

mov r1, 0
lod16 r1, r2

jmp pc
