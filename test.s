mov sp, 0xE000

push string
call print
call _builtin_lcc_strcpy16
jmp pc

string:
    .asciz "Hello world!"
