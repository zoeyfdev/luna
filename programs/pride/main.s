.bits 16
.global render_flags
.global _cstart

#define _builtin_lcc_basin_render_flags 4
#define _builtin_lcc_basin__cstart 0

var_1:
    .ptr flags_start


render_flags:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_render_flags
    sub fp, fp, r12
    push e11
    mov r1, fp + 0
    mov r2, var_1
    lod_ptr r2, r2
    str_ptr r1, r2
    mov r1, fp + 2
    mov r2, 0
    str16 r1, r2
while_stmt_0_check:
    mov r1, fp + 0
    lod_ptr r1, r1
    lod r1, r1
    mov r2, 0xFE
    cmp r1, r1, r2
    mov r3, 1
    xor r1, r1, r3
    jnz r1, while_stmt_0_body
    jmp while_stmt_0_after
while_stmt_0_body:
for_stmt_1_check:
    mov r1, fp + 2
    lod16 r1, r1
    mov r2, 40
    ilt r1, r1, r2
    jnz r1, for_stmt_1_body
    jmp for_stmt_1_after
for_stmt_1_body:
    mov r2, 0x20
    push r2
    mov r2, fp + 0
    lod_ptr r2, r2
    lod r2, r2
    push r2
    call putc
    mov r2, e6
for_stmt_1_iterator:
    mov r3, fp + 2
    mov r4, fp + 2
    lod16 r4, r4
    mov r2, r4
    inc r2
    str16 r3, r2
    jmp for_stmt_1_check
for_stmt_1_after:
    .byte 0x21    // User-defined inline assembly
    mov r3, fp + 0
    mov r5, fp + 0
    lod_ptr r5, r5
    mov r2, r5
    inc r2
    str_ptr r3, r2
if_stmt_2_check:
    mov r2, fp + 0
    lod_ptr r2, r2
    lod r2, r2
    mov r3, 0x00
    cmp r2, r2, r3
    jnz r2, if_stmt_2_success
    jmp if_stmt_2_else
if_stmt_2_success:
    jmp if_stmt_2_done
if_stmt_2_else:
if_stmt_2_done:
    jmp while_stmt_0_check
while_stmt_0_after:
.render_flags_ret:
    pop e11
    pop fp
    ret

_cstart:
    push fp
    mov r12, _builtin_lcc_basin__cstart
    sub fp, fp, r12
    int 0x10    // User-defined inline assembly
    mov r2, r1    // User-defined inline assembly
    mov r1, 1    // User-defined inline assembly
    mov r3, r1    // User-defined inline assembly
    int 11    // User-defined inline assembly
    mov r1, 0xFA41    // User-defined inline assembly
    mov r2, pit_nxt    // User-defined inline assembly
    str16 r1, r2    // User-defined inline assembly
while_stmt_3_check:
    mov r2, 1
    jnz r2, while_stmt_3_body
    jmp while_stmt_3_after
while_stmt_3_body:
    call render_flags
    mov r2, e6
    jmp while_stmt_3_check
while_stmt_3_after:
._cstart_ret:
    pop fp

