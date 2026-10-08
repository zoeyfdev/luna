#include <stdint.h>
#include <stdlib.h>
#include <string.h>

#include "../cpu/memory.h"
#include "../component/component.h"

unsigned char (*audio_read_memory)(uint32_t address);
void (*audio_write_memory)(uint32_t address, unsigned char value);

void initialize_audio(char* component) {
    void* audio_component = initialize_component(return_component_path(component));

    void (*audio_init)(unsigned char (*)(uint32_t)) = return_component_function(audio_component, "audio_init");

    audio_init(get_memory);

    audio_read_memory = return_component_function(audio_component, "audio_read_memory");
    audio_write_memory = return_component_function(audio_component, "audio_write_memory");
}
