VERSION ?= dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LD_FLAGS = -s -w \
	-X github.com/xpfo-go/akali/internal/command/version.Version=$(VERSION) \
	-X github.com/xpfo-go/akali/internal/command/version.Commit=$(COMMIT) \
	-X github.com/xpfo-go/akali/internal/command/version.BuildTime=$(BUILD_TIME)

BIN_DIR ?= ./dist

.PHONY: build-linux
build-linux:
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LD_FLAGS)" -o $(BIN_DIR)/akali_$(VERSION)_linux_amd64 .

.PHONY: build-darwin
build-darwin:
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="$(LD_FLAGS)" -o $(BIN_DIR)/akali_$(VERSION)_darwin_arm64 .

.PHONY: build-windows
build-windows:
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="$(LD_FLAGS)" -o $(BIN_DIR)/akali_$(VERSION)_win_amd64.exe .

.PHONY: build-all
build-all: build-linux build-darwin build-windows

.PHONY: test
test:
	go test ./...

.PHONY: vet
vet:
	go vet ./...
