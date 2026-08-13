.bits 32
jmp _start

#define TICK_TIME 62

_start:
    mov sp, 0x60000000 
    jmp main_setup

modulo:
    pop e11
    pop r2 // Divisor
    pop r1 // Dividend

    div r3, r1, r2
    mul r4, r3, r2
    sub e6, r3, r4

    ret

gen_random:
    pop e11
    push e11

    mov r5, 0x90909090
    lod r5, r5 // Random number
    
    push r5
    push 2
    call modulo

    jz e6, gen_random_yes
    jnz e6, gen_random_no
gen_random_yes:
    mov e6, 2
    jmp gen_random_ret
gen_random_no:
    mov e6, 0
    jmp gen_random_ret 
gen_random_ret:
    pop e11
    ret

main_setup:
    // PIT to 0.5 seconds
    mov r1, TICK_TIME
    mov r2, 0x80000013
    str32 r2, r1

    mov r1, TICK_TIME // milliseconds
    mov r2, 0x80000017
    str32 r2, r1
    // Set up IDT
    mov r1, pit_nxt
    mov r2, 0x6FFF0008
    str32 r2, r1

    mov r1, key_click
    mov r2, 0x6FFF001A
    str32 r2, r1
main:
// int 0x02
// 07 = mode
// 08 09 0A 0B = addr
    mov r1, 0x6FFF0019
    mov r2, 1
    str r1, r2 // ENABLE KEYBOARD INTERRUPT
main_render_plr:
    mov r10, GROUND
    mov r11, 1000
    add r10, r10, r11 // initial in r10
    mov r12, PLAYER_POS
    lod r12, e0 // pos in e0

    sub r10, r10, e0 // POS ON GRID IN R10

    push r1
    push r2
    push r3
    push r4

    mov r3, 2
    lod r10, r2
    cmp r4, r2, r3
    jnz r4, game_over

    mov r1, 0x01
    str r10, r1
    
    pop r4
    pop r3
    pop r2
    pop r1
main_main:
    mov r1, 0x20
    mov r2, 0x0F
    mov r3, 0x0F
    mov r4, 1000
    mov r5, GROUND

    add r4, r4, r5 // make r4 have end addr

    mov e10, pc

    lod r5, r6
    
    mov r7, 2
    cmp r8, r6, r7

    jnz r8, main_render_meteor
    jnz r6, main_render_player
    jz r6, main_render_sky
main_render_meteor:
    mov r2, 0xA0
    mov r3, 0xA0
    int 1
    jmp main_render_done
main_render_player:
    mov r2, 0b11101100
    mov r3, 0b11101100
    int 1
    jmp main_render_done
main_render_sky:
    mov r2, 0x0F
    mov r3, 0x0F
    int 1
    jmp main_render_done
main_render_done:
    inc r5
    cmp r6, r5, r4
    jz r6, e10
    
add_next:
// SHIFT ALL METEORS DOWN
    mov r10, 1000
    mov r12, GROUND
    mov e1, 40
    mov e3, 0

    add r10, r10, r12

    mov e10, pc

    lod r10, e0
    add e2, r10, e1

    str e2, e0
    str r10, e3 // clear top

    dec r10

    cmp e0, r10, r12
    jz e0, e10


// ADD NEW METEORS TO TOP
    mov r10, 0
    mov r11, 40

    mov e10, pc

    call gen_random
    jnz e6, meteor_yes
    jmp add_next_nxt
meteor_yes:
    mov e0, GROUND
    add e0, e0, r10

    mov r12, 2
    str e0, r12
add_next_nxt:
    inc r10
    cmp r12, r10, r11
    jz r12, e10

    mov r1, 0x6FFF0007
    mov r2, 1
    str r1, r2 // ENABLE PIT
pit_wait:
    hlt
    jmp pit_wait
pit_nxt:
    mov r1, 0x6FFF0007
    mov r2, 0
    str r1, r2 // DISABLE PIT
    jmp main

game_over:
    mov r1, 0x6FFF0007
    mov r2, 0
    str r1, r2 // DISABLE PIT
    mov r1, 0x6FFF0019
    mov r2, 0
    str r1, r2 // DISABLE KEYBOARD INTERRUPT

    mov r1, 16 // X
    mov r2, 8 // Y
    int 0x0c

    push game_over_txt
    call write
    mov r1, 1
    syscall

write:
    pop e11
    pop r4

    mov r2, 0xFF
    mov r3, 0x0F

    mov e10, pc

    lod r4, r1
    int 1
    inc r4
    
    jnz r1, e10
    ret

key_click: 
    pusha

    mov r1, 0x6FFF0007
    mov r2, 0
    str r1, r2 // DISABLE PIT

    mov r1, 0x6FFF0019
    mov r2, 0
    str r1, r2 // DISABLE KEYBOARD INTERRUPT

    mov r1, 0x80000012
    lod r1, r2

    mov r3, "a"
    cmp r4, r2, r3
    jnz r4, key_click_right

    mov r3, "d"
    cmp r4, r2, r3
    jnz r4, key_click_left
    
    jmp key_click_ret

key_click_left:
    mov r1, PLAYER_POS
    lod r1, r2
    dec r2

    mov r3, 1
    ilt r4, r2, r3
    jnz r4, rollover_1
    jmp rollover_1_done
rollover_1:
    mov r2, 40
rollover_1_done:
    str r1, r2
    jmp key_click_ret
key_click_right:
    mov r1, PLAYER_POS
    lod r1, r2
    inc r2

    mov r3, 40
    igt r4, r2, r3
    jnz r4, rollover_2
    jmp rollover_2_done
rollover_2:
    mov r2, 1
rollover_2_done:
    str r1, r2
    jmp key_click_ret
key_click_ret:
    mov r1, 0x6FFF0019
    mov r2, 1
    str r1, r2 // ENABLE KEYBOARD INTERRUPT

    mov r1, 0x6FFF0007
    mov r2, 1
    str r1, r2 // ENABLE PIT

    popa
    jmp irv

halt:
    jmp halt

game_over_txt:
    .asciz "Game over"

PLAYER_POS:
    .byte 20

GROUND:
