.bits 16
.global main

#define _builtin_lcc_basin_main 2

var_str_1:
    .asciz "hi"

var_ptr_2:
    .ptr var_str_1


main:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_main
    sub fp, fp, r12
    push e11
    mov r1, fp + 0
    mov r2, 1
    str16 r1, r2
for_stmt_0_check:
    mov r1, fp + 0
    lod16 r1, r1
    mov r2, 5
    ilt r1, r1, r2
    jnz r1, for_stmt_0_body
    jmp for_stmt_0_after
for_stmt_0_body:
    mov r2, var_ptr_2
    lod_ptr r2, r2
    push r2
    call printf
    mov r2, e6
for_stmt_0_iterator:
    mov r2, fp + 0
    mov r3, fp + 0
    lod16 r3, r3
    mov r4, 1
    add r3, r3, r4
    str16 r2, r3
    jmp for_stmt_0_check
for_stmt_0_after:
    mov r2, 0
    mov e6, r2
    jmp .main_ret
.main_ret:
    pop e11
    pop fp
    ret

