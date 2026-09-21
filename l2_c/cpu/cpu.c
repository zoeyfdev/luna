#include <stdint.h>
#include <time.h>

#include "registers.h"
#include "memory.h"

uint32_t get_word(uint32_t address) {
    return (uint32_t) ((uint16_t) get_memory(address) << 8 
            | (uint16_t) (get_memory(address + 1) & 0xFF));
}

uint32_t get_dword(uint32_t address) {
    return (uint32_t) get_memory(address) << 24 
        | (uint32_t) get_memory(address + 1) << 16 
        | (uint32_t) get_memory(address + 2) << 8 
        | (uint32_t) (get_memory(address + 3) & 0xFF);
}

void _sleep(int ms) {
    struct timespec ts = {
        .tv_sec = ms / 1000,
        .tv_nsec = (ms % 1000) * 1000000L
    };
    nanosleep(&ts, NULL);
}

void cpu_execute() {
    for (;;) {
        uint32_t pc = get_register(PC);
        unsigned char op = get_memory(pc);

        switch (op) {
        case 0x01: {
                // MOV
                // mov <register> (<value>/<reg>/<reg + disp>)
                unsigned char mode = get_memory(pc + 1);
                unsigned char dest = get_memory(pc + 2);
                
                switch (mode) {
                case 0x01:
                    // Immediate
                    if (!IS_XEN) {
                        set_register(dest, get_word(pc + 3));
                        set_register(PC, pc + 5);
                    } else {
                        set_register(dest, get_dword(pc + 3));
                        set_register(PC, pc + 7);
                    }
                    break;
                case 0x02:
                    // Register
                    set_register(dest, get_register(get_memory(pc + 3)));
                    set_register(PC, pc + 4);
                    break;
                case 0x03:
                    // TODO: displacement
                    break;
                }

                set_register(PC, pc + 1);
                break;
            }
        case 0x02: {
                // HLT
                // hlt
                while (get_register(PC) == pc && get_register(IR) == 0) 
                    _sleep(15);
                set_register(PC, pc + 1);
                break;
            }
        case 0x03: {
                // JMP
                // jmp (<value>/<reg>)
                unsigned char mode = get_memory(pc + 1);

                switch (mode) {
                case 0x01:
                    // Immediate
                    if (!IS_XEN)
                        set_register(PC, get_word(pc + 2));
                    else
                        set_register(PC, get_dword(pc + 2));
                    break;
                case 0x02:
                    set_register(PC, get_register(pc + 2));
                    break;
                }
                break;
            }
        case 0x04: {
                // INT
                // int <value>
                if (!IS_XEN) {
                    set_register(IR, get_register(IR) | (1 << get_word(pc + 1)));
                    set_register(PC, pc + 3);
                } else {
                    set_register(IR, get_register(IR) | (1 << get_dword(pc + 1)));
                    set_register(PC, pc + 5);
                }

                break;
            }
        case 0x05: {
                // JNZ
                // jnz <register> (<value>/<register>)
                unsigned char mode = get_memory(pc + 1);
                unsigned char reg = get_memory(pc + 2);

                if (get_register(reg) != 0) {
                    switch (mode) {
                    case 0x01:
                        // Imm
                        if (!IS_XEN)
                            set_register(PC, get_register(reg) != 0 ? get_word(pc + 3) : pc + 5);
                        else
                            set_register(PC, get_register(reg) != 0 ? get_word(pc + 3) : pc + 7);
                        break;
                    case 0x02:
                        // Reg
                        set_register(PC, get_register(reg) != 0 ? get_register(get_memory(pc + 3)) : pc + 4);
                        break;
                    }
                }
                break;
            }
        case 0x06: {
                // NOP
                // nop
                set_register(PC, pc + 1);
                break;
            }
        case 0x07: {
                // CMP
                // cmp <to> <r1> <r1>
                set_register(get_memory(pc + 1), get_register(get_memory(pc + 2)) == get_register(get_memory(pc + 3)) ? 1 : 0);
                set_register(PC, pc + 4);
                break;
            }
        case 0x08: {
                // JZ
                // jz <register> (<value>/<register>)
                unsigned char mode = get_memory(pc + 1);
                unsigned char reg = get_memory(pc + 2);

                if (get_register(reg) != 0) {
                    switch (mode) {
                    case 0x01:
                        // Imm
                        if (!IS_XEN)
                            set_register(PC, get_register(reg) == 0 ? get_word(pc + 3) : pc + 5);
                        else
                            set_register(PC, get_register(reg) == 0 ? get_word(pc + 3) : pc + 7);
                        break;
                    case 0x02:
                        // Reg
                        set_register(PC, get_register(reg) == 0 ? get_register(get_memory(pc + 3)) : pc + 4);
                        break;
                    }
                }
                break;
            }
        case 0x09: {
                // INC
                // inc <register>
                set_register(get_register(get_memory(pc + 1)), get_register(get_memory(pc + 1)) + 1);
                set_register(PC, pc + 2);
                break;
            }
        case 0x0a: {
                // DEC
                // dec <register>
                set_register(get_register(get_memory(pc + 1)), get_register(get_memory(pc + 1)) - 1);
                set_register(PC, pc + 2);
                break;
            }
        case 0x0b: {
                // PUSH
                // push <value/reg>
                unsigned char mode = get_memory(pc + 1);
                uint32_t val = 0;
                switch (mode) {
                case 0x01:
                    // Imm
                    if (!IS_XEN) {
                        val = get_word(pc + 2);
                        set_register(PC, pc + 4); 
                    } else {
                        val = get_dword(pc + 2);
                        set_register(PC, pc + 6);
                    }
                    break;
                case 0x02:
                    val = get_register(get_memory(pc + 3));
                    set_register(PC, pc + 4);
                    break;
                }
                if (!IS_XEN) { 
                    set_memory(get_register(SP), (unsigned char) val >> 8);
                    set_memory(get_register(SP) + 1, (unsigned char) val & 0xFF);

                    set_register(SP, get_register(SP) - 2);
                } else {
                    set_memory(get_register(SP), (unsigned char) val >> 24);
                    set_memory(get_register(SP) + 1, (unsigned char) val >> 16);
                    set_memory(get_register(SP) + 2, (unsigned char) val >> 8);
                    set_memory(get_register(SP) + 3, (unsigned char) val & 0xFF);

                    set_register(SP, get_register(SP) - 4); 
                }
                break;
            }
        case 0x0c: {
                // POP
                // pop <reg>
                
                unsigned char reg = get_memory(pc + 1);
                if (!IS_XEN) {
                    set_register(SP, get_register(SP) + 2);
                    set_register(reg, get_memory(get_register(SP)) << 8 
                            | get_memory(get_register(SP) + 1) & 0xFF);
                } else {
                    set_register(SP, get_register(SP) + 4);
                    set_register(reg, get_memory(get_register(SP)) << 24 
                            | get_memory(get_register(SP) + 1) << 16
                            | get_memory(get_register(SP) + 2) << 8
                            | get_memory(get_register(SP) + 3) & 0xFF);
                }
                
                set_register(PC, pc + 2);
                break;
            }
        case 0x0d: {
                // ADD
                // add <dest> <reg1> <reg2>
                unsigned char dest = get_memory(pc + 1);
                uint32_t reg1 = get_register(get_memory(pc + 2));
                uint32_t reg2 = get_register(get_memory(pc + 3));
                set_register(dest, reg1 + reg2);
                set_register(PC, pc + 4);
                break;
            }
        case 0x0e: {
                // SUB
                // SUB <dest> <reg1> <reg2>
                unsigned char dest = get_memory(pc + 1);
                uint32_t reg1 = get_register(get_memory(pc + 2));
                uint32_t reg2 = get_register(get_memory(pc + 3));
                set_register(dest, reg1 - reg2);
                set_register(PC, pc + 4);
                break;
            }
        case 0x0f: {
                // MUL
                // mul <dest> <reg1> <reg2>
                unsigned char dest = get_memory(pc + 1);
                uint32_t reg1 = get_register(get_memory(pc + 2));
                uint32_t reg2 = get_register(get_memory(pc + 3));
                set_register(dest, reg1 * reg2);
                set_register(PC, pc + 4);
                break;
            }
        case 0x10: {
                // DIV
                // div <dest> <reg1> <reg2>
                unsigned char dest = get_memory(pc + 1);
                uint32_t reg1 = get_register(get_memory(pc + 2));
                uint32_t reg2 = get_register(get_memory(pc + 3));

                if (reg2 != 0)
                    set_register(dest, reg1 / reg2);
                set_register(PC, pc + 4);
                break;
            }
        case 0x11: {
                // IGT
                // IGT <dest> <reg1> <reg2>
                unsigned char dest = get_memory(pc + 1);
                uint32_t reg1 = get_register(get_memory(pc + 2));
                uint32_t reg2 = get_register(get_memory(pc + 3));
                set_register(dest, reg1 > reg2 ? 1 : 0);
                set_register(PC, pc + 4);
                break;
            }
        case 0x12: {
                // ILT
                // ILT <dest> <reg1> <reg2>
                unsigned char dest = get_memory(pc + 1);
                uint32_t reg1 = get_register(get_memory(pc + 2));
                uint32_t reg2 = get_register(get_memory(pc + 3));
                set_register(dest, reg1 < reg2 ? 1 : 0);
                set_register(PC, pc + 4);
                break;
            }
        case 0x13: {
                // AND
                // and <dest> <reg1> <reg2>
                unsigned char dest = get_memory(pc + 1);
                uint32_t reg1 = get_register(get_memory(pc + 2));
                uint32_t reg2 = get_register(get_memory(pc + 3));
                set_register(dest, reg1 & reg2);
                set_register(PC, pc + 4);
                break;
            }
        case 0x14: {
                // OR
                // or <dest> <reg1> <reg2>
                unsigned char dest = get_memory(pc + 1);
                uint32_t reg1 = get_register(get_memory(pc + 2));
                uint32_t reg2 = get_register(get_memory(pc + 3));
                set_register(dest, reg1 | reg2);
                set_register(PC, pc + 4);
                break;
            }
        case 0x15: {
                // NOT
                // not <dest> <reg1>
                unsigned char dest = get_memory(pc + 1);
                uint32_t reg1 = get_register(get_memory(pc + 2));
                set_register(dest, reg1 ^ reg1);
                set_register(PC, pc + 3);
                break;
            }
        case 0x16: {
                // XOR
                // xor <dest> <reg1> <reg2>
                unsigned char dest = get_memory(pc + 1);
                uint32_t reg1 = get_register(get_memory(pc + 2));
                uint32_t reg2 = get_register(get_memory(pc + 3));
                set_register(dest, reg1 ^ reg2);
                set_register(PC, pc + 4);
                break;
            }
        case 0x20: {
                // MOD
                // MOD <dest> <reg1> <reg2>
                unsigned char dest = get_memory(pc + 1);
                uint32_t reg1 = get_register(get_memory(pc + 2));
                uint32_t reg2 = get_register(get_memory(pc + 3));

                if (reg2 != 0)
                    set_register(dest, reg1 % reg2);
                set_register(PC, pc + 4);
                break;
            }
        case 0x1c: {
                // SHL
                // shl <dest> <reg1> <reg2>
                unsigned char dest = get_memory(pc + 1);
                uint32_t reg1 = get_register(get_memory(pc + 2));
                uint32_t reg2 = get_register(get_memory(pc + 3));
                set_register(dest, reg1 << reg2);
                set_register(PC, pc + 4);
                break;
            }
        case 0x1d: {
                // SHR
                // shr <dest> <reg1> <reg2>
                unsigned char dest = get_memory(pc + 1);
                uint32_t reg1 = get_register(get_memory(pc + 2));
                uint32_t reg2 = get_register(get_memory(pc + 3));
                set_register(dest, reg1 >> reg2);
                set_register(PC, pc + 4);
                break;
            }
        case 0x17: {
                // LOD
                // lod <addr register> <register>
                unsigned char reg = get_memory(pc + 1);
                unsigned char val = get_memory(get_register(get_memory(pc + 2)));
                set_register(reg, val);
                set_register(PC, pc + 3);
                break;
            }
        case 0x19: {
                // LOD16
                // lod16 <addr register> <register>
                unsigned char reg = get_memory(pc + 1);
                unsigned char val = get_word(get_register(get_memory(pc + 2)));
                set_register(reg, val);
                set_register(PC, pc + 3);
                break;
            }
        case 0x1e: {
                // LOD32
                // lod32 <addr register> <register>
                unsigned char reg = get_memory(pc + 1);
                unsigned char val = get_dword(get_register(get_memory(pc + 2)));
                set_register(reg, val);
                set_register(PC, pc + 3);
                break;
            }
        case 0x1b: {
                // STR
                // str <addr register> <register>
                uint32_t addr = get_register(get_memory(pc + 1));
                uint32_t val = get_register(get_memory(pc + 2));
                set_memory(addr, val & 0xFF);
                set_register(PC, pc + 3);
                break;
            }
        case 0x18: {
                // STR16
                // str16 <addr register> <register>
                uint32_t addr = get_register(get_memory(pc + 1));
                uint32_t val = get_register(get_memory(pc + 2));
                set_memory(addr, val << 8);
                set_memory(addr + 1, val & 0xFF);
                set_register(PC, pc + 3);
                break;
            }
        case 0x1f: {
                // STR32
                // str32 <addr register> <register>
                uint32_t addr = get_register(get_memory(pc + 1));
                uint32_t val = get_register(get_memory(pc + 2));
                set_memory(addr, val << 24);
                set_memory(addr + 1, val << 16);
                set_memory(addr + 2, val << 8);
                set_memory(addr + 3, val & 0xFF);
                set_register(PC, pc + 3);
                break;
            }
        case 0x1a: {
                unsigned char mode = get_memory(pc + 1);
                switch (mode) {
                case 0x00:
                    // 16-bit mode
                    set_register(S, get_register(S) << 31 | CPU_FLAG_XEN & 0);
                    break;
                case 0x01:
                    // 32-bit mode
                    set_register(S, get_register(S) << 31 | CPU_FLAG_XEN & 1);
                    break;
                }
                set_register(PC, pc + 1);
                break;
            }
        default:
            // Illegal
            break;
        }
    }
}

void* cpu_init() {
    initialize_registers();
    initialize_memory();
    cpu_execute();
    
    return NULL;
}
