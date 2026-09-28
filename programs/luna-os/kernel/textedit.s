.bits 32
.global textedit_init

#define _builtin_lcc_basin_textedit_init 31

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
    .asciz "TextEdit v0.1 - LunaOS\n\n"


textedit_init:
    pop e11
    pop e0
    push fp
    mov r12, _builtin_lcc_basin_textedit_init
    sub fp, fp, r12
    push e11
    mov r0, fp + 0
    str_ptr r0, e0
    call save_graphics_buf
    mov r1, e6
    call video_save_cursor
    mov r1, e6
    mov r1, 0x40404040
    push r1
    call render_buf
    mov r1, e6
    mov r1, 0
    push r1
    mov r1, 0
    push r1
    call video_set_cursor
    mov r1, e6
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 4
    push r1
    push var_str_0
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 4
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    mov r1, e6
    mov r1, fp + 0
    lod_ptr r1, r1
    push r1
    mov r1, 0
    push r1
    mov r1, 0
    push r1
    call readin
    mov r1, e6
    mov r1, 0x30303030
    push r1
    call render_buf
    mov r1, e6
    call video_load_cursor
    mov r1, e6
    jmp .textedit_init_ret
.textedit_init_ret:
    pop e11
    pop fp
    ret

