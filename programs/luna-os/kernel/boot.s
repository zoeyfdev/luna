.bits 32
.global boot

#define _builtin_lcc_basin_boot 67

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

var_7:
    .ptr fzipdecode_loc

var_8:
    .ptr CRASH_SOUND

var_9:
    .ptr BOOT_SOUND

var_10:
    .ptr play_sound_loc

var_11:
    .ptr _builtin_lcc_strcpy32

var_str_0:
    .asciz "Loading resources...\n"

var_str_1:
    .asciz "Loading LunaOS...\n"

var_str_2:
    .asciz "Bootstrapping LUFS...\n"


boot:
    push fp
    mov r12, _builtin_lcc_basin_boot
    sub fp, fp, r12
    mov r1, var_11
    lod_ptr r1, r1
    push r1
    mov r1, 3
    push r1
    call targeted_load
    mov r1, e6
    mov r1, var_5
    lod_ptr r1, r1
    push r1
    mov r1, 2
    push r1
    call targeted_load
    mov r1, e6
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 0
    push r1
    push var_str_0
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
    mov r1, var_9
    lod_ptr r1, r1
    push r1
    mov r1, 43
    push r1
    call targeted_load
    mov r1, e6
    mov r1, var_6
    lod_ptr r1, r1
    push r1
    mov r1, 4
    push r1
    call targeted_load
    mov r1, e6
    mov r1, var_10
    lod_ptr r1, r1
    push r1
    mov r1, 3
    push r1
    call targeted_load
    mov r1, e6
    mov r1, var_2
    lod_ptr r1, r1
    push r1
    mov r1, 2
    push r1
    call targeted_load
    mov r1, e6
    mov r1, var_3
    lod_ptr r1, r1
    push r1
    mov r1, 2
    push r1
    call targeted_load
    mov r1, e6
    mov r1, var_7
    lod_ptr r1, r1
    push r1
    mov r1, 3
    push r1
    call targeted_load
    mov r1, e6
    mov r1, var_4
    lod_ptr r1, r1
    push r1
    mov r1, 2
    push r1
    call targeted_load
    mov r1, e6
    mov r1, var_9
    lod_ptr r1, r1
    push r1
    call fzip_decode
    mov r1, e6
    push r1
    mov r1, 41984
    push r1
    mov r1, 0
    push r1
    call play_sound
    mov r1, e6
    mov r1, var_6
    lod_ptr r1, r1
    push r1
    call fzip_decode
    mov r1, e6
    push r1
    call render_buf
    mov r1, e6
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 23
    push r1
    push var_str_1
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 23
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    mov r1, e6
    mov r1, 0xFF
    push r1
    call linear_sector_load
    mov r1, e6
    // Push allocated registers
    push r1
    // End push
    mov r1, fp + 43
    push r1
    push var_str_2
    call _builtin_lcc_strcpy32
    // Pop allocated registers
    pop r1
    // End pop
    mov r1, fp + 43
    push r1
    mov r1, 0xFF
    push r1
    mov r1, 0x00
    push r1
    call puts32
    mov r1, e6
    call fstrap
    mov r1, e6
    mov r1, 0x30303030
    push r1
    call render_buf
    mov r1, e6
    mov r1, 0
    push r1
    mov r1, 0
    push r1
    call video_set_cursor
    mov r1, e6
    jmp _cstart    // User-defined inline assembly
.boot_ret:
    pop fp
    // Function annotated as 'noreturn'

