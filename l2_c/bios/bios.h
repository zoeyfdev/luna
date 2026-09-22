#include <stdint.h>

extern char* HDD_FILE;
extern char* SD_FILE;
extern char* DVD_FILE;
extern int bios_boot_drive;

void bios_splash();
void bios_write_line(char* s);
void bios_handle_interrupt(uint32_t code);
