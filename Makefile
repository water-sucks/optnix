APP_NAME := optnix
BUILD_VAR_PKG := github.com/water-sucks/optnix/internal/build

VERSION ?= $(shell git describe --tags --always)

LDFLAGS := -X $(BUILD_VAR_PKG).Version=$(VERSION)

# Disable CGO by default. This should be a static executable.
CGO_ENABLED ?= 0

all: build

.PHONY: build
build:
	@echo "building $(APP_NAME)..."
	CGO_ENABLED=$(CGO_ENABLED) go build -o ./$(APP_NAME) -ldflags="$(LDFLAGS)" .

.PHONY: clean
clean:
	@echo "cleaning up..."
	go clean

.PHONY: test
test:
	@echo "running tests..."
	CGO_ENABLED=$(CGO_ENABLED) go test ./...

.PHONY: site
site:
	# -d is interpreted relative to the book directory.
	mdbook build ./doc -d ../site

.PHONY: serve-site
serve-site:
	mdbook serve ./doc --open
