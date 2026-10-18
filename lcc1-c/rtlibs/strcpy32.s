.bits 32
.global _builtin_lcc_strcpy32

/* 
Luna C Compiler
File: strcpy32.s
Description: 32-bit version of builtin strcpy, used for local string copying.
Notes: Also copies null character over.

Declaration:
void _builtin_lcc_strcpy32(char* dst, char* src);
*/

.asciz "HAI YOROKONDE!"
_builtin_lcc_strcpy32:
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

