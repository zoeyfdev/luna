.bits 32
#define _builtin_lcc_basin_fstrap 20
#define _builtin_lcc_basin_fwrite 20
#define _builtin_lcc_basin_flist 12
#define _builtin_lcc_basin_fgetsize 12
#define _builtin_lcc_basin_fopen 17
#define _builtin_lcc_basin_find_file 16
#define _builtin_lcc_basin_fcreate 28
#define _builtin_lcc_basin_fntf 21
#define _builtin_lcc_basin_ffnt 13
.global fstrap
.global fwrite
.global flist
.global fgetsize
.global fopen
.global find_file
.global fcreate
.global fntf
.global ffnt
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
var_168:
    .asciz "File '"
var_169:
    .ptr var_168
var_170:
    .asciz "' not found!\n"
var_171:
    .ptr var_170
var_188:
    .asciz "\n"
var_189:
    .ptr var_188

ffnt:
    pop e11
    pop e0
    push fp
    mov r12, _builtin_lcc_basin_ffnt
    sub fp, fp, r12
    push e11
    mov r1, fp + 0
    str_ptr r1, e0
    mov r1, fp + 0 // Variable name: filename, internal: fp + 0
    lod_ptr r1, r2
    mov r7, r2
    push r7
    call strlen
    mov r7, e6
    push r7
    call malloc
    mov r4, e6
    mov r7, fp + 4
    str_ptr r7, r4
    mov r1, fp + 4 // Variable name: buffer, internal: fp + 4
    lod_ptr r1, r2
    mov r4, r2
    mov r7, fp + 8
    str_ptr r7, r4
    mov r1, 0
    mov r4, r1
    mov r7, fp + 12
    str r7, r4
while_stmt_88_check:
    mov r1, fp + 0 // Variable name: filename, internal: fp + 0
    lod_ptr r1, r2
    lod r2, r2
    mov r11, r2
    mov r1, 0
    mov r5, r1
    cmp r11, r11, r5
    jz r11, while_stmt_88_body
    jmp while_stmt_88_after
while_stmt_88_body:
    mov r1, fp + 0 // Variable name: filename, internal: fp + 0
    lod_ptr r1, r2
    lod r2, r2
    mov r11, r2
    mov r1, 32
    mov r5, r1
    cmp r11, r11, r5
    jz r11, if_stmt_92
    jnz r11, else_stmt_93
if_stmt_92:
    mov r1, fp + 4 // Variable name: buffer, internal: fp + 4
    mov r2, r1
    lod_ptr r2, r2
    mov r4, r2
    mov r1, fp + 0 // Variable name: filename, internal: fp + 0
    lod_ptr r1, r2
    lod r2, r2
    mov r5, r2
    str r4, r5
    mov r1, fp + 4 // Variable name: buffer, internal: fp + 4
    mov r2, r1
    mov r4, r2
    mov r1, fp + 4 // Variable name: buffer, internal: fp + 4
    lod_ptr r1, r2
    mov r0, r2
    mov r3, r4
    mov r4, r0
    mov r5, r0
    inc r0
    str_ptr r3, r0
    jmp after_stmt_94
else_stmt_93:
    mov r1, fp + 12 // Variable name: seen, internal: fp + 12
    lod r1, r2
    mov r11, r2
    mov r1, 0
    mov r5, r1
    cmp r11, r11, r5
    jnz r11, if_stmt_97
    jz r11, after_stmt_99
if_stmt_97:
    mov r1, fp + 12 // Variable name: seen, internal: fp + 12
    mov r2, r1
    mov r4, r2
    mov r1, 1
    mov r5, r1
    str r4, r5
    mov r1, fp + 4 // Variable name: buffer, internal: fp + 4
    mov r2, r1
    lod_ptr r2, r2
    mov r4, r2
    mov r1, 46
    mov r5, r1
    str r4, r5
    mov r1, fp + 4 // Variable name: buffer, internal: fp + 4
    mov r2, r1
    mov r4, r2
    mov r1, fp + 4 // Variable name: buffer, internal: fp + 4
    lod_ptr r1, r2
    mov r0, r2
    mov r3, r4
    mov r4, r0
    mov r5, r0
    inc r0
    str_ptr r3, r0
    jmp after_stmt_99
