SRC=.
CC=gcc

OS_NAME := $(shell uname -s)
SDL_FLAGS = $(shell sdl2-config --cflags --libs)
CCFLAGS=-Wall -Wextra -std=gnu23 -Wno-type-limits -Wno-unused-parameter -Wno-implicit-fallthrough

ifeq ($(OS_NAME),Darwin)
	CC=clang
endif
ifeq ($(OS),Windows_NT)
	CC=x86_64-w64-mingw32-gcc
endif

all: luna-l2 luna-l2-components las lcc lcc1 l2ld
.PHONY: all luna-l2 luna-l2-components las lcc lcc1 l2ld

luna-l2-components: $(SRC)/l2/*
		cd l2 && $(CC) \
		hardware/g1x/g1x.c \
		-shared -fPIC \
		$(CCFLAGS) \
		-o ../components/g1x.so -g
	cd l2 && $(CC) \
		hardware/pit/pit.c \
		-shared -fPIC \
		$(CCFLAGS) \
		-o ../components/pit.so -g
	cd l2 && $(CC) \
		hardware/s1/s1.c \
		-shared -fPIC \
		$(SDL_FLAGS) \
		$(CCFLAGS) \
		-o ../components/s1.so -g

luna-l2: $(SRC)/l2/*
	cd l2 && $(CC) \
		*.c \
		util/*.c \
		cpu/*.c \
		bios/*.c \
		video/*.c \
		hwinit/*.c \
		io/*.c \
		component/*.c \
		-o ../bin/luna-l2 \
		$(SDL_FLAGS) \
		$(CCFLAGS) \
		-g

las: $(SRC)/las/* $(SRC)/lcc_info/*
	cd las && $(CC) \
		*.c \
		lexer/*.c \
		parse/*.c \
		error/*.c \
		util/*.c \
		../lcc_shared/libvector.c \
		../lcc_shared/shared.c \
		../lcc_shared/liberror.c \
		../lcc_shared/libfile.c \
		../lcc_shared/libstoi.c \
		../lcc_shared/libcpp.c \
		$(CCFLAGS) \
		-o ../bin/las \
		-g

lcc1: $(SRC)/lcc1/* $(SRC)/lcc_info/*
	cd lcc1 && go build -o ../bin/lcc1 ./lcc1.go

l2ld: $(SRC)/l2ld/*
	cd l2ld && $(CC) \
		*.c \
		../lcc_shared/libvector.c \
		../lcc_shared/libfile.c \
		-o ../bin/l2ld \
		-g \
		$(CCFLAGS)

lcc: $(SRC)/lcc/* $(SRC)/lcc_shared/*
	cd lcc && $(CC) \
		*.c \
		../lcc_shared/libvector.c \
		../lcc_shared/libfile.c \
		../lcc_shared/liberror.c \
		../lcc_shared/libcommand.c \
		../lcc_shared/shared.c \
		-o ../bin/lcc \
		-g \
		$(CCFLAGS)

lcc1-libs:
	cd lcc1/libs && lcc -c memcpy16.s
	cd lcc1/libs && lcc -c memcpy32.s
	cd lcc1/libs && lcc -c strcpy16.s
	cd lcc1/libs && lcc -c strcpy32.s
	cd lcc1/libs && sudo mv *.o /usr/local/lib/l2ld/
	sudo printf "_builtin_lcc_memcpy16 /usr/local/lib/l2ld/memcpy16.o\n_builtin_lcc_memcpy32 /usr/local/lib/l2ld/memcpy32.o\n" > /usr/local/lib/l2ld/memcpy.lib
	sudo printf "_builtin_lcc_strcpy32 /usr/local/lib/l2ld/strcpy32.o\n_builtin_lcc_strcpy16 /usr/local/lib/l2ld/strcpy16.o\n" > /usr/local/lib/l2ld/strcpy.lib

mac_qmake: luna-l2 lcc las lcc1 l2ld
	sudo cp bin/luna-l2 /Applications/"Luna L2.app"/Contents/MacOS/
	rm bin/luna-l2
	sudo cp bin/* /usr/local/bin

install:
	mkdir -p /usr/local/lib/l2/
	sudo cp -r components/* /usr/local/lib/l2
	sudo cp bin/* /usr/local/bin/

quick: # atomic install
	sh -c "make"
	sh -c "sudo make install"

clean:
	rm -rf bin/*
	rm -rf components/*
