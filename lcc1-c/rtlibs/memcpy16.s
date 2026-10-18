.bits 16
.global _builtin_lcc_memcpy16

/* 
Luna C Compiler
File: memcpy16.s
Description: 16-bit version of builtin memcpy, used for copying structs.
Notes: none.

Declaration:
void _builtin_lcc_memcpy16(void* dst, void* src, int n);
*/

_builtin_lcc_memcpy16:
    pop e11
    pop r5 // Number of bytes to copy
    pop r6 // Source address
    pop r7 // Destination address

    xor r4, r4, r4
    
    mov e10, pc

    lod r6, r8
    str r7, r8

    inc r6
    inc r7
    inc r4

    cmp r9, r5, r4
    jz r9, e10

    ret