after_stmt_99:
after_stmt_94:
    mov r1, fp + 0 // Variable name: filename, internal: fp + 0
    mov r2, r1
    mov r4, r2
    mov r1, fp + 0 // Variable name: filename, internal: fp + 0
    lod_ptr r1, r2
    mov r0, r2
    mov r3, r4
    mov r4, r0
    mov r5, r0
    inc r0
    str_ptr r3, r0
    jmp while_stmt_88_check
while_stmt_88_after:
    mov r1, fp + 8 // Variable name: bufptr, internal: fp + 8
    lod_ptr r1, r2
    mov e6, r2
    pop e11
    pop fp
    ret
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
    mov r1, fp + 0
    str_ptr r1, e0
    mov r1, 16
    mov r7, r1
    push r7
    call malloc
    mov r4, e6
    mov r7, fp + 4
    str_ptr r7, r4
    mov r1, fp + 4 // Variable name: buffer, internal: fp + 4
    lod_ptr r1, r2
    mov r4, r2
    mov r7, fp + 8
    str_ptr r7, r4
    mov r1, 0
    mov r4, r1
    mov r7, fp + 12
    str r7, r4
    mov r1, 0
    mov r4, r1
    mov r7, fp + 13
    str16 r7, r4
while_stmt_107_check:
    mov r1, fp + 0 // Variable name: name, internal: fp + 0
    lod_ptr r1, r2
    lod r2, r2
    mov r11, r2
    mov r1, 0
    mov r5, r1
    cmp r11, r11, r5
    jz r11, while_stmt_107_body
    jmp while_stmt_107_after
while_stmt_107_body:
    mov r1, fp + 0 // Variable name: name, internal: fp + 0
    lod_ptr r1, r2
    lod r2, r2
    mov r11, r2
    mov r1, 46
    mov r5, r1
    cmp r11, r11, r5
    jnz r11, if_stmt_111
    jz r11, else_stmt_112
if_stmt_111:
    mov r1, fp + 12 // Variable name: seen, internal: fp + 12
    lod r1, r2
    mov r11, r2
    mov r1, 0
    mov r5, r1
    cmp r11, r11, r5
    jnz r11, if_stmt_116
    jz r11, after_stmt_118
if_stmt_116:
    mov r1, fp + 12 // Variable name: seen, internal: fp + 12
    mov r2, r1
    mov r4, r2
    mov r1, 1
    mov r5, r1
    str r4, r5
    mov r1, 12
    mov r4, r1
    mov r7, fp + 15
    str16 r7, r4
    mov r1, fp + 15 // Variable name: initial, internal: fp + 15
    lod16 r1, r2
    mov e9, r2
    mov r1, fp + 13 // Variable name: copied, internal: fp + 13
    lod16 r1, r2
    mov e7, r2
    sub e10, e9, e7
    mov r2, e10
    mov r4, r2
    mov r7, fp + 17
    str16 r7, r4
    mov r1, 0
    mov r4, r1
    mov r7, fp + 19
    str16 r7, r4
for_stmt_121_check:
    mov r1, fp + 19 // Variable name: i, internal: fp + 19
    lod16 r1, r2
    mov r11, r2
    mov r1, fp + 17 // Variable name: toput, internal: fp + 17
    lod16 r1, r2
    mov r5, r2
    ilt r11, r11, r5
    jz r11, for_stmt_121_after
    mov r1, fp + 4 // Variable name: buffer, internal: fp + 4
    mov r2, r1
    lod_ptr r2, r2
    mov r4, r2
    mov r1, 32
    mov r5, r1
    str r4, r5
    mov r1, fp + 4 // Variable name: buffer, internal: fp + 4
    mov r2, r1
    mov r4, r2
    mov r1, fp + 4 // Variable name: buffer, internal: fp + 4
    lod_ptr r1, r2
    mov r0, r2
    mov r3, r4
    mov r4, r0
    mov r5, r0
    inc r0
    str_ptr r3, r0
    mov r1, fp + 19 // Variable name: i, internal: fp + 19
    mov r2, r1
    mov r4, r2
    mov r1, fp + 19 // Variable name: i, internal: fp + 19
    lod16 r1, r2
    mov r0, r2
    mov r3, r4
    mov r4, r0
    mov r5, r0
    inc r0
    str16 r3, r0
    jmp for_stmt_121_check
