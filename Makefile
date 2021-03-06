current_dir := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
gopath := $(current_dir)/../../../../:$(current_dir)/_vendor
match := $(match)

.PHONY: all build deps clean init

all: init deps build

build:
	GOPATH=$(gopath) go build -o $(current_dir)/bin/blimpy $(current_dir)/server/main.go

test:
ifdef match
	GOPATH=$(gopath) go test -cpu=4 github.com/zoips/blimpy -run=$(match)
else
	GOPATH=$(gopath) go test -cpu=4 github.com/zoips/blimpy
endif

deps: init
	GOPATH=$(current_dir)/_vendor $(current_dir)/_vendor/bin/gom install

clean:
	rm $(current_dir)/bin/blimpy

fmt:
	GOPATH=$(current_dir) go fmt .

init:
	GOPATH=$(current_dir)/_vendor go get github.com/mattn/gom
