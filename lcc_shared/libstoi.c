#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include <ctype.h>

char* hex_approved_chars = "0123456789abcdef";
char* bin_approved_chars = "01";

bool hex_is_approved(char c) {
    char lc = tolower(c);

    for (int i = 0; i < 16; i++)
        if (hex_approved_chars[i] == lc)
            return true;
    return false;
}

int return_hex_position(char c) {
    char lc = tolower(c);

    for (int i = 0; i < 16; i++)
        if (hex_approved_chars[i] == lc)
            return i;
    return 0;
}

bool bin_is_approved(char c) {
    char lc = tolower(c);

    for (int i = 0; i < 2; i++)
        if (bin_approved_chars[i] == lc)
            return true;
    return false;
}

int return_bin_position(char c) {
    char lc = tolower(c);

    for (int i = 0; i < 2; i++)
        if (bin_approved_chars[i] == lc)
            return i;
    return 0;
}

int64_t stoi(char* str, bool* worked) {
    if (strlen(str) == 0)
        return 0;

    int64_t atoi_result = atoi(str);
    bool str_is_zero = true;

    for (size_t i = 0; i < strlen(str); i++) {
        if (str[i] != '0') {
            str_is_zero = false;
            break;
        }
    }

    if ((atoi_result != 0) || (atoi_result == 0 && str_is_zero)) {
        *worked = true;
        return atoi_result;
    } else {
        if (strlen(str) > 2) {
            if (str[0] == '0' && tolower(str[1]) == 'x') {
                str += 2;
                int64_t final = 0;
                while (*str) {
                    if (!hex_is_approved(*str))
                        return 0;
                    final = final << 4 | return_hex_position(*str);
                    str++;
                }
                *worked = true;
                return final;
            } else if (str[0] == '0' && tolower(str[1]) == 'b') {
                str += 2;
                int64_t final = 0;
                while (*str) {
                    if (!bin_is_approved(*str))
                        return 0;
                    final = final << 1 | return_bin_position(*str);
                    str++;
                }
                *worked = true;
                return final; 
            }
        }
    }

    return 0;
}
