#!/usr/bin/env make

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

.PHONY: test
test:
	$(GOTEST) -v ./...

.PHONY: clean
clean:
	@$(GOCLEAN)
	@rm -f bin/*

.PHONY: build
build: clean bin/ghc

.PHONY: all
all: test build

bin/ghc:
	$(GOBUILD) -ldflags "-s -w" -o bin/ghc cmd/ghc/ghc.go


