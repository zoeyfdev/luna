#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#ifndef _WIN32
    char* component_base = "/usr/local/lib/l2/";
#else
    char* component_base = "C:\\Program Files\\Luna L2\\lib\\l2\\";
#endif

char* return_component_path(char* name) {
    char* path = calloc(1, strlen(component_base) + strlen(name) + 6);
    strcat(path, component_base);
    strcat(path, name);
#ifndef _WIN32
    strcat(path, ".so");
#else
    strcat(path, ".dll");
#endif

    return path;
}

#ifndef _WIN32

#include <dlfcn.h>

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