for_stmt_121_after:
    mov r1, fp + 0 // Variable name: name, internal: fp + 0
    mov r2, r1
    mov r4, r2
    mov r1, fp + 0 // Variable name: name, internal: fp + 0
    lod_ptr r1, r2
    mov r0, r2
    mov r3, r4
    mov r4, r0
    mov r5, r0
    inc r0
    str_ptr r3, r0
    jmp after_stmt_118
after_stmt_118:
    jmp after_stmt_113
else_stmt_112:
    mov r1, fp + 4 // Variable name: buffer, internal: fp + 4
    mov r2, r1
    lod_ptr r2, r2
    mov r4, r2
    mov r1, fp + 0 // Variable name: name, internal: fp + 0
    lod_ptr r1, r2
    lod r2, r2
    mov r5, r2
    str r4, r5
    mov r1, fp + 13 // Variable name: copied, internal: fp + 13
    mov r2, r1
    mov r4, r2
    mov r1, fp + 13 // Variable name: copied, internal: fp + 13
    lod16 r1, r2
    mov r0, r2
    mov r3, r4
    mov r4, r0
    mov r5, r0
    inc r0
    str16 r3, r0
    mov r1, fp + 0 // Variable name: name, internal: fp + 0
    mov r2, r1
    mov r4, r2
    mov r1, fp + 0 // Variable name: name, internal: fp + 0
    lod_ptr r1, r2
    mov r0, r2
    mov r3, r4
    mov r4, r0
    mov r5, r0
    inc r0
    str_ptr r3, r0
    mov r1, fp + 4 // Variable name: buffer, internal: fp + 4
    mov r2, r1
    mov r4, r2
    mov r1, fp + 4 // Variable name: buffer, internal: fp + 4
    lod_ptr r1, r2
    mov r0, r2
    mov r3, r4
    mov r4, r0
    mov r5, r0
    inc r0
    str_ptr r3, r0
after_stmt_113:
    jmp while_stmt_107_check
while_stmt_107_after:
    mov r1, fp + 8 // Variable name: ogbufptr, internal: fp + 8
    lod_ptr r1, r2
    mov e6, r2
    pop e11
    pop fp
    ret
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
    mov r1, fp + 0
    str_ptr r1, e0
    mov r1, fp + 4
    str32 r1, e1
    mov r1, 1564
    mov r4, r1
    mov r2, r1
    mov r4, r2
    mov r7, fp + 8
    str_ptr r7, r4
    mov r1, fp + 8 // Variable name: nfp, internal: fp + 8
    lod_ptr r1, r2
    lod_ptr r2, r2
    mov r4, r2
    mov r7, fp + 12
    str_ptr r7, r4
    mov r1, fp + 12 // Variable name: nfl, internal: fp + 12
    mov r2, r1
    lod_ptr r2, r2
    mov r4, r2
    mov r1, 1279677254
    mov r5, r1
    str32 r4, r5
    mov r1, fp + 12 // Variable name: nfl, internal: fp + 12
    mov r2, r1
    mov r4, r2
    mov r1, fp + 12 // Variable name: nfl, internal: fp + 12
    lod_ptr r1, r2
    mov e9, r2
    mov e7, 4
    add e10, e9, e7
    mov r2, e10
    mov r5, r2
    str_ptr r4, r5
    mov r1, fp + 0 // Variable name: name, internal: fp + 0
    lod_ptr r1, r2
    mov r7, r2
    push r7
    mov r1, fp + 12 // Variable name: nfl, internal: fp + 12
    lod_ptr r1, r2
    mov r7, r2
    push r7
    call strcpy
    mov r4, e6
    mov r1, fp + 0 // Variable name: name, internal: fp + 0
    lod_ptr r1, r2
    mov r7, r2
    push r7
    call strlen
    mov r4, e6
    mov r7, fp + 16
    str32 r7, r4
    mov r1, fp + 12 // Variable name: nfl, internal: fp + 12
    mov r2, r1
    mov r4, r2
    mov r1, fp + 12 // Variable name: nfl, internal: fp + 12
    lod_ptr r1, r2
    mov e9, r2
    mov r1, fp + 16 // Variable name: name_len, internal: fp + 16
    lod32 r1, r2
    mov e7, r2
    add e10, e9, e7
    mov r2, e10
    mov r5, r2
    str_ptr r4, r5
    mov r1, fp + 12 // Variable name: nfl, internal: fp + 12
    mov r2, r1
    lod_ptr r2, r2
    mov r4, r2
    mov r1, fp + 4 // Variable name: size, internal: fp + 4
    lod32 r1, r2
    mov r5, r2
    str32 r4, r5
    mov r1, fp + 12 // Variable name: nfl, internal: fp + 12
    mov r2, r1
    mov r4, r2
    mov r1, fp + 12 // Variable name: nfl, internal: fp + 12
    lod_ptr r1, r2
    mov e9, r2
    mov e7, 4
    add e10, e9, e7
    mov r2, e10
    mov r5, r2
    str_ptr r4, r5
    mov r1, 0
    mov r4, r1
    mov r7, fp + 20
    str32 r7, r4
