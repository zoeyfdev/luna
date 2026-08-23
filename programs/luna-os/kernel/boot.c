#pragma bits 32

#include "stdlib.h"
#include "images.h"
#include "fzip.h"
#include "audio.h"
#include "stub.h"
#include "stdbool.h"
#include "lufs.h"
#include "util.h"

void boot() __attribute__((noreturn)) {
    targeted_load((long int) puts32_loc, 2);

    puts32("Loading resources...\n", COLOR_WHITE, COLOR_BLACK);
    targeted_load((long int) BOOT_SOUND, 43);
    targeted_load((long int) BOOT_IMG, 4);
    targeted_load((long int) play_sound_loc, 3);
    targeted_load((long int) renderbuf_loc, 2);
    targeted_load((long int) sleep_loc, 2);
    targeted_load((long int) fzipdecode_loc, 3);
    targeted_load((long int) malloc_loc, 2);

    play_sound((void*) fzip_decode(BOOT_SOUND), 41984, false);
    render_buf((void*) fzip_decode(BOOT_IMG));

    puts32("Loading LunaOS...\n", COLOR_WHITE, COLOR_BLACK);

    linear_sector_load(0xFF);

    puts32("Bootstrapping LUFS...\n", COLOR_WHITE, COLOR_BLACK);

    fstrap();

    render_buf((void*) 0x30303030);
    video_set_cursor(0, 0);
    asm ("jmp _cstart");
}
