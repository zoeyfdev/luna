#include <stdint.h>

extern unsigned char (*audio_read_memory)(uint32_t address);
extern void (*audio_write_memory)(uint32_t address, unsigned char value);
