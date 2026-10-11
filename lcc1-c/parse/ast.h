#pragma once

#include "types.h"

typedef enum {
    AST_DECLARATION,
    AST_STATEMENT,
    AST_EXPRESSION,
    AST_LEAF
} ast_node_class;

typedef enum {
    AST_INTEGER_LITERAL,
    AST_STRING_LITERAL,
    AST_UNARY,
    AST_BINARY,
    AST_ASSIGNMENT 
} ast_node_type;

typedef enum {
    BIN_OP_PLUS,
    BIN_OP_MINUS,
    BIN_OP_MULTIPLY,
    BIN_OP_DIVIDE,
} ast_binary_op_type;

typedef struct {
    c_type type;
    char* value;
} integer_literal;

typedef struct {
    c_type type;
    char* value;
} string_literal;

typedef struct {
    c_type type;
    ast_binary_op_type operation;
    struct ast_node* left;
    struct ast_node* right;
} binary_operation;

typedef struct {
    ast_node_type type;
    union {
        integer_literal integer_literal;
        string_literal string_literal;
        binary_operation binary_operation;
    } content;
} ast_node;

typedef struct {
    ast_node** nodes;
} ast;
