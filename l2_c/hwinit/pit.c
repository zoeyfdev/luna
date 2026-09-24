#include <stdint.h>

#include "../component/component.h"
#include "../cpu/registers.h"

unsigned char (*pit_read_memory)(uint32_t);
void (*pit_write_memory)(uint32_t, unsigned char);

void initialize_pit() {
    void* pit_component = initialize_component("/usr/local/lib/l2/pit/pit.so");
    void (*pit_init)(void (*)(unsigned char, uint32_t), uint32_t (*)(unsigned char)) = return_component_function(pit_component, "pit_init");

    pit_init(set_register, get_register);

    pit_read_memory = return_component_function(pit_component, "pit_read_memory");
    pit_write_memory = return_component_function(pit_component, "pit_write_memory");
}
