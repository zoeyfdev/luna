CC=lcc

.PHONY = portable

all:
	$(CC) meteor.s -o meteor.img

portable:
	$(CC) -fpie -fpie-32 meteor.s -o meteor-portable.img -define PORTABLE 1
