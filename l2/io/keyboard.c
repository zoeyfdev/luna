#include <stdlib.h>
#include <stdint.h>

#define KEYBOARD_RAM_SIZE 1
unsigned char* KEYBOARD_MEMORY = NULL;

unsigned char char_table[][2] = {
    {'a', 'A'},
    {'b', 'B'},
    {'c', 'C'},
    {'d', 'D'},
    {'e', 'E'},
    {'f', 'F'},
    {'g', 'G'},
    {'h', 'H'},
    {'i', 'I'},
    {'j', 'J'},
    {'k', 'K'},
    {'l', 'L'},
    {'m', 'M'},
    {'n', 'N'},
    {'o', 'O'},
    {'p', 'P'},
    {'q', 'Q'},
    {'r', 'R'},
    {'s', 'S'},
    {'t', 'T'},
    {'u', 'U'},
    {'v', 'V'},
    {'w', 'W'},
    {'x', 'X'},
    {'y', 'Y'},
    {'z', 'Z'},
    {'`', '~'},
    {'1', '!'},
    {'2', '@'},
    {'3', '#'},
    {'4', '$'},
    {'5', '%'},
    {'6', '^'},
    {'7', '&'},
    {'8', '*'},
    {'9', '('},
    {'0', ')'},
    {'-', '_'},
    {'=', '+'},
    {'[', '{'},
    {']', '}'},
    {'\\', '|'},
    {';', ':'},
    {'\'', '"'},
    {',', '<'},
    {'.', '>'},
    {'/', '?'}
};

int keyboard_upper(int code) {
    for (int i = 0; i < 21 + 26; i++) {
        unsigned char* pair = char_table[i];
        if (pair[0] == code)
            return pair[1];
    }
    return code;
}

int keyboard_lower(int code) {
    for (int i = 0; i < 21 + 26; i++) {
        unsigned char* pair = char_table[i];
        if (pair[1] == code)
            return pair[0];
    }
    return code;
}

unsigned char keyboard_read_memory(uint32_t address) {
    if (address < KEYBOARD_RAM_SIZE)
        return KEYBOARD_MEMORY[address];
    return (unsigned char) rand() & 0xFF;
}

void keyboard_write_memory(uint32_t address, unsigned char value) {
    if (address < KEYBOARD_RAM_SIZE)
        KEYBOARD_MEMORY[address] = value;
}
