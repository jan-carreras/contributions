#!/usr/bin/env make

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

all: test build
test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN)
build:
	$(GOBUILD) -o bin/ghc cmd/ghc/ghc.go

