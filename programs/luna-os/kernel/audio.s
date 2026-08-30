.bits 32
.global play_sound

#define _builtin_lcc_basin_play_sound 13

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

var_str_0:
    .asciz "NEW!"
var_ptr_1:
    .ptr var_str_0

.global play_sound_loc
play_sound_loc:
play_sound:
    pop e11
    pop e2
    pop e1
    pop e0
    push fp
    mov r12, _builtin_lcc_basin_play_sound
    sub fp, fp, r12
    push e11
    mov r0, fp + 0
    str_ptr r0, e0
    mov r0, fp + 4
    str32 r0, e1
    mov r0, fp + 8
    str r0, e2
    mov r1, var_ptr_1
    lod_ptr r1, r1
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    mov r1, e6
    mov r1, var_6
    mov r2, 0x80000009
    str_ptr r1, r2
    mov r1, fp + 9
    lod_ptr r1, r1
    mov r2, 0
    str r1, r2
    mov r1, 0x80000001
    mov r2, fp + 0
    lod_ptr r2, r2
    str32 r1, r2
    mov r1, 0x80000005
    mov r2, fp + 4
    lod32 r2, r2
    str32 r1, r2
    mov r1, 0x80000000
    mov r2, 1
    str r1, r2
if_stmt_2_check:
    mov r1, fp + 8
    lod r1, r1
    mov r2, 1
    cmp r1, r1, r2
    jnz r1, if_stmt_2_success
    jmp if_stmt_2_else
if_stmt_2_success:
while_stmt_3_check:
    mov r1, fp + 9
    lod_ptr r1, r1
    lod r1, r1
    mov r2, 0
    cmp r1, r1, r2
    jnz r1, while_stmt_3_body
    jmp while_stmt_3_after
while_stmt_3_body:
    jmp while_stmt_3_check
while_stmt_3_after:
    jmp if_stmt_2_done
if_stmt_2_else:
if_stmt_2_done:
    pop e11
    pop fp
    ret

