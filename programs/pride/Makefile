CC=lcc

.PHONY = portable

all:
	$(CC) -v entry.s main.c util.s flag.s -o pride.img

portable:
	$(CC) -fpie -fpie-32 -define PORTABLE 1 entry.s main.c util.s flag.s -o pride-portable.img 
