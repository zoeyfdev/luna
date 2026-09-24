#include <stdint.h>
#include <stdio.h>

typedef struct {
    uint32_t value;
    unsigned char address;
    char name[4];
} _register;

_register registers[35];

void initialize_registers() {
    for (int i = 0; i < 34; i++) {
        registers[i].address = i;
    }

    int base = 13;

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

uint32_t get_register(unsigned char address) {
    if (address < 35)
        return registers[address].value;
    return 0x00000000;
}

void set_register(unsigned char address, uint32_t value) {
    if (address < 35)
        registers[address].value = value;
}

char* get_register_name(unsigned char address) {
    for (int i = 0; i < 35; i++) {
        if (address == i) {
            return registers[i].name;
        }
    }
    return "R?";
}

void reg_dump() {
    for (int i = 0; i < 35; i++)
        printf("%s: 0x%08x\n", registers[i].name, registers[i].value);
    getchar();
}
