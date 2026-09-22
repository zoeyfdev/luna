#include <stdint.h>

extern char* HDD_FILE;
extern char* SD_FILE;
extern char* DVD_FILE;
extern int bios_boot_drive;
extern unsigned char* BIOS_RAM;

void bios_splash();
void bios_write_line(char* s);
void bios_handle_interrupt(uint32_t code);
void bios_init();
unsigned char bios_read_memory(uint32_t address);
void bios_write_memory(uint32_t address, unsigned char value);
