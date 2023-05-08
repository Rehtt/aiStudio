cur_makefile_path := $(abspath $(lastword ./))
binName = $(shell echo $(cur_makefile_path)|awk -F '/' '{ print $$NF }')

ifeq ($(shell uname),Windows_NT)
	suffix = .exe
endif

.PHONY : build-server

build-server :
	go build -o ./bin/$(binName)$(suffix) ./cmd/server