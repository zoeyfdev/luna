#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#ifndef _WIN32
    #include <dlfcn.h>
#else
    #include <libloaderapi.h> 
#endif

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

void* return_component_function(void* handle, char* name) {
#ifndef _WIN32
    void* function_handle = dlsym(handle, name);
#else
    void* function_handle = GetProcAddress(handle, name);
#endif

    if (function_handle == NULL) {
#ifndef _WIN32
        printf("luna-l2: failed to return function '%s': %s\n", name, dlerror());
#else
        printf("luna-l2: failed to return function '%s': failed\n", name);
#endif
    }
    return function_handle;
}

void* initialize_component(char* path) {
#ifndef _WIN32
    void* handle = dlopen(path, RTLD_LAZY);
#else
    void* handle = LoadLibraryA(path);
#endif

    if (handle == NULL) {
#ifndef _WIN32
        printf("luna-l2: failed to initialize component with path '%s': %s\n", path, dlerror());
#else
        printf("luna-l2: failed to initialize component with path '%s': failed\n", path);
#endif 
        exit(1);
    }

    return handle;
}
