.bits 32
.global teststack
.global shell

#define _builtin_lcc_basin_teststack 8
#define _builtin_lcc_basin_shell 342

var_1:
    .ptr BOOT_IMG

var_2:
    .ptr PROMPTBUF

var_3:
    .ptr renderbuf_loc

var_4:
    .ptr sleep_loc

var_5:
    .ptr malloc_loc

var_6:
    .ptr puts32_loc

var_7:
    .ptr CRASH_SOUND

var_8:
    .ptr BOOT_SOUND

var_9:
    .ptr play_sound_loc

var_10:
    .byte 0

var_13:
    .byte 0

var_str_4:
    .asciz "reboot"

var_str_5:
    .asciz "Rebooting..."

var_str_7:
    .asciz "about"

var_str_8:
    .asciz "LunaOS 2.0.0\nBy Zoey Flax\n"

var_str_9:
    .asciz "\n\n"

var_str_11:
    .asciz "promptedit"

var_str_12:
    .asciz "Enter terminal prompt: "

var_str_13:
    .asciz "\n"

var_str_15:
    .asciz "open"

var_str_17:
    .asciz "Usage: open <filename>\n"

var_str_19:
    .asciz "\n"

var_str_21:
    .asciz "files"

var_str_22:
    .asciz "\n"

var_str_24:
    .asciz "shutdown"

var_str_25:
    .asciz "Shutting down...\n"

var_str_27:
    .asciz "testfault"

var_str_29:
    .asciz "clear"

var_str_31:
    .asciz "exec"

var_str_33:
    .asciz "battery"

var_str_34:
    .asciz "Battery level: "

var_str_35:
    .asciz "%\n"

var_str_37:
    .asciz "teststack"

var_str_38:
    .asciz "The first and third value should be the same.\n"

var_str_40:
    .asciz "time"

var_str_41:
    .asciz "Time: "

var_str_42:
    .asciz ":"

var_str_43:
    .asciz "\n"

var_str_44:
    .asciz "Bad command '"

var_str_45:
    .asciz "'\n"


teststack:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_teststack
    sub fp, fp, r12
    push e11
    mov r1, fp + 0
    mov r2, 0x90909090
    str_ptr r1, r2
    mov r1, fp + 4
    mov r2, fp + 0
    lod_ptr r2, r2
    lod32 r2, r2
    str32 r1, r2
    // Push allocated registers
    // Push allocated registers
    // End push
    mov r1, fp + 4
    lod32 r1, r1
    push r1
    mov r1, 1
    push r1
    call tohex
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
if_stmt_0_check:
    mov r1, var_10
    lod r1, r1
    mov r2, 0
    cmp r1, r1, r2
    jnz r1, if_stmt_0_success
    jmp if_stmt_0_else
if_stmt_0_success:
    mov r1, var_10
    mov r2, 1
    str r1, r2
    // Push allocated registers
    // Push allocated registers
    // End push
    call teststack
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jmp if_stmt_0_done
if_stmt_0_else:
    jmp .teststack_ret
if_stmt_0_done:
    // Push allocated registers
    // Push allocated registers
    // End push
    mov r1, fp + 4
    lod32 r1, r1
    push r1
    mov r1, 1
    push r1
    call tohex
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
.teststack_ret:
    pop e11
    pop fp
    ret

shell:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_shell
    sub fp, fp, r12
    push e11
while_stmt_1_check:
    mov r1, 1
    jnz r1, while_stmt_1_body
    jmp while_stmt_1_after
while_stmt_1_body:
if_stmt_2_check:
    mov r1, var_13
    lod r1, r1
    jnz r1, if_stmt_2_success
    jmp if_stmt_2_else
if_stmt_2_success:
    // Push allocated registers
    // Push allocated registers
    // End push
    mov r1, 256
    push r1
    call free
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jmp if_stmt_2_done
if_stmt_2_else:
if_stmt_2_done:
    mov r1, fp + 0
    // Push allocated registers
    // Push allocated registers
    push r1
    // End push
    mov r2, 256
    push r2
    call malloc
    // Pop saved registers
    // Pop allocated registers
    pop r1
    // End pop
    mov r2, e6
    str_ptr r1, r2
    mov r1, var_13
    mov r2, 1
    str r1, r2
    // Push allocated registers
    // Push allocated registers
    // End push
    mov r1, var_2
    lod_ptr r1, r1
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    // Push allocated registers
    // Push allocated registers
    // End push
    mov r1, fp + 0
    lod_ptr r1, r1
    push r1
    mov r1, 1
    push r1
    mov r1, 0
    push r1
    call readin
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    mov r1, fp + 4
    // Push allocated registers
    // Push allocated registers
    push r1
    // End push
    mov r2, fp + 0
    lod_ptr r2, r2
    push r2
    mov r2, 1
    push r2
    call get_word
    // Pop saved registers
    // Pop allocated registers
    pop r1
    // End pop
    mov r2, e6
    str_ptr r1, r2
