.bits 32
.global ffnt
.global fntf
.global fcreate
.global find_file
.global fopen
.global fgetsize
.global flist
.global fwrite
.global fstrap

#define _builtin_lcc_basin_ffnt 13
#define _builtin_lcc_basin_fntf 29
#define _builtin_lcc_basin_fcreate 28
#define _builtin_lcc_basin_find_file 16
#define _builtin_lcc_basin_fopen 39
#define _builtin_lcc_basin_fgetsize 12
#define _builtin_lcc_basin_flist 15
#define _builtin_lcc_basin_fwrite 20
#define _builtin_lcc_basin_fstrap 20

var_1:
    .ptr PROMPTBUF

var_2:
    .ptr renderbuf_loc

var_3:
    .ptr sleep_loc

var_4:
    .ptr malloc_loc

var_5:
    .ptr puts32_loc

var_str_13:
    .asciz "File '"

var_str_14:
    .asciz "' not found!\n"

var_str_17:
    .asciz "\n"


ffnt:
    pop e11
    pop e0
    push fp
    mov r12, _builtin_lcc_basin_ffnt
    sub fp, fp, r12
    push e11
    mov r0, fp + 0
    str_ptr r0, e0
    mov r1, fp + 4
    // Push allocated registers
    push r1
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r2, fp + 0
    lod_ptr r2, r2
    push r2
    call strlen
    // Pop allocated registers
    pop r1
    // End pop
    mov r2, e6
    push r2
    call malloc
    // Pop allocated registers
    pop r1
    // End pop
    mov r2, e6
    str_ptr r1, r2
    mov r1, fp + 8
    mov r2, fp + 4
    lod_ptr r2, r2
    str_ptr r1, r2
    mov r1, fp + 12
    mov r2, 0
    str r1, r2
while_stmt_0_check:
    mov r1, fp + 0
    lod_ptr r1, r1
    lod r1, r1
    mov r2, 0
    cmp r1, r1, r2
    mov r3, 1
    xor r1, r1, r3
    jnz r1, while_stmt_0_body
    jmp while_stmt_0_after
while_stmt_0_body:
if_stmt_1_check:
    mov r1, fp + 0
    lod_ptr r1, r1
    lod r1, r1
    mov r2, 0x20
    cmp r1, r1, r2
    mov r3, 1
    xor r1, r1, r3
    jnz r1, if_stmt_1_success
    jmp if_stmt_1_else
if_stmt_1_success:
    mov r1, fp + 4
    lod_ptr r1, r1
    mov r2, fp + 0
    lod_ptr r2, r2
    lod r2, r2
    str r1, r2
    mov r2, fp + 4
    mov r3, fp + 4
    lod_ptr r3, r3
    mov r1, r3
    inc r1
    str_ptr r2, r1
    jmp if_stmt_1_done
if_stmt_1_else:
if_stmt_2_check:
    mov r1, fp + 12
    lod r1, r1
    mov r2, 0
    cmp r1, r1, r2
    jnz r1, if_stmt_2_success
    jmp if_stmt_2_else
if_stmt_2_success:
    mov r1, fp + 12
    mov r2, 1
    str r1, r2
    mov r1, fp + 4
    lod_ptr r1, r1
    mov r2, 0x2E
    str r1, r2
    mov r2, fp + 4
    mov r3, fp + 4
    lod_ptr r3, r3
    mov r1, r3
    inc r1
    str_ptr r2, r1
    jmp if_stmt_2_done
if_stmt_2_else:
if_stmt_2_done:
if_stmt_1_done:
    mov r2, fp + 0
    mov r3, fp + 0
    lod_ptr r3, r3
    mov r1, r3
    inc r1
    str_ptr r2, r1
    jmp while_stmt_0_check
while_stmt_0_after:
    mov r1, fp + 8
    lod_ptr r1, r1
    mov e6, r1
    jmp .ffnt_ret
.ffnt_ret:
    pop e11
    pop fp
    ret

