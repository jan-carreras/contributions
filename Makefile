#!/usr/bin/env make

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
build:
	$(GOBUILD) -o bin/ghc cmd/ghc/ghc.go
all: test build

