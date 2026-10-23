#include <stdio.h>

#ifndef _WIN32
    #include <fcntl.h>
#else
    #include <fileapi.h>
    #include <minwinbase.h>
    #include <handleapi.h>
#endif

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
        set_register(R0, 1);
        return false;
    }

#ifndef _WIN32
    int fd = open(file, O_RDONLY, 0);
    bool fail = fd == -1;
#else
    void* fd = CreateFileA(file, GENERIC_READ, FILE_SHARE_READ,
            NULL, OPEN_EXISTING, FILE_ATTRIBUTE_NORMAL, NULL);
    OVERLAPPED* offset = calloc(1, sizeof(OVERLAPPED));
    offset->OffsetHigh = (((dest_sector * 512) & 0xFFFFFFFF00000000) >> 32);
    offset->Offset = ((dest_sector * 512) & 0xFFFFFFFF);

    bool fail = fd == INVALID_HANDLE_VALUE;
#endif

    if (fail) {
        printf("luna-l2: could not load/reload block device with drive ID %d\n", drive);
        set_register(R0, 1);
        return false;
    }

#ifndef _WIN32
    pread(fd, MEMORY + (sector * 512), 512, dest_sector * 512);
    close(fd);
#else
    ReadFile(fd, MEMORY + (sector * 512), 512, NULL, offset);
    CloseHandle(fd);
    free(offset);
#endif

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
        set_register(R0, 1);
        return false;
    }

#ifndef _WIN32
    int fd = open(file, O_RDWR | O_SYNC, 0);
    bool fail = fd == -1;
#else
    void* fd = CreateFileA(file, GENERIC_READ | GENERIC_WRITE, 
            FILE_SHARE_READ | FILE_SHARE_WRITE,
            NULL, OPEN_EXISTING, FILE_ATTRIBUTE_NORMAL, NULL);
    OVERLAPPED* offset = calloc(1, sizeof(OVERLAPPED));
    offset->OffsetHigh = (((dest_sector * 512) & 0xFFFFFFFF00000000) >> 32);
    offset->Offset = ((dest_sector * 512) & 0xFFFFFFFF);

    bool fail = fd == INVALID_HANDLE_VALUE;
#endif

    if (fail) {
        printf("luna-l2: could not load/reload block device with drive ID %d\n", drive);
        set_register(R0, 1);
        return false;
    }

#ifndef _WIN32
    pwrite(fd, MEMORY + (sector * 512), 512, dest_sector * 512);
    close(fd);
#else
    WriteFile(fd, MEMORY + (sector * 512), 512, NULL, offset);
    CloseHandle(fd);
    free(offset);
#endif

    set_register(R0, 0);
    return true;
}