fntf:
    pop e11
    pop e0
    push fp
    mov r12, _builtin_lcc_basin_fntf
    sub fp, fp, r12
    push e11
    mov r0, fp + 0
    str_ptr r0, e0
    mov r1, fp + 4
    // Push allocated registers
    push r1
    // End push
    mov r2, 16
    push r2
    call malloc
    // Pop allocated registers
    pop r1
    // End pop
    mov r2, e6
    str_ptr r1, r2
    mov r1, fp + 8
    mov r2, fp + 4
    lod_ptr r2, r2
    str_ptr r1, r2
    mov r1, fp + 12
    mov r2, 0
    str r1, r2
    mov r1, fp + 13
    mov r2, 0
    str32 r1, r2
while_stmt_3_check:
    mov r1, fp + 0
    lod_ptr r1, r1
    lod r1, r1
    mov r2, 0
    cmp r1, r1, r2
    mov r3, 1
    xor r1, r1, r3
    jnz r1, while_stmt_3_body
    jmp while_stmt_3_after
while_stmt_3_body:
if_stmt_4_check:
    mov r1, fp + 0
    lod_ptr r1, r1
    lod r1, r1
    mov r2, 0x2e
    cmp r1, r1, r2
    jnz r1, if_stmt_4_success
    jmp if_stmt_4_else
if_stmt_4_success:
if_stmt_5_check:
    mov r1, fp + 12
    lod r1, r1
    mov r2, 0
    cmp r1, r1, r2
    jnz r1, if_stmt_5_success
    jmp if_stmt_5_else
if_stmt_5_success:
    mov r1, fp + 12
    mov r2, 1
    str r1, r2
    mov r1, fp + 17
    mov r2, 12
    str32 r1, r2
    mov r1, fp + 21
    mov r2, fp + 17
    lod32 r2, r2
    mov r3, fp + 13
    lod32 r3, r3
    sub r2, r2, r3
    str32 r1, r2
for_stmt_6_init:
    mov r1, fp + 25
    mov r2, 0
    str32 r1, r2
for_stmt_6_check:
    mov r1, fp + 25
    lod32 r1, r1
    mov r2, fp + 21
    lod32 r2, r2
    ilt r1, r1, r2
    jnz r1, for_stmt_6_body
    jmp for_stmt_6_after
for_stmt_6_body:
    mov r2, fp + 4
    lod_ptr r2, r2
    mov r3, 0x20
    str r2, r3
    mov r3, fp + 4
    mov r4, fp + 4
    lod_ptr r4, r4
    mov r2, r4
    inc r2
    str_ptr r3, r2
for_stmt_6_iterator:
    mov r3, fp + 25
    mov r4, fp + 25
    lod32 r4, r4
    mov r2, r4
    inc r2
    str32 r3, r2
    jmp for_stmt_6_check
for_stmt_6_after:
    mov r3, fp + 0
    mov r5, fp + 0
    lod_ptr r5, r5
    mov r2, r5
    inc r2
    str_ptr r3, r2
    jmp if_stmt_5_done
if_stmt_5_else:
if_stmt_5_done:
    jmp if_stmt_4_done
if_stmt_4_else:
    mov r2, fp + 4
    lod_ptr r2, r2
    mov r3, fp + 0
    lod_ptr r3, r3
    lod r3, r3
    str r2, r3
    mov r3, fp + 13
    mov r5, fp + 13
    lod32 r5, r5
    mov r2, r5
    inc r2
    str32 r3, r2
    mov r3, fp + 0
    mov r5, fp + 0
    lod_ptr r5, r5
    mov r2, r5
    inc r2
    str_ptr r3, r2
    mov r3, fp + 4
    mov r5, fp + 4
    lod_ptr r5, r5
    mov r2, r5
    inc r2
    str_ptr r3, r2
if_stmt_4_done:
    jmp while_stmt_3_check
while_stmt_3_after:
    // Push allocated registers
    push r1
    push r4
    // End push
    mov r2, 16
    push r2
    call free
    // Pop allocated registers
    pop r4
    pop r1
    // End pop
    mov r2, e6
    mov r2, fp + 8
    lod_ptr r2, r2
    mov e6, r2
    jmp .fntf_ret
.fntf_ret:
    pop e11
    pop fp
    ret

