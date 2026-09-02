.bits 32
.global CRASH_SOUND
.global BOOT_SOUND

CRASH_SOUND: 
    .embed "audio/crash.pcm.fz"
BOOT_SOUND:
    .embed "audio/boot.pcm.fz"
