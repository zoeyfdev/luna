.bits 16
.global main

#define _builtin_lcc_basin_main 2


main:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_main
    sub fp, fp, r12
    push e11
    mov r1, fp + 0
    mov r2, 1
    str16 r1, r2
    mov r2, fp + 0
    mov r3, fp + 0
    lod16 r3, r3
    mov r1, r3
    dec r1
    str16 r2, r1
    mov r1, 0
    mov e6, r1
    jmp .main_ret
.main_ret:
    pop e11
    pop fp
    ret

