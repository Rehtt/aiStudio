cur_makefile_path := $(abspath $(lastword ./))
binName = $(shell echo $(cur_makefile_path)|awk -F '/' '{ print $$NF }')

ifeq ($(shell uname),Windows_NT)
	suffix = .exe
endif

.PHONY : build-server

build-server :
	go build -o ./bin/server$(suffix) ./cmd/server

build-discord :
	go build -o ./bin/discord$(suffix) ./cmd/discord

update-kit:
	go get -u github.com/Rehtt/Kit@master
	go get -u github.com/Rehtt/Kit/web@master