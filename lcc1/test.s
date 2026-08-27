.bits 16
.global foo

#define _builtin_lcc_basin_foo 4

foo:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_foo
    sub fp, fp, r12
    push e11
    mov r1, fp + 0
    mov r2, 1
    str_ptr r1, r2
    mov r1, fp + 2
    mov r2, 2
    str_ptr r1, r2
    mov r1, fp + 0
    mov r2, fp + 2
    lod_ptr r2, r2
    lod16 r2, r2
    str_ptr r1, r2
    pop e11
    pop fp
    ret

