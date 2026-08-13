CC=lcc

.PHONY = portable

all:
	$(CC) meteor.s -o meteor.img

portable:
	$(CC) -fpie -fpie-32 meteor-portable.s -o meteor-portable.img
