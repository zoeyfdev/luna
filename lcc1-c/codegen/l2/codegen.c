#include <stdio.h>

#include "../../parse/ast.h"
#include "../../parse/types.h"
#include "../../../lcc_shared/libvector.h"
#include "registers.h"
#include "types.h"
#include "util.h"
#include "../../../lcc_shared/liberror.h"

vector* header;
vector* stream;

void write_header(char* buf, bool indent) {
    if (indent)
        for (int i = 0; i < 4; i++) {
            vec_grow(header, 1);
            ((char*) header->data)[header->next] = 0x20;
        }

    while (*buf) {
        vec_grow(header, 1);
        ((char*) header->data)[header->next] = *buf;
        buf++;
    }
}

void write(char* buf, bool indent) {
    if (indent)
        for (int i = 0; i < 4; i++) {
            vec_grow(stream, 1);
            ((char*) stream->data)[stream->next] = 0x20;
        }

    while (*buf) {
        vec_grow(stream, 1);
        ((char*) stream->data)[stream->next] = *buf;
        buf++;
    }
}

codegen_return* codegen_leaf(ast_node* leaf, bool read) {
    if (leaf->type == AST_INTEGER_LITERAL) {
        char* _register = take_register();

        write("mov ", true);
        write(_register, false);
        write(", ", false);
        write(leaf->content.integer_literal.value, false);
        write("\n", false);

        codegen_return* r = new_cgr();

        r->_register = _register;
        r->type_info = leaf->content.integer_literal.type_info;
        return r;
    } else if (leaf->type == AST_STRING_LITERAL) {
        char* name = generate_generic();
        write_header(name, false);
        write_header(":\n", false);

        write_header(".asciz \"", true);
        write_header(leaf->content.string_literal.value, false);
        write_header("\"\n\n", false);

        char* name2 = generate_generic();
        write_header(name, false);
        write_header(":\n", false);

        write_header(".ptr ", true);
        write_header(name, false);
        write_header("\n\n", false);

        char* r = take_register();
        write("mov ", true);
        write(r, false);
        write(", ", false);
        write(name2, false);
        write("\n", false);

        if (read)



        codegen_return* ret = new_cgr();
        ret->_register = r;
        ret->type_info = leaf->content.string_literal.type_info;
        return ret;
    }

    lcc_ice(NULL, 0, "no return value for codegen_leaf()", NULL);
}

codegen_return* codegen_expression(ast_node* expression, bool read) {
    if (expression->type == AST_INTEGER_LITERAL ||
            expression->type == AST_STRING_LITERAL) {
        return codegen_leaf(expression, read);
    } else if (expression->type == AST_ASSIGNMENT) {
        codegen_return* target = codegen_expression(expression->content.assignment.target, false);
        codegen_return* value = codegen_expression(expression->content.assignment.value, true);
       
        emit_store(&value->type_info, target->_register, value->_register);
        free_register(target->_register);

        return value;
    }

    lcc_ice(NULL, 0, "no return value for codegen_expression()", NULL);
}

void codegen_statement(ast_node* statement) {
    if (statement->type == AST_ASSEMBLY_STATEMENT) {
        write(statement->content.assembly.string, true);
        write("    // Inline assembly\n", false);
    } else if (statement->type == AST_STATEMENT_EXPRESSION) {
        free_register(codegen_expression(statement->content.statement_expression.expression, true)->_register);
    }
}

void codegen_decl(ast_node* decl) {
    if (decl->type == AST_ASSEMBLY_DECL) {
        write_header(decl->content.assembly.string, false);
        write_header("    // Inline assembly\n", false);
    } else if (decl->type == AST_VARIABLE) {
        variable v = decl->content.variable;

        if (!v.type_info.external) { 
            if (v.kind == KIND_VARIABLE) {
                write_header(v.name, false);
                write_header(":\n", false);
                if ((v.scope == 0 || v.type_info._static) && !v._register) {
                    // const codegen
                    switch (v.type_info.type) {
                    case TYPE_I8:
                        write_header(".byte 0x00\n", true);
                        break;
                    case TYPE_I16:
                        write_header(".word 0x0000\n", true);
                        break;
                    case TYPE_I32:
                        write_header(".dword 0x00000000\n", true);
                        break;
                    }
                    write_header("\n", false);
                } 
            } else if (v.kind == KIND_FUNCTION) {
                write(v.name, false);
                write(":\n", false);
                for (int i = 0; i < v.children->elements; i++) {
                    ast_node* statement = ((ast_node**) v.children->data)[i];
                    codegen_statement(statement);
                }
            }
        }
    } 
}

char* codegen_l2(ast* tu, bool is_32) {
    regs_init();
    header = vec_init(sizeof(char*), 0);
    stream = vec_init(sizeof(char*), 0);

    if (!is_32)
        write_header(".bits 16\n\n", false);
    else
        write_header(".bits 32\n\n", false);

   
    for (int i = 0; i < tu->declarations->elements; i++) {
        ast_node* decl = ((ast_node**) tu->declarations->data)[i];
 
        if (decl->type == AST_VARIABLE) {
            if (!decl->content.variable.type_info._static && !decl->content.variable.type_info.external) {
                write_header(".global ", false);
                write_header(decl->content.variable.name, false);
                write_header("\n", false);
            }
        }
    }
    write_header("\n", false);

    for (int i = 0; i < tu->declarations->elements; i++) {
        ast_node* decl = ((ast_node**) tu->declarations->data)[i];
 
        codegen_decl(decl);
    }

    vec_grow(stream, 1);
    ((char*) stream->data)[stream->next] = 0;

    vector* final = vec_init(sizeof(char*), 0);
    for (int i = 0; i < header->elements; i++) {
        vec_grow(final, 1);
        ((char*) final->data)[final->next] = ((char*) header->data)[i];
    }
    for (int i = 0; i < stream->elements; i++) {
        vec_grow(final, 1);
        ((char*) final->data)[final->next] = ((char*) stream->data)[i];
    }

    return (char*) final->data;
}
