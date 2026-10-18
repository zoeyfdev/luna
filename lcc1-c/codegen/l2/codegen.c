#include <stdio.h>

#include "../../parse/ast.h"
#include "../../parse/types.h"
#include "../../../lcc_shared/libvector.h"
#include "registers.h"
#include "types.h"

vector* stream;

codegen_return* write(char* buf, bool indent) {
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

codegen_return* codegen_leaf(ast_node* leaf) {
    if (leaf->type == AST_INTEGER_LITERAL) {
        
    }
}

codegen_return* codegen_expression(ast_node* expression) {
     
}

void codegen_statement(ast_node* statement) {
    if (statement->type == AST_ASSEMBLY_STATEMENT) {
        write(statement->content.assembly.string, true);
        write("    // Inline assembly\n", true);
    } else if (statement->type == AST_STATEMENT_EXPRESSION) {
        codegen_expression(statement->content.statement_expression.expression);
    }
}

void codegen_decl(ast_node* decl) {
    if (decl->type == AST_ASSEMBLY_DECL) {
        write(decl->content.assembly.string, false);
        write("    // Inline assembly\n", false);
    } else if (decl->type == AST_VARIABLE) {
        variable v = decl->content.variable;

        if (!v.type_info.external) {
            write(v.name, false);
            write(":\n", false);
            if (v.kind == KIND_VARIABLE) {
                if ((v.scope == 0 || v.type_info._static) && !v._register) {
                    // const codegen
                    switch (v.type_info.type) {
                    case TYPE_I8:
                        write(".byte 0x00\n", true);
                        break;
                    case TYPE_I16:
                        write(".word 0x0000\n", true);
                        break;
                    case TYPE_I32:
                        write(".dword 0x00000000\n", true);
                        break;
                    }
                    write("\n", false);
                } 
            } else if (v.kind == KIND_FUNCTION) {
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
    stream = vec_init(sizeof(char*), 0);

    if (!is_32)
        write(".bits 16\n", false);
    else
        write(".bits 32\n", false);

    
    for (int i = 0; i < tu->declarations->elements; i++) {
        ast_node* decl = ((ast_node**) tu->declarations->data)[i];
 
        if (decl->type == AST_VARIABLE) {
            if (decl->content.variable.kind == KIND_FUNCTION) {
                if (!decl->content.variable.type_info._static && !decl->content.variable.type_info.external) {
                    write(".global ", false);
                    write(decl->content.variable.name, false);
                    write("\n", false);
                }
            }
        }
        codegen_decl(decl);

        write("\n", false);
    }

    vec_grow(stream, 1);
    ((char*) stream->data)[stream->next] = 0;

    return (char*) stream->data;
}
