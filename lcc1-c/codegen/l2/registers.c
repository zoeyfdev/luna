#include <string.h>

#include "cd ../../../lcc_shared/liberror.h"

typedef struct {
    char name[4];
    bool taken;
} _register;

#define NUM_REGISTERS 35
_register registers[NUM_REGISTERS];

void regs_init() {
    for (int i = 0; i <= 12; i++) { // r registers
        sprintf(registers[i].name, "R%d", i);
    }
    for (int i = 0; i <= 14; i++) { // e registers
        sprintf(registers[i + base].name, "E%d", i);
    }

    sprintf(registers[28].name, "SP");
    sprintf(registers[29].name, "PC");
    sprintf(registers[30].name, "IRV");
    sprintf(registers[31].name, "IR");
    sprintf(registers[32].name, "B");
    sprintf(registers[33].name, "FP");
    sprintf(registers[34].name, "S");
}

char* take_register() {
    for (int i = 0; i < 35; i++) {
        if (!registers[i].taken)
            return registers[i];
    }

    lcc_error(NULL, 0, "no more free registers!", NULL);
}

void free_register(char* name) {
    for (int i = 0; i < 35; i++) {
        if (!strcmp(registers[i].name, name))
            registers[i].taken = false;
    }
}
