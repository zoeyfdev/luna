SRC=.
CC=gcc

BUILD_MODE=
MODULE_ENDING=so
WIN_COMPONENT_ROOT=
DEBUG=

OS_NAME := $(shell uname -s)
SDL_FLAGS = $(shell pkg-config --cflags --libs sdl2)
CCFLAGS=-Wall -Wextra -std=gnu23 -Wno-type-limits -Wno-unused-parameter -Wno-implicit-fallthrough -Wno-switch

ifeq ($(DEBUG),true)
	CCFLAGS+=-g
endif
ifeq ($(BUILD_MODE),windows)
	CC=x86_64-w64-mingw32-gcc
	SDL_FLAGS=-Wl,-Bstatic $(shell x86_64-w64-mingw32-pkg-config sdl2 --libs --cflags) -Wl,-Bdynamic -lsetupapi -lole32 -lwinmm -loleaut32 -limm32 -lversion 
	CCFLAGS+=-Wl,-Bstatic -lgcc -lwinpthread -Wl,-Bdynamic
	ifneq ($(strip $(WIN_COMPONENT_ROOT)),)
		CCFLAGS+=-DWIN_COMPONENT_BASE=$(WIN_COMPONENT_ROOT)
	endif
	MODULE_ENDING=dll
endif

all: always luna-l2 luna-l2-components las lcc lcc1 l2ld
.PHONY: all luna-l2 luna-l2-components las lcc lcc1 l2ld

always:
	mkdir -p bin
	mkdir -p components

luna-l2-components:
		cd l2 && $(CC) \
		hardware/g1x/g1x.c \
		-shared -fPIC \
		$(CCFLAGS) \
		-o ../components/g1x.$(MODULE_ENDING)
	cd l2 && $(CC) \
		hardware/pit/pit.c \
		-shared -fPIC \
		$(CCFLAGS) \
		-o ../components/pit.$(MODULE_ENDING)
	cd l2 && $(CC) \
		hardware/s1/s1.c \
		-shared -fPIC \
		$(SDL_FLAGS) \
		$(CCFLAGS) \
		-o ../components/s1.$(MODULE_ENDING)

luna-l2:
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
		$(CCFLAGS)

las:
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
		-o ../bin/las

lcc1:
	cd lcc1-c && $(CC) \
		*.c \
		lex/*.c \
		parse/*.c \
		codegen/l2/*.c \
		env/*.c \
		typecheck/*.c \
		../lcc_shared/shared.c \
		../lcc_shared/libvector.c \
		../lcc_shared/libcpp.c \
		../lcc_shared/liberror.c \
		../lcc_shared/libstoi.c \
		../lcc_shared/libfile.c \
		-o ../bin/lcc1-c \
		$(CCFLAGS)
	#cd lcc1 && go build -o ../bin/lcc1 ./lcc1.go

l2ld:
	cd l2ld && $(CC) \
		*.c \
		../lcc_shared/libvector.c \
		../lcc_shared/libfile.c \
		-o ../bin/l2ld \
		$(CCFLAGS)

lcc:
	cd lcc && $(CC) \
		*.c \
		../lcc_shared/libvector.c \
		../lcc_shared/libfile.c \
		../lcc_shared/liberror.c \
		../lcc_shared/libcommand.c \
		../lcc_shared/shared.c \
		-o ../bin/lcc \
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
ifneq ($(BUILD_MODE),windows)
	@mkdir -p /usr/local/lib/l2/
	@sudo cp -r components/* /usr/local/lib/l2
	@sudo cp bin/* /usr/local/bin/
else
	@mkdir -p C:\\"Program Files"\\"Luna L2"\\
	@mkdir -p C:\\"Program Files"\\"Luna L2"\\lib
	@mkdir -p C:\\"Program Files"\\"Luna L2"\\lib\\l2
	@mkdir -p C:\\"Program Files"\\"Luna L2"\\bin
	@mkdir -p C:\\"Program Files"\\"Luna L2"\\include
	@cp -r bin/* C:\\"Program Files"\\"Luna L2"\\bin\\
	@cp -r components/* C:\\"Program Files"\\"Luna L2"\\lib\\l2\\
	@setx PATH "%PATH%;C:\Program Files\Luna L2\bin" /M
endif

quick: # atomic install
	sh -c "make"
	sh -c "sudo make install"

clean:
	rm -rf bin/*
	rm -rf components/*
