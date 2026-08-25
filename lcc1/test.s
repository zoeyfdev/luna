.bits 16
.global thing
.global main

#define _builtin_lcc_basin_thing 2
#define _builtin_lcc_basin_main 2
var_str_0:
    .asciz "Hello world!"
var_ptr_1:
    .ptr var_str_0

thing:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_thing
    sub fp, fp, r12
    pop e0
    push e11
    mov r0, fp + 0
    str16 r0, e0
    pop e11
    pop fp
    ret

main:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_main
    sub fp, fp, r12
    pop e0
    push e11
    mov r0, fp + 0
    str16 r0, e0
    mov r1, var_ptr_1
    push r1
    call thing
    mov r1, e6
    pop e11
    pop fp
    ret

