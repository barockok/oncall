SHELL := /usr/bin/bash

APP := engine-go
PKG := github.com/barockok/oncall/engine-go
CMD := $(PKG)/cmd/server
GO ?= go
PORT ?= 8080

.PHONY: all build run tidy lint test clean

all: build

build:
	$(GO) build -ldflags "-X $(PKG)/internal/version.version=$$(git describe --tags --always 2>/dev/null || echo dev)" ./cmd/server

run:
	PORT=$(PORT) $(GO) run ./cmd/server

tidy:
	$(GO) mod tidy

clean:
	rm -f server