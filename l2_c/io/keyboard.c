#include <stdlib.h>

unsigned char* KEYBOARD_MEMORY = NULL;

unsigned char char_table[][2] = {
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
    for (int i = 0; i < 21; i++) {
        unsigned char* pair = char_table[i];
        if (pair[0] == code)
            return pair[1];
    }
    return code;
}

int keyboard_lower(int code) {
    for (int i = 0; i < 21; i++) {
        unsigned char* pair = char_table[i];
        if (pair[1] == code)
            return pair[0];
    }
    return code;
}
