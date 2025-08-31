APP_NAME := optnix
BUILD_VAR_PKG := github.com/water-sucks/optnix/internal/build

VERSION ?= $(shell (git describe --tags --always || echo unknown))
GIT_REVISION ?= $(shell (git rev-parse HEAD || echo main))

LDFLAGS := -X $(BUILD_VAR_PKG).Version=$(VERSION)

GENERATED_MODULE_DOCS := doc/src/usage/generated-module.md
NIX_MODULE := nix/modules/nixos.nix
GITHUB_URL := https://github.com/water-sucks/optnix/blob/$(GIT_REVISION)/$(NIX_MODULE)

MANPAGES_SRC := $(wildcard doc/man/*.scd)
MANPAGES := $(patsubst doc/man/%.scd,%,$(MANPAGES_SRC))

# Disable CGO by default. This should be a static executable.
CGO_ENABLED ?= 0

all: build $(MANPAGES)

.PHONY: build
build:
	@echo "building $(APP_NAME)..."
	CGO_ENABLED=$(CGO_ENABLED) go build -o ./$(APP_NAME) -ldflags="$(LDFLAGS)" .

.PHONY: clean
clean:
	@echo "cleaning up..."
	go clean
	rm -rf ./nixos site/ $(MANPAGES)

.PHONY: check
check:
	@echo "running checks..."
	golangci-lint run

.PHONY: test
test:
	@echo "running tests..."
	CGO_ENABLED=$(CGO_ENABLED) go test ./...

site: $(GENERATED_MODULE_DOCS)
	# -d is interpreted relative to the book directory.
	mdbook build ./doc -d ../site

.PHONY: serve-site
serve-site: $(GENERATED_MODULE_DOCS)
	mdbook serve ./doc --open

$(GENERATED_MODULE_DOCS): $(NIX_MODULE)
	nix-options-doc -f markdown -p $< --strip-prefix | \
		tail -n +4 | \
		sed -E 's|\(#L([0-9]+)\)|('"$(GITHUB_URL)"'#L\1)|g' \
		> $@

man: $(MANPAGES)

%: doc/man/%.scd
	scdoc < $< > $@
