#!/usr/bin/env make

GOCMD=go

.PHONY: test
test:
	$(GOCMD) test -v ./...

.PHONY: clean
clean:
	@$(GOCMD) clean
	@rm -f bin/*

.PHONY: build
build: clean bin/ghc

.PHONY: install
install:
	$(GOCMD) install cmd/ghc/ghc.go

.PHONY: all
all: test build

bin/ghc:
	$(GOCMD) build -ldflags "-s -w" -o bin/ghc cmd/ghc/ghc.go