for_stmt_131_check:
    mov r1, fp + 20 // Variable name: i, internal: fp + 20
    lod32 r1, r2
    mov r11, r2
    mov r1, fp + 4 // Variable name: size, internal: fp + 4
    lod32 r1, r2
    mov r5, r2
    ilt r11, r11, r5
    jz r11, for_stmt_131_after
    mov r1, fp + 12 // Variable name: nfl, internal: fp + 12
    mov r2, r1
    lod_ptr r2, r2
    mov r4, r2
    mov r1, 0
    mov r5, r1
    str32 r4, r5
    mov r1, fp + 12 // Variable name: nfl, internal: fp + 12
    mov r2, r1
    mov r4, r2
    mov r1, fp + 12 // Variable name: nfl, internal: fp + 12
    lod_ptr r1, r2
    mov r0, r2
    mov r3, r4
    mov r4, r0
    mov r5, r0
    inc r0
    str_ptr r3, r0
    mov r1, fp + 20 // Variable name: i, internal: fp + 20
    mov r2, r1
    mov r4, r2
    mov r1, fp + 20 // Variable name: i, internal: fp + 20
    lod32 r1, r2
    mov r0, r2
    mov r3, r4
    mov r4, r0
    mov r5, r0
    inc r0
    str32 r3, r0
    jmp for_stmt_131_check
for_stmt_131_after:
    mov r1, fp + 12 // Variable name: nfl, internal: fp + 12
    lod_ptr r1, r2
    mov e9, r2
    mov e7, 512
    div e10, e9, e7
    mov r2, e10
    mov r4, r2
    mov r7, fp + 24
    str32 r7, r4
    mov r1, fp + 24 // Variable name: sector, internal: fp + 24
    lod32 r1, r2
    mov r7, r2
    push r7
    call save_sector
    mov r4, e6
    mov r1, fp + 8 // Variable name: nfp, internal: fp + 8
    mov r2, r1
    lod_ptr r2, r2
    mov r4, r2
    mov r1, fp + 8 // Variable name: nfp, internal: fp + 8
    lod_ptr r1, r2
    lod_ptr r2, r2
    mov e9, r2
    mov r1, fp + 4 // Variable name: size, internal: fp + 4
    lod32 r1, r2
    mov e7, r2
    add e10, e9, e7
    mov r2, e10
    mov r5, r2
    str_ptr r4, r5
    mov r1, 3
    mov r7, r1
    push r7
    call save_sector
    mov r4, e6
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
    mov r1, fp + 0
    str_ptr r1, e0
    mov r1, 1560
    mov r4, r1
    mov r2, r1
    mov r4, r2
    mov r7, fp + 4
    str_ptr r7, r4
    mov r1, fp + 4 // Variable name: fsp, internal: fp + 4
    lod_ptr r1, r2
    lod_ptr r2, r2
    mov r4, r2
    mov r7, fp + 8
    str_ptr r7, r4
while_stmt_140_check:
    mov r1, 1
    mov r11, r1
    jnz r11, while_stmt_140_body
    jmp while_stmt_140_after
