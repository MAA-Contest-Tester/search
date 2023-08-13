UNAME := $(shell uname)
ifeq ($(UNAME), Darwin)
	FLAGS=-ldflags '-s -extldflags "-sectcreate __TEXT __info_plist Info.plist"'
else
	FLAGS=-tags netgo -ldflags '-extldflags "-static"'
endif

out/psearch:
	mkdir -p out
	go build $(FLAGS) -o out/psearch .

clean:
	rm -rf out/psearch

all: out/psearch
make: out/psearch