if_stmt_3_check:
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 8
    push r1
    push var_str_4
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 8
    push r1
    mov r1, fp + 4
    lod_ptr r1, r1
    push r1
    call strcmp
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jnz r1, if_stmt_3_success
    jmp if_stmt_3_else
if_stmt_3_success:
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 15
    push r1
    push var_str_5
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 15
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    mov r1, 0    // User-defined inline assembly
    int 0xf    // User-defined inline assembly
    jmp if_stmt_3_done
if_stmt_3_else:
if_stmt_3_done:
if_stmt_6_check:
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 28
    push r1
    push var_str_7
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 28
    push r1
    mov r1, fp + 4
    lod_ptr r1, r1
    push r1
    call strcmp
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jnz r1, if_stmt_6_success
    jmp if_stmt_6_else
if_stmt_6_success:
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 34
    push r1
    push var_str_8
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 34
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 63
    push r1
    push var_str_9
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 63
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jmp while_stmt_1_check
    jmp if_stmt_6_done
if_stmt_6_else:
if_stmt_6_done:
if_stmt_10_check:
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 68
    push r1
    push var_str_11
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 68
    push r1
    mov r1, fp + 4
    lod_ptr r1, r1
    push r1
    call strcmp
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jnz r1, if_stmt_10_success
    jmp if_stmt_10_else
if_stmt_10_success:
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 79
    push r1
    push var_str_12
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 79
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    // Push allocated registers
    // Push allocated registers
    // End push
    mov r1, var_2
    lod_ptr r1, r1
    push r1
    mov r1, 0
    push r1
    mov r1, 0
    push r1
    call readin
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    // Push allocated registers
    // Push allocated registers
    // End push
    mov r1, var_2
    lod_ptr r1, r1
    push r1
    mov r1, 0
    push r1
    call save_buffer
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 103
    push r1
    push var_str_13
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 103
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jmp while_stmt_1_check
    jmp if_stmt_10_done
if_stmt_10_else:
if_stmt_10_done:
if_stmt_14_check:
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 106
    push r1
    push var_str_15
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 106
    push r1
    mov r1, fp + 4
    lod_ptr r1, r1
    push r1
    call strcmp
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jnz r1, if_stmt_14_success
    jmp if_stmt_14_else
if_stmt_14_success:
    mov r1, fp + 111
    // Push allocated registers
    // Push allocated registers
    push r1
    // End push
    mov r2, fp + 0
    lod_ptr r2, r2
    push r2
    mov r2, 2
    push r2
    call get_word
    // Pop saved registers
    // Pop allocated registers
    pop r1
    // End pop
    mov r2, e6
    str_ptr r1, r2
    // Push allocated registers
    // Push allocated registers
    // End push
    mov r1, 64
    push r1
    call malloc
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
if_stmt_16_check:
    // Push allocated registers
    // Push allocated registers
    // End push
    mov r1, fp + 111
    lod_ptr r1, r1
    push r1
    call strlen
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    mov r2, 0
    cmp r1, r1, r2
    jnz r1, if_stmt_16_success
    jmp if_stmt_16_else
if_stmt_16_success:
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 115
    push r1
    push var_str_17
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 115
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    // Push allocated registers
    // Push allocated registers
    // End push
    mov r1, 64
    push r1
    call free
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jmp while_stmt_1_check
    jmp if_stmt_16_done
