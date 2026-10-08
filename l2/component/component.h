#pragma once

extern char* component_base;

void* initialize_component(char* path);
void* return_component_function(void* handle, char* name);
char* return_component_path(char* name);
