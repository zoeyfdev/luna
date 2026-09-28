.bits 32
.global printf
.global tohex
.global pause
.global kernel_panic
.global video_set_cursor
.global get_cursor_x
.global get_cursor_y
.global video_save_cursor
.global video_load_cursor
.global query_drive_inserted
.global reboot
.global load_sector
.global load_executable
.global app_error
.global get_word
.global atoi
.global toint

#define _builtin_lcc_basin_printf 4
#define _builtin_lcc_basin_tohex 11
#define _builtin_lcc_basin_pause 34
#define _builtin_lcc_basin_kernel_panic 139
#define _builtin_lcc_basin_video_set_cursor 8
#define _builtin_lcc_basin_get_cursor_x 0
#define _builtin_lcc_basin_get_cursor_y 0
#define _builtin_lcc_basin_video_save_cursor 0
#define _builtin_lcc_basin_video_load_cursor 0
#define _builtin_lcc_basin_query_drive_inserted 1
#define _builtin_lcc_basin_reboot 0
#define _builtin_lcc_basin_load_sector 9
#define _builtin_lcc_basin_load_executable 112
#define _builtin_lcc_basin_app_error 73
#define _builtin_lcc_basin_get_word 20
#define _builtin_lcc_basin_atoi 20
#define _builtin_lcc_basin_toint 7

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
    .ptr CRASH_SOUND

var_7:
    .ptr BOOT_SOUND

var_8:
    .ptr play_sound_loc

var_9:
    .ptr fzipdecode_loc

var_str_0:
    .asciz "0x"

var_str_1:
    .asciz "\n"

var_str_2:
    .asciz "Press any key to continue...\n\n"

var_str_3:
    .asciz "System error\n\nYour PC ran into an error and needs to\nbe restarted.\n\nPress any key to reboot.\n\n\n"

var_str_4:
    .asciz "Instruction: 0x"

var_str_5:
    .asciz "\n"

var_str_6:
    .asciz "Location: 0x"

var_str_7:
    .asciz "\n"

var_11:
    .dword 0

var_12:
    .dword 0

var_str_9:
    .asciz "Error! "

var_str_10:
    .asciz "Please insert a disc into the DVD\ndrive and try again.\n"

var_str_12:
    .asciz "Error! "

var_str_13:
    .asciz "Invalid executable file format.\n"

var_str_14:
    .asciz "Error! "

var_str_15:
    .asciz "Executable automatically\nterminated due to instruction fault.\n"

var_str_24:
    .asciz "\n"


printf:
    pop e11
    pop e0
    push fp
    mov r12, _builtin_lcc_basin_printf
    sub fp, fp, r12
    push e11
    mov r0, fp + 0
    str_ptr r0, e0
    mov r1, fp + 0
    lod_ptr r1, r1
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    mov r1, e6
.printf_ret:
    pop e11
    pop fp
    ret

tohex:
    pop e11
    pop e1
    pop e0
    push fp
    mov r12, _builtin_lcc_basin_tohex
    sub fp, fp, r12
    push e11
    mov r0, fp + 0
    str32 r0, e0
    mov r0, fp + 4
    str r0, e1
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 5
    push r1
    push var_str_0
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 5
    push r1
    mov r1, 255
    push r1
    mov r1, 0
    push r1
    call puts32
    mov r1, e6
    mov r1, fp + 0
    lod32 r1, r1
    push r1
    mov r1, fp + 4
    lod r1, r1
    push r1
    mov r1, 11
    push r1
    call malloc
    mov r1, e6
    push r1
    call itoa
    mov r1, e6
    push r1
    mov r1, 255
    push r1
    mov r1, 0
    push r1
    call puts32
    mov r1, e6
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 8
    push r1
    push var_str_1
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 8
    push r1
    mov r1, 255
    push r1
    mov r1, 0
    push r1
    call puts32
    mov r1, e6
    mov r1, 11
    push r1
    call free
    mov r1, e6
.tohex_ret:
    pop e11
    pop fp
    ret

pause:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_pause
    sub fp, fp, r12
    push e11
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 0
    push r1
    push var_str_2
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 0
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    mov r1, e6
    mov r1, fp + 33
    // Push allocated registers
    push r1
    // End push
    call wait_for_key
    // Pop allocated registers
    pop r1
    // End pop
    mov r2, e6
    str r1, r2
    mov r1, fp + 33
    lod r1, r1
    mov e6, r1
    jmp .pause_ret
