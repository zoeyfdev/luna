#include <string.h>
#include <stdio.h>

#include "../../../lcc_shared/liberror.h"

typedef struct {
    char name[4];
    bool taken;
} _register;

#define NUM_REGISTERS 18
_register registers[NUM_REGISTERS];

void regs_init() {
    int base = 1;
    for (int i = 0; i < 12; i++, base++) { // r registers
        sprintf(registers[i].name, "r%d", base);
    }

    base = 7;
    for (int i = 12; i < 18; i++, base++) { // e registers
        sprintf(registers[i].name, "e%d", base);
    }
}

char* take_register() {
    for (int i = 0; i < 35; i++) {
        if (!registers[i].taken) {
            registers[i].taken = true;
            return registers[i].name;
        }
    }

    lcc_error(NULL, 0, "no more free registers!", NULL);
    return "";
}

void free_register(char* name) {
    for (int i = 0; i < 35; i++) {
        if (!strcmp(registers[i].name, name))
            registers[i].taken = false;
    }
}