while_stmt_140_body:
    mov r1, fp + 8 // Variable name: fp, internal: fp + 8
    lod_ptr r1, r2
    lod32 r2, r2
    mov r11, r2
    mov r1, 1279677254
    mov r5, r1
    cmp r11, r11, r5
    jz r11, if_stmt_144
    jnz r11, after_stmt_146
if_stmt_144:
    jmp while_stmt_140_after
    jmp after_stmt_146
after_stmt_146:
    mov r1, fp + 8 // Variable name: fp, internal: fp + 8
    mov r2, r1
    mov r4, r2
    mov r1, fp + 8 // Variable name: fp, internal: fp + 8
    lod_ptr r1, r2
    mov e9, r2
    mov e7, 4
    add e10, e9, e7
    mov r2, e10
    mov r5, r2
    str_ptr r4, r5
    mov r1, fp + 0 // Variable name: name, internal: fp + 0
    lod_ptr r1, r2
    mov r7, r2
    push r7
    mov r1, fp + 8 // Variable name: fp, internal: fp + 8
    lod_ptr r1, r2
    mov r7, r2
    push r7
    call strcmp
    mov r11, e6
    mov r1, 1
    mov r5, r1
    cmp r11, r11, r5
    jnz r11, if_stmt_149
    jz r11, else_stmt_150
if_stmt_149:
    mov r1, fp + 8 // Variable name: fp, internal: fp + 8
    mov r2, r1
    mov r4, r2
    mov r1, fp + 8 // Variable name: fp, internal: fp + 8
    lod_ptr r1, r2
    mov e9, r2
    mov e7, 20
    add e10, e9, e7
    mov r2, e10
    mov r5, r2
    str_ptr r4, r5
    mov r1, fp + 8 // Variable name: fp, internal: fp + 8
    lod_ptr r1, r2
    mov e6, r2
    pop e11
    pop fp
    ret
    jmp after_stmt_151
else_stmt_150:
    mov r1, fp + 8 // Variable name: fp, internal: fp + 8
    mov r2, r1
    mov r4, r2
    mov r1, fp + 8 // Variable name: fp, internal: fp + 8
    lod_ptr r1, r2
    mov e9, r2
    mov e7, 16
    add e10, e9, e7
    mov r2, e10
    mov r5, r2
    str_ptr r4, r5
    mov r1, fp + 8 // Variable name: fp, internal: fp + 8
    lod_ptr r1, r2
    lod32 r2, r2
    mov r4, r2
    mov r7, fp + 12
    str32 r7, r4
    mov r1, fp + 8 // Variable name: fp, internal: fp + 8
    mov r2, r1
    mov r4, r2
    mov r1, fp + 8 // Variable name: fp, internal: fp + 8
    lod_ptr r1, r2
    mov e9, r2
    mov e7, 4
    add e10, e9, e7
    mov r2, e10
    mov r5, r2
    str_ptr r4, r5
    mov r1, fp + 8 // Variable name: fp, internal: fp + 8
    mov r2, r1
    mov r4, r2
    mov r1, fp + 8 // Variable name: fp, internal: fp + 8
    lod_ptr r1, r2
    mov e9, r2
    mov r1, fp + 12 // Variable name: size, internal: fp + 12
    lod32 r1, r2
    mov e7, r2
    add e10, e9, e7
    mov r2, e10
    mov r5, r2
    str_ptr r4, r5
after_stmt_151:
    jmp while_stmt_140_check
while_stmt_140_after:
    mov r1, 0
    mov e6, r1
    pop e11
    pop fp
    ret
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
    mov r1, fp + 0
    str_ptr r1, e0
    mov r1, fp + 4
    str r1, e1
    mov r1, fp + 0 // Variable name: filename, internal: fp + 0
    lod_ptr r1, r2
    mov r7, r2
    push r7
    call find_file
    mov r4, e6
    mov r7, fp + 5
    str_ptr r7, r4
    mov r1, fp + 5 // Variable name: faddr, internal: fp + 5
    lod_ptr r1, r2
    mov r11, r2
    mov r1, 0
    mov r5, r1
    cmp r11, r11, r5
    jnz r11, if_stmt_160
    jz r11, after_stmt_162