.pause_ret:
    pop e11
    pop fp
    ret

kernel_panic:
    push fp
    mov r12, _builtin_lcc_basin_kernel_panic
    sub fp, fp, r12
    push r2    // User-defined inline assembly
    push r1    // User-defined inline assembly
    mov r1, var_6
    lod_ptr r1, r1
    push r1
    call fzip_decode
    mov r1, e6
    push r1
    mov r1, 164352
    push r1
    mov r1, 0
    push r1
    call play_sound
    mov r1, e6
    mov r1, 0x80808080
    push r1
    call screen_fill
    mov r1, e6
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 0
    push r1
    push var_str_3
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 0
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x80
    push r1
    call puts32
    mov r1, e6
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 104
    push r1
    push var_str_4
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 104
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x80
    push r1
    call puts32
    mov r1, e6
    pop e9    // User-defined inline assembly
    mov r1, e9
    push r1
    mov r1, 1
    push r1
    mov r1, 11
    push r1
    call malloc
    mov r1, e6
    push r1
    call itoa
    mov r1, e6
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x80
    push r1
    call puts32
    mov r1, e6
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 120
    push r1
    push var_str_5
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 120
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x80
    push r1
    call puts32
    mov r1, e6
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 123
    push r1
    push var_str_6
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 123
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x80
    push r1
    call puts32
    mov r1, e6
    pop e9    // User-defined inline assembly
    mov r1, e9
    push r1
    mov r1, 1
    push r1
    mov r1, 11
    push r1
    call malloc
    mov r1, e6
    push r1
    call itoa
    mov r1, e6
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x80
    push r1
    call puts32
    mov r1, e6
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 136
    push r1
    push var_str_7
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 136
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x80
    push r1
    call puts32
    mov r1, e6
    call wait_for_key
    mov r1, e6
    int 0x10    // User-defined inline assembly
    int 0xf    // User-defined inline assembly
.kernel_panic_ret:
    pop fp
    // Function annotated as 'noreturn'

video_set_cursor:
    pop e11
    pop e1
    pop e0
    push fp
    mov r12, _builtin_lcc_basin_video_set_cursor
    sub fp, fp, r12
    push e11
    mov r0, fp + 0
    str32 r0, e0
    mov r0, fp + 4
    str32 r0, e1
    mov r1, e0    // User-defined inline assembly
    mov r2, e1    // User-defined inline assembly
    int 0x0c    // User-defined inline assembly
.video_set_cursor_ret:
    pop e11
    pop fp
    ret

get_cursor_x:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_get_cursor_x
    sub fp, fp, r12
    push e11
    int 0xe    // User-defined inline assembly
    mov e6, r1    // User-defined inline assembly
.get_cursor_x_ret:
    pop e11
    pop fp
    ret

get_cursor_y:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_get_cursor_y
    sub fp, fp, r12
    push e11
    int 0xe    // User-defined inline assembly
    mov e6, r2    // User-defined inline assembly
.get_cursor_y_ret:
    pop e11
    pop fp
    ret

video_save_cursor:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_video_save_cursor
    sub fp, fp, r12
    push e11
    mov r1, var_11
    // Push allocated registers
    push r1
    // End push
    call get_cursor_x
    // Pop allocated registers
    pop r1
    // End pop
    mov r2, e6
    str32 r1, r2
    mov r1, var_12
    // Push allocated registers
    push r1
    // End push
    call get_cursor_y
    // Pop allocated registers
    pop r1
    // End pop
    mov r2, e6
    str32 r1, r2
.video_save_cursor_ret:
    pop e11
    pop fp
    ret

video_load_cursor:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_video_load_cursor
    sub fp, fp, r12
    push e11
    mov r1, var_11
    lod32 r1, r1
    push r1
    mov r1, var_12
    lod32 r1, r1
    push r1
    call video_set_cursor
    mov r1, e6
.video_load_cursor_ret:
    pop e11
    pop fp
    ret

