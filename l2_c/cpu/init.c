#include <stdlib.h>
#include <stdio.h>

#include "../video/video.h"
#include "../util/psleep.h"
#include "../bios/bios.h"
#include "../bios/disk.h"

#include "registers.h"
#include "memory.h"
#include "cpu.h"

void* cpu_poweron(void* VOID) {
    initialize_registers();
    initialize_memory();

    while (!VIDEO_READY)
        psleep(15);

    bios_splash();
    
    int drive = 0;
    bool boot_ok = false;

    // try drives
    bool result;

boot_try:
    switch (drive) {
    case 0:
        bios_write_line("Booting from hard disk...");
        result = load_sector(0, 0, 0);
        break;
    case 1:
        bios_write_line("Booting from SD...");
        result = load_sector(1, 0, 0);
        break;
    case 2:
        bios_write_line("Booting from DVD...");
        result = load_sector(2, 0, 0);
        break;
    }

    if (result == false) {
        if (drive != 2) {
            bios_write_line("Could not read the boot disk\n");
            drive++;
            goto boot_try;
        } else {
            bios_write_line("No bootable device");
        }
    } else {
        boot_ok = true;
        bios_boot_drive = drive;
    }

    if (boot_ok == true)
        cpu_execute();
    else
        for (;;) psleep(15);
    return NULL;
}


