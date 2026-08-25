.bits 16
.global TestFunc
.global main

#define _builtin_lcc_basin_TestFunc 0
#define _builtin_lcc_basin_main 0

var_1:
    .word 1

TestFunc:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_TestFunc
    sub fp, fp, r12
    push e11
    pop e11
    pop fp
    ret

main:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_main
    sub fp, fp, r12
    push e11
    call TestFunc
    mov r1, e6
    mov e6, r1
    ret
    pop e11
    pop fp
    ret

