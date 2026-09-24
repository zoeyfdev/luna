#ifndef _WIN32

#include <dlfcn.h>
#include <stdio.h>
#include <stdlib.h>

void* return_component_function(void* handle, char* name) {
    void* function_handle = dlsym(handle, name);
    if (function_handle == NULL) {
        printf("luna-l2: failed to return function '%s': %s", name, dlerror());
    }
    return function_handle;
}

void* initialize_component(char* path) {
    void* handle = dlopen(path, RTLD_LAZY);
    if (handle == NULL) {
        printf("luna-l2: failed to initialize component with path '%s': %s\n", path, dlerror());
        exit(1);
    }

    return handle;
}


#endif
