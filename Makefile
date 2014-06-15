current_dir := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
gopath := $(current_dir)/src:$(current_dir)/_vendor

.PHONY: all build deps clean init

all: init deps build

build:
	GOPATH=$(gopath) go build -o $(current_dir)/bin/blimpy $(current_dir)/src/main.go

deps:
	GOPATH=$(current_dir)/_vendor $(current_dir)/_vendor/bin/gom install

clean:
	rm $(current_dir)/bin/blimpy

init:
	GOPATH=$(current_dir)/_vendor go get github.com/mattn/gom
