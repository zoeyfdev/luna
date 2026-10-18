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
    AST_STMT_EXPRESSION,

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

// Leaves

typedef struct {
    c_type type;
    char* value;
} integer_literal;

typedef struct {
    c_type type;
    char* value;
} string_literal;

// Expressions

typedef struct {
    c_type type;
    ast_binary_op_type operation;
    struct ast_node* left;
    struct ast_node* right;
} binary_operation;

typedef struct {
    c_type value;
    struct ast_node* target;
    struct ast_node* value;
} assignment;

// Statements

typedef struct {
    struct ast_node* statement;
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

typedef struct {
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
} ast_node;

typedef struct {
    vector* declarations;
} ast;
