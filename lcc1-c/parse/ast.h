#pragma once

#include "types.h"
#include "../../lcc_shared/libvector.h"

/*
typedef enum {
    AST_DECLARATION,
    AST_STATEMENT,
    AST_EXPRESSION,
    AST_LEAF
} ast_node_class;
*/

typedef enum {
    AST_INTEGER_LITERAL,
    AST_STRING_LITERAL,

    AST_UNARY,

    AST_BINARY,
    AST_ASSIGNMENT,
    AST_FUNCTION_CALL,

    AST_ASSEMBLY_STATEMENT,
    AST_STATEMENT_EXPRESSION,

    AST_VARIABLE,
    AST_ASSEMBLY_DECL
} ast_node_type;

typedef enum {
    BIN_OP_PLUS,
    BIN_OP_MINUS,
    BIN_OP_MULTIPLY,
    BIN_OP_DIVIDE,
    BIN_OP_EQUAL
} ast_binary_op_type;

typedef enum {
    KIND_VARIABLE,
    KIND_FUNCTION
} ast_variable_kind;

typedef struct ast_node ast_node;

// Leaves

typedef struct {
    c_type type_info;
    char* value;
} integer_literal;

typedef struct {
    c_type type_info;
    char* value;
} string_literal;

// Expressions

typedef struct {
    c_type type_info;
    ast_binary_op_type operation;
    ast_node* left;
    ast_node* right;
} binary_operation;

typedef struct {
    c_type type_info;
    ast_node* target;
    ast_node* value;
} assignment;

// Statements

typedef struct {
    ast_node* expression;
} statement_expression;

// Declarations

typedef struct {
    char* name;
    char** attributes;
    char* internal;

    c_type type_info;

    ast_variable_kind kind;

    int scope;
    uint32_t basin_size;

    bool requires_const;
    bool _register;

    vector* parameters;
    
    vector* children;
} variable;

typedef struct {
    char* string;
} assembly;

// Broad types

struct ast_node {
    ast_node_type type;
    union {
        integer_literal integer_literal;
        string_literal string_literal;

        binary_operation binary_operation;
        assignment assignment;

        assembly assembly;
        statement_expression statement_expression;

        variable variable;
    } content;
};

typedef struct {
    vector* declarations;
} ast;
