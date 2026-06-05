# flowState-cli local development tasks.
#
# Local builds go into ./dist (gitignored) — NEVER into the repo root — so a
# stale binary can't shadow the npm-installed `flowstate` on your PATH.
# Use `flowstate --version` to confirm which binary actually runs.

BINARY      := flowstate
DIST_DIR    := dist
VERSION     := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT      := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
LDFLAGS     := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)

.PHONY: build run test cover clean install-local

## build: compile the binary into ./dist
build:
	@mkdir -p $(DIST_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY) ./cmd/flowState
	@echo "built $(DIST_DIR)/$(BINARY) ($(VERSION), $(COMMIT))"

## run: build then run from ./dist
run: build
	./$(DIST_DIR)/$(BINARY)

## test: run the full test suite
test:
	go test ./...

## cover: run tests with coverage summary
cover:
	go test -cover ./...

## clean: remove build output
clean:
	rm -rf $(DIST_DIR)

## install-local: install the dev build onto your PATH via `go install`
install-local:
	go install -ldflags "$(LDFLAGS)" ./cmd/flowState
