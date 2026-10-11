#pragma once

typedef enum {
    TOK_NUMBER,
    TOK_PLUS,
    TOK_MINUS,
    TOK_STAR,
    TOK_SLASH,
    TOK_L_PAREN,
    TOK_R_PAREN,
} token_type;

typedef struct {
    token_type type;
    int line;
    char* filename; 
    char* value;
    char* displayed_value;
} token;
