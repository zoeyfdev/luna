.bits 32
#define _builtin_lcc_basin__cstart 0
.global _cstart
var_53:
    .ptr PROMPTBUF
var_54:
    .ptr renderbuf_loc
var_55:
    .ptr sleep_loc
var_56:
    .ptr malloc_loc
var_57:
    .ptr puts32_loc
var_101:
    .ptr BOOT_IMG
var_105:
    .asciz "NOTEPAD.SYS"
var_106:
    .ptr var_105
var_110:
    .asciz "NOTEPAD.SYS"
var_111:
    .ptr var_110
var_112:
    .asciz "Welcome to "
var_113:
    .ptr var_112
var_114:
    .asciz "Luna"
var_115:
    .ptr var_114
var_116:
    .asciz "OS!\n"
var_117:
    .ptr var_116

_cstart:
    push fp
    mov r12, _builtin_lcc_basin__cstart
    sub fp, fp, r12
    mov r1, var_106
    lod_ptr r1, r2
    mov r7, r2
    push r7
    call fntf
    mov r7, e6
    push r7
    mov r1, 0
    mov r7, r1
    push r7
    call fopen
    mov r11, e6
    mov r1, r11
    mov e8, 0
    add r1, r1, e8
    mov r11, r1
    lod_ptr r11, r11
    mov r1, 0
    mov r5, r1
    cmp r11, r11, r5
    jnz r11, if_stmt_107
    jz r11, after_stmt_109
if_stmt_107:
    mov r1, var_111
    lod_ptr r1, r2
    mov r7, r2
    push r7
    call fntf
    mov r7, e6
    push r7
    mov r1, 256
    mov r7, r1
    push r7
    call fcreate
    mov r4, e6
    jmp after_stmt_109
after_stmt_109:
    mov r1, var_113
    lod_ptr r1, r2
    mov r7, r2
    push r7
    mov r1, 255
    mov r7, r1
    push r7
    mov r1, 0
    mov r7, r1
    push r7
    call puts32
    mov r4, e6
    mov r1, var_115
    lod_ptr r1, r2
    mov r7, r2
    push r7
    mov r1, 95
    mov r7, r1
    push r7
    mov r1, 0
    mov r7, r1
    push r7
    call puts32
    mov r4, e6
    mov r1, var_117
    lod_ptr r1, r2
    mov r7, r2
    push r7
    mov r1, 255
    mov r7, r1
    push r7
    mov r1, 0
    mov r7, r1
    push r7
    call puts32
    mov r4, e6
while_stmt_118_check:
    mov r1, 1
    mov r11, r1
    jnz r11, while_stmt_118_body
    jmp while_stmt_118_after
while_stmt_118_body:
    call shell
    mov r4, e6
    jmp while_stmt_118_check
while_stmt_118_after:
    pop fp
