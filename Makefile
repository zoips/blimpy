current_dir := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
gopath := $(current_dir)/src:$(current_dir)/_vendor

.PHONY: all build deps clean

all: deps build

build:
	GOPATH=$(gopath) go build -o $(current_dir)/bin/blimpy $(current_dir)/src/main.go

deps:
	GOPATH=$(gopath) $(current_dir)/bin/gom install

clean:
	rm $(current_dir)/bin/blimpy

