.bits 32
#define _builtin_lcc_basin_play_sound 13
.global play_sound

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
    mov r1, fp + 0
    str_ptr r1, e0
    mov r1, fp + 4
    str32 r1, e1
    mov r1, fp + 8
    str r1, e2
    mov r1, 2147483657
    mov r4, r1
    mov r2, r1
    mov r4, r2
    mov r7, fp + 9
    str_ptr r7, r4
    mov r1, fp + 9 // Variable name: done_flag, internal: fp + 9
    mov r2, r1
    lod_ptr r2, r2
    mov r4, r2
    mov r1, 0
    mov r5, r1
    str r4, r5
    mov r1, 2147483649
    mov r4, r1
    mov r2, r1
    mov r4, r2
    mov r1, fp + 0 // Variable name: buffer, internal: fp + 0
    lod_ptr r1, r2
    mov r5, r2
    str32 r4, r5
    mov r1, 2147483653
    mov r4, r1
    mov r2, r1
    mov r4, r2
    mov r1, fp + 4 // Variable name: size, internal: fp + 4
    lod32 r1, r2
    mov r5, r2
    str32 r4, r5
    mov r1, 2147483648
    mov r4, r1
    mov r2, r1
    mov r4, r2
    mov r1, 1
    mov r5, r1
    str r4, r5
    mov r1, fp + 8 // Variable name: block, internal: fp + 8
    lod r1, r2
    mov r11, r2
    mov r1, 1
    mov r5, r1
    cmp r11, r11, r5
    jnz r11, if_stmt_8
    jz r11, after_stmt_10
if_stmt_8:
while_stmt_11_check:
    mov r1, fp + 9 // Variable name: done_flag, internal: fp + 9
    lod_ptr r1, r2
    lod r2, r2
    mov r11, r2
    mov r1, 0
    mov r5, r1
    cmp r11, r11, r5
    jnz r11, while_stmt_11_body
    jmp while_stmt_11_after
while_stmt_11_body:
    jmp while_stmt_11_check
while_stmt_11_after:
    jmp after_stmt_10
after_stmt_10:
    pop e11
    pop fp
    ret
    pop e11
    pop fp
    ret
