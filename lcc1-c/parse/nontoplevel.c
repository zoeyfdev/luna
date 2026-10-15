#include "ast.h"
#include "types.h"
#include "../../lcc_shared/libvector.h"

void parse_local(vector* elements, int start) {
    for (int i = start; i < tokens->elements; i++) {
        token* t = ((token**) tokens->data)[i];
    }
}
