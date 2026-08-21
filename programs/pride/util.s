#ifdef PORTABLE
    .bits 32
#else
    .bits 16
#endif

.global pit_nxt
.global sleep
.global putc
.global key_click

sleep:
    pop e11
    pop r4 // Seconds
    push e11

    mov r5, 0 

    mov e10, pc

    call pit_handler
    inc r5

    cmp r6, r4, r5
    jz r6, e10

    pop e11
    ret

pit_handler:
    pop e11

    #ifdef PORTABLE
        mov r1, 0x6FFF0007 // 0xFA3E for 16 bit
    #else
        mov r1, 0xFA3E
    #endif

    mov r2, 1
    str r1, r2
pit_wait:
    hlt
    jmp pit_wait
pit_nxt:
    #ifdef PORTABLE
        mov r1, 0x6FFF0007 // 0xFA3E for 16 bit
    #else
        mov r1, 0xFA3E
    #endif

    mov r2, 0
    str r1, r2

    ret

putc:
    pop e11
    pop r2
    pop r1

    mov r3, r2

    int 1

    ret

#ifdef PORTABLE
key_click:
    pusha
    mov r1, 0x6FFF0019
    mov r2, 0
    str r1, r2 // DISABLE KEYBOARD INTERRUPT

    mov r1, 0x80000012
    lod r1, r2

    mov r3, "q"
    cmp r4, r2, r3
    jnz r4, kc_exit
    jmp kc_ret
kc_exit:
    popa
    mov r1, 1
    syscall
kc_ret:
    popa
    jmp irv
#endif