query_drive_inserted:
    pop e11
    pop e0
    push fp
    mov r12, _builtin_lcc_basin_query_drive_inserted
    sub fp, fp, r12
    push e11
    mov r0, fp + 0
    str r0, e0
    mov r1, e0    // User-defined inline assembly
    int 0x3    // User-defined inline assembly
    mov e6, r1    // User-defined inline assembly
.query_drive_inserted_ret:
    pop e11
    pop fp
    ret

reboot:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_reboot
    sub fp, fp, r12
    push e11
    int 0x10    // User-defined inline assembly
    int 0xf    // User-defined inline assembly
.reboot_ret:
    pop e11
    pop fp
    ret

load_sector:
    pop e11
    pop e2
    pop e1
    pop e0
    push fp
    mov r12, _builtin_lcc_basin_load_sector
    sub fp, fp, r12
    push e11
    mov r0, fp + 0
    str r0, e0
    mov r0, fp + 1
    str_ptr r0, e1
    mov r0, fp + 5
    str32 r0, e2
    mov r2, e0    // User-defined inline assembly
    mov r1, e1    // User-defined inline assembly
    mov r3, e2    // User-defined inline assembly
    int 0x0b    // User-defined inline assembly
.load_sector_ret:
    pop e11
    pop fp
    ret

load_executable:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_load_executable
    sub fp, fp, r12
    push e11
if_stmt_8_check:
    mov r1, 2
    push r1
    call query_drive_inserted
    mov r1, e6
    mov r2, 0
    cmp r1, r1, r2
    jnz r1, if_stmt_8_success
    jmp if_stmt_8_else
if_stmt_8_success:
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 0
    push r1
    push var_str_9
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 0
    push r1
    mov r1, 0x80
    push r1
    mov r1, 0x00
    push r1
    call puts32
    mov r1, e6
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 8
    push r1
    push var_str_10
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 8
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    mov r1, e6
    jmp .load_executable_ret
    jmp if_stmt_8_done
if_stmt_8_else:
if_stmt_8_done:
    mov r1, fp + 66
    // Push allocated registers
    push r1
    // End push
    call ASLR_generate_address
    // Pop allocated registers
    pop r1
    // End pop
    mov r2, e6
    str_ptr r1, r2
    mov r2, fp + 66
    mov r3, fp + 66
    lod_ptr r3, r3
    mov r1, r3
    inc r1
    str_ptr r2, r1
    mov r1, 2
    push r1
    mov r1, fp + 66
    lod_ptr r1, r1
    mov r2, 512
    div r1, r1, r2
    push r1
    mov r1, 0
    push r1
    call load_sector
    mov r1, e6
    mov r1, 2
    push r1
    mov r1, fp + 66
    lod_ptr r1, r1
    mov r2, 512
    div r1, r1, r2
    mov r2, 1
    add r1, r1, r2
    push r1
    mov r1, 1
    push r1
    call load_sector
    mov r1, e6
    mov r1, 2
    push r1
    mov r1, fp + 66
    lod_ptr r1, r1
    mov r2, 512
    div r1, r1, r2
    mov r2, 2
    add r1, r1, r2
    push r1
    mov r1, 2
    push r1
    call load_sector
    mov r1, e6
if_stmt_11_check:
    mov r1, fp + 66
    lod_ptr r1, r1
    lod32 r1, r1
    mov r2, 0x4C325049
    cmp r1, r1, r2
    mov r3, 1
    xor r1, r1, r3
    jnz r1, if_stmt_11_success
    jmp if_stmt_11_else
if_stmt_11_success:
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 70
    push r1
    push var_str_12
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 70
    push r1
    mov r1, 0x80
    push r1
    mov r1, 0x00
    push r1
    call puts32
    mov r1, e6
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 78
    push r1
    push var_str_13
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 78
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    mov r1, e6
    jmp .load_executable_ret
    jmp if_stmt_11_done
if_stmt_11_else:
if_stmt_11_done:
    mov r1, fp + 66
    lod_ptr r1, r1
    push r1
    call lexec_core
    mov r1, e6
.load_executable_ret:
    pop e11
    pop fp
    ret

app_error:
    push fp
    mov r12, _builtin_lcc_basin_app_error
    sub fp, fp, r12
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 0
    push r1
    push var_str_14
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 0
    push r1
    mov r1, 0x80
    push r1
    mov r1, 0x00
    push r1
    call puts32
    mov r1, e6
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 8
    push r1
    push var_str_15
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 8
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    mov r1, e6
    jmp lexec_done