if_stmt_16_else:
if_stmt_16_done:
    mov r1, fp + 140
    // Push allocated registers
    // Push allocated registers
    push r1
    // End push
    // Push allocated registers
    // Push allocated registers
    push r1
    // End push
    mov r2, fp + 111
    lod_ptr r2, r2
    push r2
    call fntf
    // Pop saved registers
    // Pop allocated registers
    pop r1
    // End pop
    mov r2, e6
    push r2
    call fgetsize
    // Pop saved registers
    // Pop allocated registers
    pop r1
    // End pop
    mov r2, e6
    str32 r1, r2
    mov r1, fp + 144
    // Push allocated registers
    // Push allocated registers
    push r1
    // End push
    // Push allocated registers
    // Push allocated registers
    push r1
    // End push
    mov r2, fp + 111
    lod_ptr r2, r2
    push r2
    call fntf
    // Pop saved registers
    // Pop allocated registers
    pop r1
    // End pop
    mov r2, e6
    push r2
    mov r2, 1
    push r2
    call fopen
    // Pop saved registers
    // Pop allocated registers
    pop r1
    // End pop
    mov r2, e6
    str_ptr r1, r2
    mov r1, fp + 148
    mov r2, fp + 144
    // Above should not be loaded
    lod_ptr r2, r2
    mov r3, 0
    add r2, r2, r3
    lod_ptr r2, r2
    str_ptr r1, r2
if_stmt_18_check:
    mov r1, fp + 148
    lod_ptr r1, r1
    mov r2, 0x00000000
    cmp r1, r1, r2
    jnz r1, if_stmt_18_success
    jmp if_stmt_18_else
if_stmt_18_success:
    // Push allocated registers
    // Push allocated registers
    // End push
    mov r1, 64
    push r1
    call free
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jmp while_stmt_1_check
    jmp if_stmt_18_done
if_stmt_18_else:
if_stmt_18_done:
    mov r1, fp + 152
    // Push allocated registers
    // Push allocated registers
    push r1
    // End push
    mov r2, fp + 140
    lod32 r2, r2
    push r2
    call malloc
    // Pop saved registers
    // Pop allocated registers
    pop r1
    // End pop
    mov r2, e6
    str_ptr r1, r2
    mov r1, fp + 156
    // Push allocated registers
    // Push allocated registers
    push r1
    // End push
    mov r2, fp + 148
    lod_ptr r2, r2
    push r2
    mov r2, fp + 152
    lod_ptr r2, r2
    push r2
    call strcpy
    // Pop saved registers
    // Pop allocated registers
    pop r1
    // End pop
    mov r2, e6
    str_ptr r1, r2
    mov r1, fp + 156
    lod_ptr r1, r1
    mov r2, 0
    str r1, r2
    // Push allocated registers
    // Push allocated registers
    // End push
    mov r1, fp + 152
    lod_ptr r1, r1
    push r1
    call textedit_init
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    // Push allocated registers
    // End push
    mov r1, fp + 111
    lod_ptr r1, r1
    push r1
    call fntf
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    push r1
    mov r1, fp + 152
    lod_ptr r1, r1
    push r1
    call fwrite
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    // Push allocated registers
    // Push allocated registers
    // End push
    mov r1, 64
    push r1
    call free
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    // Push allocated registers
    // Push allocated registers
    // End push
    mov r1, fp + 140
    lod32 r1, r1
    push r1
    call free
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 160
    push r1
    push var_str_19
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 160
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jmp while_stmt_1_check
    jmp if_stmt_14_done
if_stmt_14_else:
if_stmt_14_done:
if_stmt_20_check:
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 163
    push r1
    push var_str_21
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 163
    push r1
    mov r1, fp + 4
    lod_ptr r1, r1
    push r1
    call strcmp
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jnz r1, if_stmt_20_success
    jmp if_stmt_20_else
if_stmt_20_success:
    // Push allocated registers
    // Push allocated registers
    // End push
    call flist
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 169
    push r1
    push var_str_22
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 169
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jmp while_stmt_1_check
    jmp if_stmt_20_done
if_stmt_20_else:
if_stmt_20_done:
if_stmt_23_check:
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 172
    push r1
    push var_str_24
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 172
    push r1
    mov r1, fp + 4
    lod_ptr r1, r1
    push r1
    call strcmp
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jnz r1, if_stmt_23_success
    jmp if_stmt_23_else
if_stmt_23_success:
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 181
    push r1
    push var_str_25
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 181
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    int 0x11    // User-defined inline assembly
    jmp if_stmt_23_done
if_stmt_23_else:
if_stmt_23_done:
if_stmt_26_check:
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 200
    push r1
    push var_str_27
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 200
    push r1
    mov r1, fp + 4
    lod_ptr r1, r1
    push r1
    call strcmp
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jnz r1, if_stmt_26_success
    jmp if_stmt_26_else