fcreate:
    pop e11
    pop e1
    pop e0
    push fp
    mov r12, _builtin_lcc_basin_fcreate
    sub fp, fp, r12
    push e11
    mov r0, fp + 0
    str_ptr r0, e0
    mov r0, fp + 4
    str32 r0, e1
    mov r2, fp + 8
    mov r3, 0x61C
    str_ptr r2, r3
    mov r2, fp + 12
    mov r3, fp + 8
    lod_ptr r3, r3
    lod_ptr r3, r3
    str_ptr r2, r3
    mov r2, fp + 12
    lod_ptr r2, r2
    mov r3, 0x4c465346
    str32 r2, r3
    mov r2, fp + 12
    mov r3, fp + 12
    lod_ptr r3, r3
    mov r5, 4
    add r3, r3, r5
    str_ptr r2, r3
    // Push allocated registers
    push r1
    push r4
    // End push
    mov r2, fp + 0
    lod_ptr r2, r2
    push r2
    mov r2, fp + 12
    lod_ptr r2, r2
    push r2
    call strcpy
    // Pop allocated registers
    pop r4
    pop r1
    // End pop
    mov r2, e6
    mov r2, fp + 16
    // Push allocated registers
    push r1
    push r2
    push r4
    // End push
    mov r3, fp + 0
    lod_ptr r3, r3
    push r3
    call strlen
    // Pop allocated registers
    pop r4
    pop r2
    pop r1
    // End pop
    mov r3, e6
    str32 r2, r3
    mov r2, fp + 12
    mov r3, fp + 12
    lod_ptr r3, r3
    mov r5, fp + 16
    lod32 r5, r5
    add r3, r3, r5
    str_ptr r2, r3
    mov r2, fp + 12
    lod_ptr r2, r2
    mov r3, fp + 4
    lod32 r3, r3
    str32 r2, r3
    mov r2, fp + 12
    mov r3, fp + 12
    lod_ptr r3, r3
    mov r5, 4
    add r3, r3, r5
    str_ptr r2, r3
for_stmt_7_init:
    mov r2, fp + 20
    mov r3, 0
    str32 r2, r3
for_stmt_7_check:
    mov r2, fp + 20
    lod32 r2, r2
    mov r3, fp + 4
    lod32 r3, r3
    ilt r2, r2, r3
    jnz r2, for_stmt_7_body
    jmp for_stmt_7_after
for_stmt_7_body:
    mov r3, fp + 12
    lod_ptr r3, r3
    mov r5, 0x00
    str32 r3, r5
    mov r5, fp + 12
    mov r6, fp + 12
    lod_ptr r6, r6
    mov r3, r6
    inc r3
    str_ptr r5, r3
for_stmt_7_iterator:
    mov r5, fp + 20
    mov r6, fp + 20
    lod32 r6, r6
    mov r3, r6
    inc r3
    str32 r5, r3
    jmp for_stmt_7_check
for_stmt_7_after:
    mov r3, fp + 24
    mov r5, fp + 12
    lod_ptr r5, r5
    mov r7, 512
    div r5, r5, r7
    str32 r3, r5
    // Push allocated registers
    push r1
    push r2
    push r4
    push r6
    // End push
    mov r3, fp + 24
    lod32 r3, r3
    push r3
    call save_sector
    // Pop allocated registers
    pop r6
    pop r4
    pop r2
    pop r1
    // End pop
    mov r3, e6
    mov r3, fp + 8
    lod_ptr r3, r3
    mov r5, fp + 8
    lod_ptr r5, r5
    lod_ptr r5, r5
    mov r7, fp + 4
    lod32 r7, r7
    add r5, r5, r7
    str_ptr r3, r5
    // Push allocated registers
    push r1
    push r2
    push r4
    push r6
    // End push
    mov r3, 3
    push r3
    call save_sector
    // Pop allocated registers
    pop r6
    pop r4
    pop r2
    pop r1
    // End pop
    mov r3, e6
.fcreate_ret:
    pop e11
    pop fp
    ret

find_file:
    pop e11
    pop e0
    push fp
    mov r12, _builtin_lcc_basin_find_file
    sub fp, fp, r12
    push e11
    mov r0, fp + 0
    str_ptr r0, e0
    mov r3, fp + 4
    mov r5, 0x618
    str_ptr r3, r5
    mov r3, fp + 8
    mov r5, fp + 4
    lod_ptr r5, r5
    lod_ptr r5, r5
    str_ptr r3, r5
