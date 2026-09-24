#include <stdint.h>

#include "../component/component.c"
#include "../video/nrgba/nrgba.h"

unsigned char (*v_read_video_memory)(uint32_t);
void (*v_write_video_memory)(uint32_t, unsigned char);
void (*v_print_char)(unsigned char, unsigned char, unsigned char);
void (*v_set_cursor)(int, int);
void (*v_get_cursor)(int*, int*);
void (*v_gpu_reset)();
nrgba_image* (*return_framebuffer)();

void initialize_video() {
    void* video_component = initialize_component("/usr/local/lib/l2/video/g1x.so");

    nrgba_image* (*return_framebuffer)();
    void (*initialize_component)();
 
    initialize_component = return_component_function(video_component, "initialize_component");

    v_read_video_memory = return_component_function(video_component, "read_video_memory");
    v_write_video_memory = return_component_function(video_component, "write_video_memory");
    v_print_char = return_component_function(video_component, "print_char");
    v_set_cursor = return_component_function(video_component, "set_cursor");
    v_get_cursor = return_component_function(video_component, "get_cursor");
    v_gpu_reset = return_component_function(video_component, "gpu_reset");
    return_framebuffer = return_component_function(video_component, "return_framebuffer");

    (*initialize_component)();
}
