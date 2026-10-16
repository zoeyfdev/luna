#pragma once

typedef enum {
    TOK_NONE,

    TOK_PLUS,
    TOK_MINUS,
    TOK_SLASH,

    TOK_STAR,
    TOK_AMPERSAND,
    TOK_EQUAL,

    TOK_L_PAREN,
    TOK_R_PAREN,
    TOK_L_CURLY,
    TOK_R_CURLY,

    TOK_ASM,

    TOK_LONG,
    TOK_SHORT,
    TOK_SIGNED,
    TOK_UNSIGNED,
    TOK_EXTERN,
    TOK_STATIC,
    TOK_CONST,

    TOK_NUMBER,
    TOK_IDENTIFIER,
    TOK_STRING
} token_type;

typedef struct {
    token_type type;
    int line;
    char* filename; 
    char* value;
    char* displayed_value;
} token;