while_stmt_8_check:
    mov r3, 1
    jnz r3, while_stmt_8_body
    jmp while_stmt_8_after
while_stmt_8_body:
if_stmt_9_check:
    mov r3, fp + 8
    lod_ptr r3, r3
    lod32 r3, r3
    mov r5, 0x4c465346
    cmp r3, r3, r5
    mov r7, 1
    xor r3, r3, r7
    jnz r3, if_stmt_9_success
    jmp if_stmt_9_else
if_stmt_9_success:
    jmp while_stmt_8_after
    jmp if_stmt_9_done
if_stmt_9_else:
if_stmt_9_done:
    mov r3, fp + 8
    mov r5, fp + 8
    lod_ptr r5, r5
    mov r7, 4
    add r5, r5, r7
    str_ptr r3, r5
if_stmt_10_check:
    // Push allocated registers
    push r1
    push r2
    push r4
    push r6
    // End push
    mov r3, fp + 0
    lod_ptr r3, r3
    push r3
    mov r3, fp + 8
    lod_ptr r3, r3
    push r3
    call strcmp
    // Pop allocated registers
    pop r6
    pop r4
    pop r2
    pop r1
    // End pop
    mov r3, e6
    mov r5, 1
    cmp r3, r3, r5
    jnz r3, if_stmt_10_success
    jmp if_stmt_10_else
if_stmt_10_success:
    mov r3, fp + 8
    mov r5, fp + 8
    lod_ptr r5, r5
    mov r7, 20
    add r5, r5, r7
    str_ptr r3, r5
    mov r3, fp + 8
    lod_ptr r3, r3
    mov e6, r3
    jmp .find_file_ret
    jmp if_stmt_10_done
if_stmt_10_else:
    mov r3, fp + 8
    mov r5, fp + 8
    lod_ptr r5, r5
    mov r7, 16
    add r5, r5, r7
    str_ptr r3, r5
    mov r3, fp + 12
    mov r5, fp + 8
    lod_ptr r5, r5
    lod32 r5, r5
    str32 r3, r5
    mov r3, fp + 8
    mov r5, fp + 8
    lod_ptr r5, r5
    mov r7, 4
    add r5, r5, r7
    str_ptr r3, r5
    mov r3, fp + 8
    mov r5, fp + 8
    lod_ptr r5, r5
    mov r7, fp + 12
    lod32 r7, r7
    add r5, r5, r7
    str_ptr r3, r5
if_stmt_10_done:
    jmp while_stmt_8_check
while_stmt_8_after:
    mov r3, 0
    mov e6, r3
    jmp .find_file_ret
.find_file_ret:
    pop e11
    pop fp
    ret

fopen:
    pop e11
    pop e1
    pop e0
    push fp
    mov r12, _builtin_lcc_basin_fopen
    sub fp, fp, r12
    push e11
    mov r0, fp + 0
    str_ptr r0, e0
    mov r0, fp + 4
    str r0, e1
    mov r3, fp + 5
    // Push allocated registers
    push r1
    push r2
    push r3
    push r4
    push r6
    // End push
    mov r5, fp + 0
    lod_ptr r5, r5
    push r5
    call find_file
    // Pop allocated registers
    pop r6
    pop r4
    pop r3
    pop r2
    pop r1
    // End pop
    mov r5, e6
    str_ptr r3, r5
if_stmt_11_check:
    mov r3, fp + 5
    lod_ptr r3, r3
    mov r5, 0x00000000
    cmp r3, r3, r5
    jnz r3, if_stmt_11_success
    jmp if_stmt_11_else
if_stmt_11_success:
if_stmt_12_check:
    mov r3, fp + 4
    lod r3, r3
    mov r5, 1
    cmp r3, r3, r5
    jnz r3, if_stmt_12_success
    jmp if_stmt_12_else