if_stmt_26_success:
    mov r1, 4    // User-defined inline assembly
    mov r2, pc    // User-defined inline assembly
    int 0x07    // User-defined inline assembly
    jmp if_stmt_26_done
if_stmt_26_else:
if_stmt_26_done:
if_stmt_28_check:
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 210
    push r1
    push var_str_29
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 210
    push r1
    mov r1, fp + 4
    lod_ptr r1, r1
    push r1
    call strcmp
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jnz r1, if_stmt_28_success
    jmp if_stmt_28_else
if_stmt_28_success:
    // Push allocated registers
    // Push allocated registers
    // End push
    mov r1, 0x40404040
    push r1
    call render_buf
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    // Push allocated registers
    // Push allocated registers
    // End push
    mov r1, 0
    push r1
    mov r1, 0
    push r1
    call video_set_cursor
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jmp while_stmt_1_check
    jmp if_stmt_28_done
if_stmt_28_else:
if_stmt_28_done:
if_stmt_30_check:
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 216
    push r1
    push var_str_31
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 216
    push r1
    mov r1, fp + 4
    lod_ptr r1, r1
    push r1
    call strcmp
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jnz r1, if_stmt_30_success
    jmp if_stmt_30_else
if_stmt_30_success:
    // Push allocated registers
    // Push allocated registers
    // End push
    call load_executable
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jmp while_stmt_1_check
    jmp if_stmt_30_done
if_stmt_30_else:
if_stmt_30_done:
if_stmt_32_check:
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 221
    push r1
    push var_str_33
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 221
    push r1
    mov r1, fp + 4
    lod_ptr r1, r1
    push r1
    call strcmp
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jnz r1, if_stmt_32_success
    jmp if_stmt_32_else
if_stmt_32_success:
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 229
    push r1
    push var_str_34
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 229
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    // Push allocated registers
    // End push
    mov r1, 0x80000026
    lod r1, r1
    push r1
    call atoi
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 245
    push r1
    push var_str_35
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 245
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jmp while_stmt_1_check
    jmp if_stmt_32_done
if_stmt_32_else:
if_stmt_32_done:
if_stmt_36_check:
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 249
    push r1
    push var_str_37
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 249
    push r1
    mov r1, fp + 4
    lod_ptr r1, r1
    push r1
    call strcmp
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jnz r1, if_stmt_36_success
    jmp if_stmt_36_else
if_stmt_36_success:
    // Push allocated registers
    // Push allocated registers
    // End push
    call teststack
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 259
    push r1
    push var_str_38
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 259
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jmp while_stmt_1_check
    jmp if_stmt_36_done
if_stmt_36_else:
if_stmt_36_done:
if_stmt_39_check:
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 307
    push r1
    push var_str_40
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 307
    push r1
    mov r1, fp + 4
    lod_ptr r1, r1
    push r1
    call strcmp
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jnz r1, if_stmt_39_success
    jmp if_stmt_39_else
if_stmt_39_success:
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 312
    push r1
    push var_str_41
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 312
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    // Push allocated registers
    // End push
    mov r1, 0x80000022
    lod r1, r1
    push r1
    call atoi
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 319
    push r1
    push var_str_42
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 319
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    // Push allocated registers
    // End push
    mov r1, 0x80000021
    lod r1, r1
    push r1
    call atoi
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 321
    push r1
    push var_str_43
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 321
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jmp while_stmt_1_check
    jmp if_stmt_39_done
if_stmt_39_else:
if_stmt_39_done:
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 324
    push r1
    push var_str_44
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 324
    push r1
    mov r1, 0xE9
    push r1
    mov r1, 0x00
    push r1
    call puts32
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    // Push allocated registers
    // Push allocated registers
    // End push
    mov r1, fp + 4
    lod_ptr r1, r1
    push r1
    mov r1, 0xE9
    push r1
    mov r1, 0x00
    push r1
    call puts32
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    // Push allocated registers
    // Push allocated registers
    // End push
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 338
    push r1
    push var_str_45
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 338
    push r1
    mov r1, 0xE9
    push r1
    mov r1, 0x00
    push r1
    call puts32
    // Pop saved registers
    // Pop allocated registers
    // End pop
    mov r1, e6
    jmp while_stmt_1_check
while_stmt_1_after:
    jmp .shell_ret
.shell_ret:
    pop e11
    pop fp
    ret