if_stmt_160:
    mov r1, fp + 4 // Variable name: complain_on_not_found, internal: fp + 4
    lod r1, r2
    mov r11, r2
    mov r1, 1
    mov r5, r1
    cmp r11, r11, r5
    jnz r11, if_stmt_165
    jz r11, after_stmt_167
if_stmt_165:
    mov r1, var_169
    lod_ptr r1, r2
    mov r7, r2
    push r7
    mov r1, 233
    mov r7, r1
    push r7
    mov r1, 0
    mov r7, r1
    push r7
    call puts32
    mov r4, e6
    mov r1, fp + 0 // Variable name: filename, internal: fp + 0
    lod_ptr r1, r2
    mov r7, r2
    push r7
    call ffnt
    mov r7, e6
    push r7
    mov r1, 233
    mov r7, r1
    push r7
    mov r1, 0
    mov r7, r1
    push r7
    call puts32
    mov r4, e6
    mov r1, var_171
    lod_ptr r1, r2
    mov r7, r2
    push r7
    mov r1, 233
    mov r7, r1
    push r7
    mov r1, 0
    mov r7, r1
    push r7
    call puts32
    mov r4, e6
    mov r1, fp + 9 // Variable name: f, internal: fp + 9
    mov r2, r1
    mov e6, r2
    pop e11
    pop fp
    ret
    jmp after_stmt_167
after_stmt_167:
    jmp after_stmt_162
after_stmt_162:
    mov r1, fp + 9 // Variable name: f, internal: fp + 9
    mov e8, 0
    add r1, r1, e8
    mov r2, r1
    mov r4, r2
    mov r1, fp + 5 // Variable name: faddr, internal: fp + 5
    lod_ptr r1, r2
    mov r5, r2
    str_ptr r4, r5
    mov r1, fp + 9 // Variable name: f, internal: fp + 9
    mov e8, 4
    add r1, r1, e8
    mov r2, r1
    mov r4, r2
    mov r1, fp + 0 // Variable name: filename, internal: fp + 0
    lod_ptr r1, r2
    mov r7, r2
    push r7
    call ffnt
    mov r5, e6
    str_ptr r4, r5
    mov r1, fp + 9 // Variable name: f, internal: fp + 9
    mov r2, r1
    mov e6, r2
    pop e11
    pop fp
    ret
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
    mov r1, fp + 0
    str_ptr r1, e0
    mov r1, fp + 0 // Variable name: filename, internal: fp + 0
    lod_ptr r1, r2
    mov r7, r2
    push r7
    mov r1, 1
    mov r7, r1
    push r7
    call fopen
    mov r4, e6
    mov r7, fp + 4
    str_ptr r7, r4
    mov r1, fp + 4 // Variable name: f, internal: fp + 4
    lod_ptr r1, r1
    mov e8, 0
    add r1, r1, e8
    lod_ptr r1, r2
    mov r4, r2
    mov r7, fp + 8
    str_ptr r7, r4
    mov r1, fp + 8 // Variable name: fptr, internal: fp + 8
    mov r2, r1
    mov r4, r2
    mov r1, fp + 8 // Variable name: fptr, internal: fp + 8
    lod_ptr r1, r2
    mov e9, r2
    mov e7, 4
    sub e10, e9, e7
    mov r2, e10
    mov r5, r2
    str_ptr r4, r5
    mov r1, fp + 8 // Variable name: fptr, internal: fp + 8
    lod_ptr r1, r2
    lod32 r2, r2
    mov e6, r2
    pop e11
    pop fp
    ret
    pop e11
    pop fp
    ret
flist:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_flist
    sub fp, fp, r12
    push e11
    mov r1, 1560
    mov r4, r1
    mov r2, r1
    mov r4, r2
    mov r7, fp + 0
    str_ptr r7, r4
    mov r1, fp + 0 // Variable name: fsp, internal: fp + 0
    lod_ptr r1, r2
    lod_ptr r2, r2
    mov r4, r2
    mov r7, fp + 4
    str_ptr r7, r4
while_stmt_181_check:
    mov r1, 1
    mov r11, r1
    jnz r11, while_stmt_181_body
    jmp while_stmt_181_after
