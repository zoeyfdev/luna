.bits 16
.global main

#define _builtin_lcc_basin_main 0

var_1:
    .word 1

main:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_main
    sub fp, fp, r12
    push e11
    mov r1, var_1
    mov r2, 1
    str r1, r2
    mov r1, var_1
    lod16 r1, r1
    mov e6, r1
    ret
    pop e11
    pop fp
    ret
