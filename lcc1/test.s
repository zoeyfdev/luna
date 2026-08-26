.bits 16
.global foo

#define _builtin_lcc_basin_foo 0
var_1:
    .ptr 2

var_2:
    .ptr 1


foo:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_foo
    sub fp, fp, r12
    push e11
    mov r1, var_1
    mov r2, var_2
    lod16 r2, r2
    str r1, r2
    pop e11
    pop fp
    ret