while_stmt_181_body:
    mov r1, fp + 4 // Variable name: fp, internal: fp + 4
    lod_ptr r1, r2
    lod32 r2, r2
    mov r11, r2
    mov r1, 1279677254
    mov r5, r1
    cmp r11, r11, r5
    jz r11, if_stmt_185
    jnz r11, after_stmt_187
if_stmt_185:
    jmp while_stmt_181_after
    jmp after_stmt_187
after_stmt_187:
    mov r1, fp + 4 // Variable name: fp, internal: fp + 4
    mov r2, r1
    mov r4, r2
    mov r1, fp + 4 // Variable name: fp, internal: fp + 4
    lod_ptr r1, r2
    mov e9, r2
    mov e7, 4
    add e10, e9, e7
    mov r2, e10
    mov r5, r2
    str_ptr r4, r5
    mov r1, fp + 4 // Variable name: fp, internal: fp + 4
    lod_ptr r1, r2
    mov r7, r2
    push r7
    call ffnt
    mov r7, e6
    push r7
    mov r1, 255
    mov r7, r1
    push r7
    mov r1, 0
    mov r7, r1
    push r7
    call puts32
    mov r4, e6
    mov r1, var_189
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
    mov r1, fp + 4 // Variable name: fp, internal: fp + 4
    mov r2, r1
    mov r4, r2
    mov r1, fp + 4 // Variable name: fp, internal: fp + 4
    lod_ptr r1, r2
    mov e9, r2
    mov e7, 16
    add e10, e9, e7
    mov r2, e10
    mov r5, r2
    str_ptr r4, r5
    mov r1, fp + 4 // Variable name: fp, internal: fp + 4
    lod_ptr r1, r2
    lod32 r2, r2
    mov r4, r2
    mov r7, fp + 8
    str32 r7, r4
    mov r1, fp + 4 // Variable name: fp, internal: fp + 4
    mov r2, r1
    mov r4, r2
    mov r1, fp + 4 // Variable name: fp, internal: fp + 4
    lod_ptr r1, r2
    mov e9, r2
    mov e7, 4
    add e10, e9, e7
    mov r2, e10
    mov r5, r2
    str_ptr r4, r5
    mov r1, fp + 4 // Variable name: fp, internal: fp + 4
    mov r2, r1
    mov r4, r2
    mov r1, fp + 4 // Variable name: fp, internal: fp + 4
    lod_ptr r1, r2
    mov e9, r2
    mov r1, fp + 8 // Variable name: size, internal: fp + 8
    lod32 r1, r2
    mov e7, r2
    add e10, e9, e7
    mov r2, e10
    mov r5, r2
    str_ptr r4, r5
    jmp while_stmt_181_check
while_stmt_181_after:
    pop e11
    pop fp
    ret
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
    mov r1, fp + 0
    str_ptr r1, e0
    mov r1, fp + 4
    str_ptr r1, e1
    mov r1, fp + 0 // Variable name: name, internal: fp + 0
    lod_ptr r1, r2
    mov r7, r2
    push r7
    mov r1, 1
    mov r7, r1
    push r7
    call fopen
    mov r4, e6
    mov r7, fp + 8
    str_ptr r7, r4
    mov r1, fp + 8 // Variable name: f, internal: fp + 8
    lod_ptr r1, r1
    mov e8, 0
    add r1, r1, e8
    lod_ptr r1, r2
    mov r4, r2
    mov r7, fp + 12
    str_ptr r7, r4
    mov r1, fp + 12 // Variable name: cptr, internal: fp + 12
    lod_ptr r1, r2
    mov r11, r2
    mov r1, 0
    mov r5, r1
    cmp r11, r11, r5
    jnz r11, if_stmt_199
    jz r11, after_stmt_201
if_stmt_199:
    pop e11
    pop fp
    ret
    jmp after_stmt_201
