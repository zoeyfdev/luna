.bits 16
.global foo

#define _builtin_lcc_basin_foo 4
var_str_0:
    .asciz "Hello world!"
var_ptr_1:
    .ptr var_str_0

foo:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_foo
    sub fp, fp, r12
    push e11
    mov r1, fp + 0
    mov r2, 1
    str16 r1, r2
    mov r1, fp + 2
    mov r2, 2
    str_ptr r1, r2
if_stmt_0_check:
    mov r1, fp + 0
    lod16 r1, r1
    mov r2, 1
    cmp r1, r1, r2
    jnz r1, if_stmt_0_success
if_stmt_0_success:
    mov r1, var_ptr_1
    lod_ptr r1, r1
    push r1
    call print
    mov r1, e6
    pop e11
    pop fp
    ret

