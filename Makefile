# Makefile for go-cfprefs

BASEDIR ?= $(PWD)
SRCDIR ?= $(BASEDIR)
DISTDIR ?= $(BASEDIR)/dist

APPNAME ?= cfprefs
APPVER ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

TIMESTAMP := $(shell date +%s)

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

APPEXE = $(APPNAME)
ifeq ($(GOOS),windows)
	APPEXE = $(APPNAME).exe
endif


.PHONY: all
all: test build


.PHONY: init
init:
	mkdir -p "$(SRCDIR)/tmp"
	cd $(SRCDIR) && go mod download
	cd $(SRCDIR) && go mod tidy


.PHONY: build
build: cli


.PHONY: cli
cli: init
	mkdir -p $(DISTDIR)
	cd $(SRCDIR) && \
	CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-ldflags "-X main.version=$(APPVER) -s -w" \
		-o $(DISTDIR)/$(APPEXE) cli/main.go


.PHONY: unit-test
unit-test: init
	cd $(SRCDIR) && go test -v ./...


.PHONY: test
test: unit-test


.PHONY: test
coverage: test
	cd $(SRCDIR) && go test -v -coverprofile=tmp/coverage.out ./...
	cd $(SRCDIR) && go tool cover -html=tmp/coverage.out -o $(DISTDIR)/coverage.html


.PHONY: static-checks
static-checks: unit-test


.PHONY: preflight
preflight: static-checks


.PHONY: clean
clean:
	cd $(SRCDIR) && go clean


.PHONY: clobber
clobber: clean
	rm -Rf "$(SRCDIR)/dist"
	cd $(SRCDIR) && go clean -modcache
