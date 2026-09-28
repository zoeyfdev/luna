mov sp, 0xEFFF

push string
call print
jmp pc

string:
    .asciz "Hello world!"
