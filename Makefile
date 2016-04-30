export GOPATH := $(shell pwd)
default: build

init:
	bower install
	cd src/main && go get

clean:
	rm -f bin/server bin/main bin/LA-Hacks

build: init clean
	go build -o bin/LA-Hacks src/main/main.go 

run: build
	@-pkill LA-Hacks
	bin/LA-Hacks>>log.txt 2>&1 &
