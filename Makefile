SRC=./
CC=gcc
CCFLAGS=-Wall -Wextra -Wimplicit-fallthrough -std=gnu23 -Wno-type-limits -Wno-unused-parameter
EMUFLAGS=-Wimplicit-fallthrough

# No windows support (for now until I figure out how to get DLLs to work)

all: luna-l2 las lcc lcc1 l2ld
.PHONY: clean install l2ld lcc lcc1 las

luna-l2: $(SRC)/l2/*
	cd l2 && $(CC) \
		hardware/g1x/g1x.c \
		-shared -fPIC \
		$(CCFLAGS) \
		-o ../components/video/g1x.so -g
	cd l2 && $(CC) \
		hardware/pit/pit.c \
		-shared -fPIC \
		$(CCFLAGS) \
		-o ../components/pit/pit.so -g
	cd l2 && $(CC) \
		hardware/s1/s1.c \
		-shared -fPIC \
		$(shell sdl2-config --cflags --libs) \
		$(CCFLAGS) \
		-o ../components/audio/s1.so -g
	cd l2 && $(CC) \
		luna_l2.c \
		util/*.c \
		cpu/*.c \
		bios/*.c \
		video/*.c \
		hwinit/*.c \
		io/*.c \
		component/*.c \
		-o ../bin/luna-l2 \
		$(shell sdl2-config --cflags --libs) \
		$(CCFLAGS) \
		-g

las: $(SRC)/las/* $(SRC)/lcc_info/*
	cd las && go build -o ../bin/las ./las.go

lcc1: $(SRC)/lcc1/* $(SRC)/lcc_info/*
	cd lcc1 && go build -o ../bin/lcc1 ./lcc1.go

lcc1-libs:
	cd lcc1/libs && lcc -c memcpy16.s
	cd lcc1/libs && lcc -c memcpy32.s
	cd lcc1/libs && lcc -c strcpy16.s
	cd lcc1/libs && lcc -c strcpy32.s
	cd lcc1/libs && sudo mv *.o /usr/local/lib/l2ld/
	sudo printf "_builtin_lcc_memcpy16 /usr/local/lib/l2ld/memcpy16.o\n_builtin_lcc_memcpy32 /usr/local/lib/l2ld/memcpy32.o\n" > /usr/local/lib/l2ld/memcpy.lib
	sudo printf "_builtin_lcc_strcpy32 /usr/local/lib/l2ld/strcpy32.o\n_builtin_lcc_strcpy16 /usr/local/lib/l2ld/strcpy16.o\n" > /usr/local/lib/l2ld/strcpy.lib

lcc: $(SRC)/lcc/* $(SRC)/lcc_info/*
	cd lcc && go build -o ../bin/lcc ./lcc.go

l2ld:
	cd l2ld && go build -o ../bin/l2ld ./l2ld.go

macos-installer:
	cd l2 && CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o ../Mac/amd64/usr/local/bin/"Luna L2"/Contents/MacOS/luna-l2 luna_l2.go
	cd lcc && GOOS=darwin GOARCH=amd64 go build -o ../Mac/amd64/usr/local/bin/lcc lcc.go
	cd las && GOOS=darwin GOARCH=amd64 go build -o ../Mac/amd64/usr/local/bin/las las.go
	cd lcc1 && GOOS=darwin GOARCH=amd64 go build -o ../Mac/amd64/usr/local/bin/lcc1 lcc1.go
	cd l2ld && GOOS=darwin GOARCH=amd64 go build -o ../Mac/amd64/usr/local/bin/l2ld l2ld.go
	cd l2 && CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o ../Mac/arm64/usr/local/bin/"Luna L2"/Contents/MacOS/luna-l2 luna_l2.go
	cd lcc && GOOS=darwin GOARCH=arm64 go build -o ../Mac/arm64/usr/local/bin/lcc lcc.go
	cd las && GOOS=darwin GOARCH=arm64 go build -o ../Mac/arm64/usr/local/bin/las las.go
	cd lcc1 && GOOS=darwin GOARCH=arm64 go build -o ../Mac/arm64/usr/local/bin/lcc1 lcc1.go
	cd l2ld && GOOS=darwin GOARCH=arm64 go build -o ../Mac/arm64/usr/local/bin/l2ld l2ld.go
	pkgbuild \
		--root Mac/amd64 \
		--install-location / \
		--identifier com.alexfdev0.lunal2.amd64 \
		--version 1.0 \
		--scripts Mac/scripts \
		build/"Luna L2 (amd64).pkg"
	pkgbuild \
		--root Mac/arm64 \
		--install-location / \
		--identifier com.alexfdev0.lunal2.arm64 \
		--version 1.0 \
		--scripts Mac/scripts \
		build/"Luna L2 (arm64).pkg"

mac_qmake:
	mkdir -p /usr/local/lib/l2/
	cd l2 && go build -buildmode=plugin -o ../components/audio/s1.so ./audio/hardware/s1.go
	cd l2 && go build -buildmode=plugin -o ../components/video/g1x.so ./video/hardware/g1x.go
	cd l2 && go build -buildmode=plugin -o ../components/video/g1.so ./video/hardware/g1.go
	sudo cp -r components/* /usr/local/lib/l2/
	cd l2 && CGO_ENABLED=1 GOOS=darwin go build -o /Applications/"Luna L2.app"/Contents/MacOS/luna-l2 luna_l2.go
	cd lcc && GOOS=darwin go build -o /usr/local/bin/lcc lcc.go
	cd las && GOOS=darwin go build -o /usr/local/bin/las las.go
	cd lcc1 && GOOS=darwin go build -o /usr/local/bin/lcc1 lcc1.go
	cd l2ld && GOOS=darwin go build -o /usr/local/bin/l2ld l2ld.go	

windows-installer:
	cd l2 && CGO_LDFLAGS="-lmingw32 -lSDL2" CGO_CFLAGS="-D_REENTRANT" CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui" -o ../Windows/luna-l2.exe -x luna_l2.go
	cd lcc && GOOS=windows GOARCH=amd64 go build -o ../Windows/lcc.exe lcc.go
	cd las && GOOS=windows GOARCH=amd64 go build -o ../Windows/las.exe las.go
	cd lcc1 && GOOS=windows GOARCH=amd64 go build -o ../Windows/lcc1.exe lcc1.go
	cd l2ld && GOOS=windows GOARCH=amd64 go build -o ../Windows/l2ld.exe l2ld.go
	cd Windows && wixl -v msi.xml -o "../build/Luna L2.msi"

install: lcc1-libs
	mkdir -p /usr/local/lib/l2/
	sudo cp -r components/* /usr/local/lib/l2
	sudo cp bin/* /usr/local/bin/

clean:
	rm -f /usr/local/bin/luna-l2
	rm -f /usr/local/bin/las
	rm -f /usr/local/bin/lcc1
	rm -f /usr/local/bin/lcc
	rm -f /usr/local/bin/l2ld
