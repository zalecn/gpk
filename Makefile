
# computes the current branch name for snapshot deploy
BRANCH        := $(shell git rev-parse --abbrev-ref HEAD)
ifeq ($(BRANCH),HEAD)
	BRANCH:=master
endif
VERSION :=$(BRANCH) #default beahviour
#computes the current packages names (information lies in the .gpk file)
PKG           := $(shell gpk ? -n)

.PHONY: all init generate compile deploy clean

all: compile

compile:;	
	GOPATH=`pwd` go install -ldflags '-X ericaro.net/gopack/cmds.GopackageVersion $(VERSION)' ./src/...

install: compile
	sudo cp ./bin/gpk /usr/bin/gpk

clean:;	rm -Rf bin/*
 