if_stmt_12_success:
    // Push allocated registers
    push r1
    push r2
    push r4
    push r6
    // End push
    // Push allocated registers
    push r1
    push r2
    push r3
    push r4
    push r6
    // End push
    mov r3, fp + 17
    push r3
    push var_str_13
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r6
    pop r4
    pop r3
    pop r2
    pop r1
    // End pop
    mov r3, fp + 17
    push r3
    mov r3, 0xE9
    push r3
    mov r3, 0x00
    push r3
    call puts32
    // Pop allocated registers
    pop r6
    pop r4
    pop r2
    pop r1
    // End pop
    mov r3, e6
    // Push allocated registers
    push r1
    push r2
    push r4
    push r6
    // End push
    // Push allocated registers
    push r1
    push r2
    push r4
    push r6
    // End push
    mov r3, fp + 0
    lod_ptr r3, r3
    push r3
    call ffnt
    // Pop allocated registers
    pop r6
    pop r4
    pop r2
    pop r1
    // End pop
    mov r3, e6
    push r3
    mov r3, 0xE9
    push r3
    mov r3, 0x00
    push r3
    call puts32
    // Pop allocated registers
    pop r6
    pop r4
    pop r2
    pop r1
    // End pop
    mov r3, e6
    // Push allocated registers
    push r1
    push r2
    push r4
    push r6
    // End push
    // Push allocated registers
    push r1
    push r2
    push r3
    push r4
    push r6
    // End push
    mov r3, fp + 24
    push r3
    push var_str_14
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r6
    pop r4
    pop r3
    pop r2
    pop r1
    // End pop
    mov r3, fp + 24
    push r3
    mov r3, 0xE9
    push r3
    mov r3, 0x00
    push r3
    call puts32
    // Pop allocated registers
    pop r6
    pop r4
    pop r2
    pop r1
    // End pop
    mov r3, e6
    mov r3, fp + 9
    // Above should not be loaded
    mov r5, 0
    add r3, r3, r5
    mov r5, 0x00000000
    str_ptr r3, r5
    mov r3, fp + 9
    mov e6, r3
    jmp .fopen_ret
    jmp if_stmt_12_done
if_stmt_12_else:
    mov r3, fp + 9
    // Above should not be loaded
    mov r5, 0
    add r3, r3, r5
    mov r5, 0x00000000
    str_ptr r3, r5
    mov r3, fp + 9
    mov e6, r3
    jmp .fopen_ret
if_stmt_12_done:
    jmp if_stmt_11_done
if_stmt_11_else:
if_stmt_11_done:
    mov r3, fp + 9
    // Above should not be loaded
    mov r5, 0
    add r3, r3, r5
    mov r5, fp + 5
    lod_ptr r5, r5
    str_ptr r3, r5
    mov r3, fp + 9
    // Above should not be loaded
    mov r5, 4
    add r3, r3, r5
    // Push allocated registers
    push r1
    push r2
    push r3
    push r4
    push r6
    // End push
    mov r5, fp + 0
    lod_ptr r5, r5
    push r5
    call ffnt
    // Pop allocated registers
    pop r6
    pop r4
    pop r3
    pop r2
    pop r1
    // End pop
    mov r5, e6
    str_ptr r3, r5
    mov r3, fp + 9
    mov e6, r3
    jmp .fopen_ret
.fopen_ret:
    pop e11
    pop fp
    ret

fgetsize:
    pop e11
    pop e0
    push fp
    mov r12, _builtin_lcc_basin_fgetsize
    sub fp, fp, r12
    push e11
    mov r0, fp + 0
    str_ptr r0, e0
    mov r3, fp + 4
    // Push allocated registers
    push r1
    push r2
    push r3
    push r4
    push r6
    // End push
    mov r5, fp + 0
    lod_ptr r5, r5
    push r5
    mov r5, 1
    push r5
    call fopen
    // Pop allocated registers
    pop r6
    pop r4
    pop r3
    pop r2
    pop r1
    // End pop
    mov r5, e6
    str_ptr r3, r5
    mov r3, fp + 8
    mov r5, fp + 4
    // Above should not be loaded
    lod_ptr r5, r5
    mov r7, 0
    add r5, r5, r7
    lod_ptr r5, r5
    str_ptr r3, r5
    mov r3, fp + 8
    mov r5, fp + 8
    lod_ptr r5, r5
    mov r7, 4
    sub r5, r5, r7
    str_ptr r3, r5
    mov r3, fp + 8
    lod_ptr r3, r3
    lod32 r3, r3
    mov e6, r3
    jmp .fgetsize_ret
.fgetsize_ret:
    pop e11
    pop fp
    ret

