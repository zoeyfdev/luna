.bits 32
.global _cstart

#define _builtin_lcc_basin__cstart 47

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

var_6:
    .ptr BOOT_IMG

var_str_1:
    .asciz "NOTEPAD.SYS"

var_str_2:
    .asciz "NOTEPAD.SYS"

var_str_3:
    .asciz "Welcome to "

var_str_4:
    .asciz "Luna"

var_str_5:
    .asciz "OS!\n"


_cstart:
    push fp
    mov r12, _builtin_lcc_basin__cstart
    sub fp, fp, r12
if_stmt_0_check:
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 0
    push r1
    push var_str_1
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 0
    push r1
    call fntf
    mov r1, e6
    push r1
    mov r1, 0
    push r1
    call fopen
    mov r1, e6
    // Above should not be loaded
    mov r2, 0
    add r1, r1, r2
    lod_ptr r1, r1
    mov r2, 0x00000000
    cmp r1, r1, r2
    jnz r1, if_stmt_0_success
    jmp if_stmt_0_else
if_stmt_0_success:
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 12
    push r1
    push var_str_2
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 12
    push r1
    call fntf
    mov r1, e6
    push r1
    mov r1, 256
    push r1
    call fcreate
    mov r1, e6
    jmp if_stmt_0_done
if_stmt_0_else:
if_stmt_0_done:
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 24
    push r1
    push var_str_3
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 24
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    mov r1, e6
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 36
    push r1
    push var_str_4
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 36
    push r1
    mov r1, 0x5F
    push r1
    mov r1, 0x00
    push r1
    call puts32
    mov r1, e6
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 41
    push r1
    push var_str_5
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 41
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    mov r1, e6
while_stmt_6_check:
    mov r1, 1
    jnz r1, while_stmt_6_body
    jmp while_stmt_6_after
while_stmt_6_body:
    call shell
    mov r1, e6
    jmp while_stmt_6_check
while_stmt_6_after:
._cstart_ret:
    pop fp
    // Function annotated as 'noreturn'