.app_error_ret:
    pop fp
    // Function annotated as 'noreturn'

get_word:
    pop e11
    pop e1
    pop e0
    push fp
    mov r12, _builtin_lcc_basin_get_word
    sub fp, fp, r12
    push e11
    mov r0, fp + 0
    str_ptr r0, e0
    mov r0, fp + 4
    str32 r0, e1
    mov r1, fp + 8
    // Push allocated registers
    push r1
    // End push
    mov r2, 1024
    push r2
    call malloc
    // Pop allocated registers
    pop r1
    // End pop
    mov r2, e6
    str_ptr r1, r2
    mov r1, fp + 12
    mov r2, fp + 8
    lod_ptr r2, r2
    str_ptr r1, r2
    mov r1, fp + 16
    mov r2, 1
    str32 r1, r2
if_stmt_16_check:
    mov r1, fp + 0
    lod_ptr r1, r1
    push r1
    call strlen
    mov r1, e6
    mov r2, 0
    cmp r1, r1, r2
    jnz r1, if_stmt_16_success
    jmp if_stmt_16_else
if_stmt_16_success:
    mov r1, 1024
    push r1
    call free
    mov r1, e6
    mov r1, fp + 8
    lod_ptr r1, r1
    mov r2, 0
    str r1, r2
    mov r1, fp + 12
    lod_ptr r1, r1
    mov e6, r1
    jmp .get_word_ret
    jmp if_stmt_16_done
if_stmt_16_else:
if_stmt_16_done:
while_stmt_17_check:
    mov r1, fp + 0
    lod_ptr r1, r1
    lod r1, r1
    mov r2, 0x00
    cmp r1, r1, r2
    mov r3, 1
    xor r1, r1, r3
    jnz r1, while_stmt_17_body
    jmp while_stmt_17_after
while_stmt_17_body:
if_stmt_18_check:
    mov r1, fp + 0
    lod_ptr r1, r1
    lod r1, r1
    mov r2, 0x20
    cmp r1, r1, r2
    jnz r1, if_stmt_18_success
    jmp if_stmt_18_else
if_stmt_18_success:
if_stmt_19_check:
    mov r1, fp + 16
    lod32 r1, r1
    mov r2, fp + 4
    lod32 r2, r2
    cmp r1, r1, r2
    jnz r1, if_stmt_19_success
    jmp if_stmt_19_else
if_stmt_19_success:
    jmp while_stmt_17_after
    jmp if_stmt_19_done
if_stmt_19_else:
    mov r2, fp + 16
    mov r3, fp + 16
    lod32 r3, r3
    mov r1, r3
    inc r1
    str32 r2, r1
    mov r2, fp + 0
    mov r3, fp + 0
    lod_ptr r3, r3
    mov r1, r3
    inc r1
    str_ptr r2, r1
if_stmt_19_done:
    jmp if_stmt_18_done
if_stmt_18_else:
if_stmt_18_done:
if_stmt_20_check:
    mov r1, fp + 16
    lod32 r1, r1
    mov r2, fp + 4
    lod32 r2, r2
    cmp r1, r1, r2
    jnz r1, if_stmt_20_success
    jmp if_stmt_20_else
if_stmt_20_success:
    mov r1, fp + 8
    lod_ptr r1, r1
    mov r2, fp + 0
    lod_ptr r2, r2
    lod r2, r2
    str r1, r2
    mov r2, fp + 8
    mov r3, fp + 8
    lod_ptr r3, r3
    mov r1, r3
    inc r1
    str_ptr r2, r1
    jmp if_stmt_20_done
if_stmt_20_else:
if_stmt_20_done:
    mov r2, fp + 0
    mov r3, fp + 0
    lod_ptr r3, r3
    mov r1, r3
    inc r1
    str_ptr r2, r1
    jmp while_stmt_17_check
while_stmt_17_after:
    mov r1, fp + 8
    lod_ptr r1, r1
    mov r2, 0
    str r1, r2
    mov r1, 1024
    push r1
    call free
    mov r1, e6
    mov r1, fp + 12
    lod_ptr r1, r1
    mov e6, r1
    jmp .get_word_ret
