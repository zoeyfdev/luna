#include <stdlib.h>

#include "../io/keyboard.h"

void initialize_keyboard() {
    KEYBOARD_MEMORY = calloc(1, 1);
}
