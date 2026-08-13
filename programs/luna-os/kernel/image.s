.bits 32
.global fzip_decode
.global fzipdecode_loc

fzipdecode_loc:
fzip_decode:
    // .byte 0x21
    pop e11
    pop r5 // image
    push e11

    mov r10, 4

    lod32 r5, r11 // size in r11

    push r11
    call malloc

    push e6

    mov r12, e6 // buffer in r12

    add r5, r5, r10 // advance past size

    mov r7, 0 // count
    mov e0, 0 // global count

    mov e10, pc

    lod32 r5, r6 // Load current n into r6
    add r5, r5, r10 // advance past n
    lod r5, r9 // char in r9
    inc r5

fzd_inner:
    cmp r8, r7, r6
    jnz r8, fzd_next
    
    str r12, r9 // store byte
    inc r7
    inc e0
    inc r12
    jmp fzd_inner

fzd_next:
    // .byte 0x21
    mov r7, 0
    cmp e1, e0, r11
    jnz e1, fzd_ret
    jmp e10
fzd_ret:
    pop e6
    pop e11
    ret
