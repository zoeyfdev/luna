#pragma once

extern unsigned char (*pit_read_memory)(uint32_t);
extern void (*pit_write_memory)(uint32_t, unsigned char);
