VERSION = $(shell git describe --tags)
VER = $(shell git describe --tags --abbrev=0)
DATE = $(shell date -u '+%Y-%m-%d_%H:%M:%S%Z')

FLAG_MODULE = GO111MODULE=on
FLAGS_SHARED = $(FLAG_MODULE) GOARCH=amd64
NO_C = CGO_ENABLED=0
FLAGS_LINUX = $(FLAGS_SHARED) GOOS=linux
FLAGS_MAC = $(FLAGS_SHARED) GOOS=darwin
FLAGS_MAC_ARM = GO111MODULE=on $GOARCH=arm64 GOOS=darwin
FLAGS_WIN = $(FLAGS_SHARED) GOOS=windows
FLAGS_LD=-ldflags "-s -w -X github.com/gnames/gnparser.Build=${DATE} \
                  -X github.com/gnames/gnparser.Version=${VERSION}"
GOCMD = go
GOBUILD = $(GOCMD) build $(FLAGS_LD)
GOINSTALL = $(GOCMD) install $(FLAGS_LD)
GOCLEAN = $(GOCMD) clean
GOGET = $(GOCMD) get

RELEASE_DIR ?= "/tmp"
BUILD_DIR ?= "."
CLIB_DIR ?= "."

test: deps install
	$(FLAG_MODULE) go test -race ./...

deps:
	$(GOCMD) mod download;

tools: deps
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

peg:
	cd ent/parser; \
	peg query.peg; \
	goimports -w query.peg.go;

