mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))

all:
	GOPATH=$(current_dir)/src:$(current_dir)/_vendor go build -o bin/blimpy src/main.go

.PHONY: all
