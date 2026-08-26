.bits 16
#define _builtin_lcc_basin_foo 2
.global foo

foo:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_foo
    sub fp, fp, r12
    push e11
    mov r1, 1
    mov r4, r1
    mov r7, fp + 0
    str16 r7, r4
    mov r1, fp + 0 // Variable name: a, internal: fp + 0
    mov r2, r1
    lod16 r2, r2
    mov r4, r2
    mov r1, 1
    mov r5, r1
    str16 r4, r5
    pop e11
    pop fp
    ret