flist:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_flist
    sub fp, fp, r12
    push e11
    mov r3, fp + 0
    mov r5, 0x618
    str_ptr r3, r5
    mov r3, fp + 4
    mov r5, fp + 0
    lod_ptr r5, r5
    lod_ptr r5, r5
    str_ptr r3, r5
while_stmt_15_check:
    mov r3, 1
    jnz r3, while_stmt_15_body
    jmp while_stmt_15_after
while_stmt_15_body:
if_stmt_16_check:
    mov r3, fp + 4
    lod_ptr r3, r3
    lod32 r3, r3
    mov r5, 0x4c465346
    cmp r3, r3, r5
    mov r7, 1
    xor r3, r3, r7
    jnz r3, if_stmt_16_success
    jmp if_stmt_16_else
if_stmt_16_success:
    jmp while_stmt_15_after
    jmp if_stmt_16_done
if_stmt_16_else:
if_stmt_16_done:
    mov r3, fp + 4
    mov r5, fp + 4
    lod_ptr r5, r5
    mov r7, 4
    add r5, r5, r7
    str_ptr r3, r5
    // Push allocated registers
    push r1
    push r2
    push r4
    push r6
    // End push
    // Push allocated registers
    push r1
    push r2
    push r4
    push r6
    // End push
    mov r3, fp + 4
    lod_ptr r3, r3
    push r3
    call ffnt
    // Pop allocated registers
    pop r6
    pop r4
    pop r2
    pop r1
    // End pop
    mov r3, e6
    push r3
    mov r3, 0xFF
    push r3
    mov r3, 0x00
    push r3
    call puts32
    // Pop allocated registers
    pop r6
    pop r4
    pop r2
    pop r1
    // End pop
    mov r3, e6
    // Push allocated registers
    push r1
    push r2
    push r4
    push r6
    // End push
    // Push allocated registers
    push r1
    push r2
    push r3
    push r4
    push r6
    // End push
    mov r3, fp + 8
    push r3
    push var_str_17
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r6
    pop r4
    pop r3
    pop r2
    pop r1
    // End pop
    mov r3, fp + 8
    push r3
    mov r3, 0xFF
    push r3
    mov r3, 0x00
    push r3
    call puts32
    // Pop allocated registers
    pop r6
    pop r4
    pop r2
    pop r1
    // End pop
    mov r3, e6
    mov r3, fp + 4
    mov r5, fp + 4
    lod_ptr r5, r5
    mov r7, 16
    add r5, r5, r7
    str_ptr r3, r5
    mov r3, fp + 11
    mov r5, fp + 4
    lod_ptr r5, r5
    lod32 r5, r5
    str32 r3, r5
    mov r3, fp + 4
    mov r5, fp + 4
    lod_ptr r5, r5
    mov r7, 4
    add r5, r5, r7
    str_ptr r3, r5
    mov r3, fp + 4
    mov r5, fp + 4
    lod_ptr r5, r5
    mov r7, fp + 11
    lod32 r7, r7
    add r5, r5, r7
    str_ptr r3, r5
    jmp while_stmt_15_check
while_stmt_15_after:
    jmp .flist_ret
.flist_ret:
    pop e11
    pop fp
    ret

fwrite:
    pop e11
    pop e1
    pop e0
    push fp
    mov r12, _builtin_lcc_basin_fwrite
    sub fp, fp, r12
    push e11
    mov r0, fp + 0
    str_ptr r0, e0
    mov r0, fp + 4
    str_ptr r0, e1
    mov r3, fp + 8
    // Push allocated registers
    push r1
    push r2
    push r3
    push r4
    push r6
    // End push
    mov r5, fp + 0
    lod_ptr r5, r5
    push r5
    mov r5, 1
    push r5
    call fopen
    // Pop allocated registers
    pop r6
    pop r4
    pop r3
    pop r2
    pop r1
    // End pop
    mov r5, e6
    str_ptr r3, r5
    mov r3, fp + 12
    mov r5, fp + 8
    // Above should not be loaded
    lod_ptr r5, r5
    mov r7, 0
    add r5, r5, r7
    lod_ptr r5, r5
    str_ptr r3, r5
if_stmt_18_check:
    mov r3, fp + 12
    lod_ptr r3, r3
    mov r5, 0
    cmp r3, r3, r5
    jnz r3, if_stmt_18_success
    jmp if_stmt_18_else