after_stmt_201:
    mov r1, fp + 4 // Variable name: content, internal: fp + 4
    lod_ptr r1, r2
    mov r7, r2
    push r7
    mov r1, fp + 12 // Variable name: cptr, internal: fp + 12
    lod_ptr r1, r2
    mov r7, r2
    push r7
    call strcpy
    mov r4, e6
    mov r1, fp + 12 // Variable name: cptr, internal: fp + 12
    lod_ptr r1, r2
    mov e9, r2
    mov e7, 512
    div e10, e9, e7
    mov r2, e10
    mov r4, r2
    mov r7, fp + 16
    str32 r7, r4
    mov r1, fp + 16 // Variable name: sec, internal: fp + 16
    lod32 r1, r2
    mov r7, r2
    push r7
    call save_sector
    mov r4, e6
    mov r1, fp + 16 // Variable name: sec, internal: fp + 16
    lod32 r1, r2
    mov e9, r2
    mov e7, 1
    sub e10, e9, e7
    mov r2, e10
    mov r7, r2
    push r7
    call save_sector
    mov r4, e6
    mov r1, fp + 16 // Variable name: sec, internal: fp + 16
    lod32 r1, r2
    mov e9, r2
    mov e7, 1
    add e10, e9, e7
    mov r2, e10
    mov r7, r2
    push r7
    call save_sector
    mov r4, e6
    pop e11
    pop fp
    ret
fstrap:
    pop e11
    push fp
    mov r12, _builtin_lcc_basin_fstrap
    sub fp, fp, r12
    push e11
    mov r1, 1560
    mov r4, r1
    mov r2, r1
    mov r4, r2
    mov r7, fp + 0
    str_ptr r7, r4
    mov r1, fp + 0 // Variable name: fptr, internal: fp + 0
    lod_ptr r1, r2
    lod_ptr r2, r2
    mov r4, r2
    mov r7, fp + 4
    str_ptr r7, r4
while_stmt_207_check:
    mov r1, 1
    mov r11, r1
    jnz r11, while_stmt_207_body
    jmp while_stmt_207_after
while_stmt_207_body:
    mov r1, fp + 4 // Variable name: fstart_addr, internal: fp + 4
    lod_ptr r1, r2
    lod32 r2, r2
    mov r11, r2
    mov r1, 1279677254
    mov r5, r1
    cmp r11, r11, r5
    jz r11, if_stmt_211
    jnz r11, after_stmt_213
if_stmt_211:
    jmp while_stmt_207_after
    jmp after_stmt_213
after_stmt_213:
    mov r1, fp + 4 // Variable name: fstart_addr, internal: fp + 4
    lod_ptr r1, r2
    mov e9, r2
    mov e7, 512
    div e10, e9, e7
    mov r2, e10
    mov r4, r2
    mov r7, fp + 8
    str32 r7, r4
    mov r1, fp + 4 // Variable name: fstart_addr, internal: fp + 4
    mov r2, r1
    mov r4, r2
    mov r1, fp + 4 // Variable name: fstart_addr, internal: fp + 4
    lod_ptr r1, r2
    mov e9, r2
    mov e7, 20
    add e10, e9, e7
    mov r2, e10
    mov r5, r2
    str_ptr r4, r5
    mov r1, fp + 4 // Variable name: fstart_addr, internal: fp + 4
    lod_ptr r1, r2
    lod32 r2, r2
    mov r4, r2
    mov r7, fp + 12
    str32 r7, r4
    mov r1, fp + 12 // Variable name: size, internal: fp + 12
    lod32 r1, r2
    mov e9, r2
    mov e7, 512
    div e10, e9, e7
    mov r2, e10
    mov r4, r2
    mov r7, fp + 16
    str32 r7, r4
    mov r1, fp + 16 // Variable name: sectors, internal: fp + 16
    lod32 r1, r2
    mov r7, r2
    push r7
    mov r1, fp + 8 // Variable name: osector, internal: fp + 8
    lod32 r1, r2
    mov r7, r2
    push r7
    call offset_sec_load
    mov r4, e6
    mov r1, fp + 0 // Variable name: fptr, internal: fp + 0
    mov r2, r1
    mov r4, r2
    mov r1, fp + 0 // Variable name: fptr, internal: fp + 0
    lod_ptr r1, r2
    mov e9, r2
    mov r1, fp + 12 // Variable name: size, internal: fp + 12
    lod32 r1, r2
    mov e7, r2
    add e10, e9, e7
    mov r2, e10
    mov r5, r2
    str_ptr r4, r5
    jmp while_stmt_207_check
while_stmt_207_after:
    pop e11
    pop fp
    ret
