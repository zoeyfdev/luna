#include <stdio.h>

#include "../../parse/ast.h"
#include "../../parse/types.h"

vector* stream;

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

char* codegen_l2(ast* tu, bool is_32) {
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
        } else if (decl->type == AST_ASSEMBLY) {
            write(decl->content.assembly.string, false);
            write("    // User-defined inline assembly\n", false);
        }

        write("\n", false);
    }

    vec_grow(stream, 1);
    ((char*) stream->data)[stream->next] = 0;

    return (char*) stream->data;
}
