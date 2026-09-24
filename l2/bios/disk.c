#include <stdio.h>
#include <fcntl.h>
#include <unistd.h>

#include "bios.h"
#include "../util/psleep.h"
#include "../cpu/registers.h"
#include "../cpu/memory.h"

bool load_sector(int drive, int sector, int dest_sector) {
    char* file;
    switch (drive) {
    case 0:
        file = HDD_FILE;
        psleep(12);
        break;
    case 1:
        file = SD_FILE;
        psleep(2);
        break;
    case 2:
        file = DVD_FILE;
        psleep(110);
        break;
    }

    if (file == NULL) {
        printf("luna-l2: invalid drive %d for disk!\n", drive);
        return false;
    }

    int fd = open(file, O_RDONLY, 0);
    if (fd == -1) {
        printf("luna-l2: could not load/reload block device with drive ID %d\n", drive);
        set_register(R0, 1);
        return false;
    }

    pread(fd, MEMORY + (sector * 512), 512, dest_sector * 512);

    close(fd);

    set_register(R0, 0);
    return true;
}

bool write_sector(int drive, int sector, int dest_sector) {
    char* file;
    switch (drive) {
    case 0:
        file = HDD_FILE;
        psleep(12);
        break;
    case 1:
        file = SD_FILE;
        psleep(2);
        break;
    case 2:
        file = DVD_FILE;
        psleep(110);
        break;
    }

    if (file == NULL) {
        printf("luna-l2: invalid drive %d for disk!\n", drive);
        return false;
    }

    int fd = open(file, O_RDWR | O_SYNC, 0);
    if (fd == -1) {
        printf("luna-l2: could not load/reload block device with drive ID %d\n", drive);
        set_register(R0, 1);
        return false;
    }

    pwrite(fd, MEMORY + (sector * 512), 512, dest_sector * 512);

    close(fd);

    set_register(R0, 0);
    return true;
}