.get_word_ret:
    pop e11
    pop fp
    ret

atoi:
    pop e11
    pop e0
    push fp
    mov r12, _builtin_lcc_basin_atoi
    sub fp, fp, r12
    push e11
    mov r0, fp + 0
    str32 r0, e0
    mov r1, fp + 4
    // Push allocated registers
    push r1
    // End push
    mov r2, 32
    push r2
    call malloc
    // Pop allocated registers
    pop r1
    // End pop
    mov r2, e6
    str_ptr r1, r2
    mov r1, fp + 8
    // Push allocated registers
    push r1
    // End push
    mov r2, 32
    push r2
    call malloc
    // Pop allocated registers
    pop r1
    // End pop
    mov r2, e6
    str_ptr r1, r2
    mov r1, fp + 12
    mov r2, 0
    str32 r1, r2
if_stmt_21_check:
    mov r1, fp + 0
    lod32 r1, r1
    mov r2, 0
    cmp r1, r1, r2
    jnz r1, if_stmt_21_success
    jmp if_stmt_21_else
if_stmt_21_success:
    mov r1, fp + 4
    lod_ptr r1, r1
    mov r2, 0x30
    str r1, r2
    mov r1, fp + 4
    lod_ptr r1, r1
    mov e6, r1
    jmp .atoi_ret
    jmp if_stmt_21_done
if_stmt_21_else:
if_stmt_21_done:
while_stmt_22_check:
    mov r1, fp + 0
    lod32 r1, r1
    mov r2, 0
    igt r1, r1, r2
    jnz r1, while_stmt_22_body
    jmp while_stmt_22_after
while_stmt_22_body:
    mov r1, fp + 8
    lod_ptr r1, r1
    mov r2, fp + 12
    lod32 r2, r2
    add r1, r1, r2
    // Load ignored
    mov r2, 0x30
    mov r3, fp + 0
    lod32 r3, r3
    mov r4, 10
    mod r3, r3, r4
    add r2, r2, r3
    str r1, r2
    mov r2, fp + 12
    mov r3, fp + 12
    lod32 r3, r3
    mov r1, r3
    inc r1
    str32 r2, r1
    mov r1, fp + 0
    mov r2, fp + 0
    lod32 r2, r2
    mov r3, 10
    div r2, r2, r3
    str32 r1, r2
    jmp while_stmt_22_check
while_stmt_22_after:
    mov r1, fp + 16
    mov r2, 0
    str32 r1, r2
while_stmt_23_check:
    mov r1, fp + 12
    lod32 r1, r1
    mov r2, 0
    igt r1, r1, r2
    jnz r1, while_stmt_23_body
    jmp while_stmt_23_after
while_stmt_23_body:
    mov r2, fp + 12
    mov r3, fp + 12
    lod32 r3, r3
    mov r1, r3
    dec r1
    str32 r2, r1
    mov r1, fp + 4
    lod_ptr r1, r1
    mov r3, fp + 16
    mov r4, fp + 16
    lod32 r4, r4
    mov r2, r4
    inc r2
    str32 r3, r2
    add r1, r1, r4
    // Load ignored
    mov r2, fp + 8
    lod_ptr r2, r2
    mov r3, fp + 12
    lod32 r3, r3
    add r2, r2, r3
    lod r2, r2
    str r1, r2
    jmp while_stmt_23_check
while_stmt_23_after:
    mov r1, fp + 4
    lod_ptr r1, r1
    mov r2, fp + 16
    lod32 r2, r2
    add r1, r1, r2
    // Load ignored
    mov r2, 0x00
    str r1, r2
    mov r1, fp + 4
    lod_ptr r1, r1
    mov e6, r1
    jmp .atoi_ret
.atoi_ret:
    pop e11
    pop fp
    ret

toint:
    pop e11
    pop e0
    push fp
    mov r12, _builtin_lcc_basin_toint
    sub fp, fp, r12
    push e11
    mov r0, fp + 0
    str32 r0, e0
    mov r1, fp + 0
    lod32 r1, r1
    push r1
    call atoi
    mov r1, e6
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
    mov r1, fp + 4
    push r1
    push var_str_24
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
.toint_ret:
    pop e11
    pop fp
    ret