if_stmt_18_success:
    jmp .fwrite_ret
    jmp if_stmt_18_done
if_stmt_18_else:
if_stmt_18_done:
    // Push allocated registers
    push r1
    push r2
    push r4
    push r6
    // End push
    mov r3, fp + 4
    lod_ptr r3, r3
    push r3
    mov r3, fp + 12
    lod_ptr r3, r3
    push r3
    call strcpy
    // Pop allocated registers
    pop r6
    pop r4
    pop r2
    pop r1
    // End pop
    mov r3, e6
    mov r3, fp + 16
    mov r5, fp + 12
    lod_ptr r5, r5
    mov r7, 512
    div r5, r5, r7
    str32 r3, r5
    // Push allocated registers
    push r1
    push r2
    push r4
    push r6
    // End push
    mov r3, fp + 16
    lod32 r3, r3
    push r3
    call save_sector
    // Pop allocated registers
    pop r6
    pop r4
    pop r2
    pop r1
    // End pop
    mov r3, e6
    // Push allocated registers
    push r1
    push r2
    push r4
    push r6
    // End push
    mov r3, fp + 16
    lod32 r3, r3
    mov r5, 1
    sub r3, r3, r5
    push r3
    call save_sector
    // Pop allocated registers
    pop r6
    pop r4
    pop r2
    pop r1
    // End pop
    mov r3, e6
    // Push allocated registers
    push r1
    push r2
    push r4
    push r6
    // End push
    mov r3, fp + 16
    lod32 r3, r3
    mov r5, 1
    add r3, r3, r5
    push r3
    call save_sector
    // Pop allocated registers
    pop r6
    pop r4
    pop r2
    pop r1
    // End pop
    mov r3, e6
.fwrite_ret:
    pop e11
    pop fp
    ret

fstrap:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_fstrap
    sub fp, fp, r12
    push e11
    mov r3, fp + 0
    mov r5, 0x618
    str_ptr r3, r5
    mov r3, fp + 4
    mov r5, fp + 0
    lod_ptr r5, r5
    lod_ptr r5, r5
    str_ptr r3, r5
while_stmt_19_check:
    mov r3, 1
    jnz r3, while_stmt_19_body
    jmp while_stmt_19_after
while_stmt_19_body:
if_stmt_20_check:
    mov r3, fp + 4
    lod_ptr r3, r3
    lod32 r3, r3
    mov r5, 0x4c465346
    cmp r3, r3, r5
    mov r7, 1
    xor r3, r3, r7
    jnz r3, if_stmt_20_success
    jmp if_stmt_20_else
if_stmt_20_success:
    jmp while_stmt_19_after
    jmp if_stmt_20_done
if_stmt_20_else:
if_stmt_20_done:
    mov r3, fp + 8
    mov r5, fp + 4
    lod_ptr r5, r5
    mov r7, 512
    div r5, r5, r7
    str32 r3, r5
    mov r3, fp + 4
    mov r5, fp + 4
    lod_ptr r5, r5
    mov r7, 20
    add r5, r5, r7
    str_ptr r3, r5
    mov r3, fp + 12
    mov r5, fp + 4
    lod_ptr r5, r5
    lod32 r5, r5
    str32 r3, r5
    mov r3, fp + 16
    mov r5, fp + 12
    lod32 r5, r5
    mov r7, 512
    div r5, r5, r7
    mov r7, 1
    add r5, r5, r7
    str32 r3, r5
    // Push allocated registers
    push r1
    push r2
    push r4
    push r6
    // End push
    mov r3, fp + 16
    lod32 r3, r3
    push r3
    mov r3, fp + 8
    lod32 r3, r3
    push r3
    call offset_sec_load
    // Pop allocated registers
    pop r6
    pop r4
    pop r2
    pop r1
    // End pop
    mov r3, e6
    mov r3, fp + 0
    mov r5, fp + 0
    lod_ptr r5, r5
    mov r7, fp + 12
    lod32 r7, r7
    add r5, r5, r7
    str_ptr r3, r5
    jmp while_stmt_19_check
while_stmt_19_after:
.fstrap_ret:
    pop e11
    pop fp
    ret

