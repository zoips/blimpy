current_dir := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
gopath := $(current_dir)/src:$(current_dir)/_vendor

.PHONY: build deps clean

deps:
	GOPATH=$(gopath) $(current_dir)/bin/gom install

clean:
	rm $(current_dir)/bin/blimpy

build: deps
	GOPATH=$(gopath) go build -o $(current_dir)/bin/blimpy $(current_dir)/src/main.go


