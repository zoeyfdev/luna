.bits 16
.global _builtin_lcc_strcpy16

/* 
Luna C Compiler
File: strcpy16.s
Description: 16-bit version of builtin strcpy, used for local string copying.
Notes: Also copies null character over.

Declaration:
void _builtin_lcc_strcpy16(char* dst, char* src);
*/

_builtin_lcc_strcpy16:
    pop e11 // Return address
    pop r2 // Destination
    pop r1 // Source

    mov e10, pc

    lod r2, r3
    str r1, r3

    inc r1
    inc r2
    
    jnz r3, e10
    ret


