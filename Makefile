GO := go
GO_BUILD_FLAGS =
GO_TEST_FLAGS =
GO_BUILD = CGO_ENABLED=0 $(GO) build $(GO_BUILD_FLAGS)
GO_TEST = $(GO) test $(GO_TEST_FLAGS) -timeout=$(TEST_TIMEOUT)
INSTALL = $(QUIET)install
BINDIR ?= /usr/local/bin
SUBDIRS_TRIANGLE_CLI := .
TARGET=triangle
VERSION=$(shell cat VERSION)
# homebrew uses the github release's tarball of the source that does not contain the '.git' directory.
GIT_BRANCH = $(shell command -v git >/dev/null 2>&1 && git rev-parse --abbrev-ref HEAD 2> /dev/null)
GIT_HASH = $(shell command -v git >/dev/null 2>&1 && git rev-parse --short HEAD 2> /dev/null)
GO_TAGS ?=
IMAGE_REPOSITORY ?= quay.io/khulnasoft/triangle
IMAGE_TAG ?= $(if $(findstring -dev,$(VERSION)),latest,v$(VERSION))
CONTAINER_ENGINE ?= docker
RELEASE_UID ?= $(shell id -u)
RELEASE_GID ?= $(shell id -g)

TEST_TIMEOUT ?= 5s

# renovate: datasource=docker depName=golangci/golangci-lint
GOLANGCILINT_WANT_VERSION = v1.55.1
GOLANGCILINT_IMAGE_SHA = sha256:c4e67eb904109ade78e2f38d98a424502f016db5676409390469bcdafea0f57d
GOLANGCILINT_VERSION = $(shell golangci-lint version 2>/dev/null)

# renovate: datasource=docker depName=library/golang
GOLANG_IMAGE_VERSION = 1.21.3-alpine3.18
GOLANG_IMAGE_SHA = sha256:96a8a701943e7eabd81ebd0963540ad660e29c3b2dc7fb9d7e06af34409e9ba6

# Add the ability to override variables
-include Makefile.override

all: triangle

triangle:
	$(MAKE) -C $(SUBDIRS_TRIANGLE_CLI) triangle-bin

triangle-bin:
	$(GO_BUILD) $(if $(GO_TAGS),-tags $(GO_TAGS)) -ldflags "-w -s -X 'github.com/khulnasoft/triangle/pkg.GitBranch=${GIT_BRANCH}' -X 'github.com/khulnasoft/triangle/pkg.GitHash=$(GIT_HASH)' -X 'github.com/khulnasoft/triangle/pkg.Version=${VERSION}'" -o $(TARGET) $(SUBDIRS_TRIANGLE_CLI)

release:
	$(CONTAINER_ENGINE) run --rm --workdir /triangle --volume `pwd`:/triangle docker.io/library/golang:$(GOLANG_IMAGE_VERSION)@$(GOLANG_IMAGE_SHA) \
		sh -c "apk add --no-cache setpriv make git && \
			/usr/bin/setpriv --reuid=$(RELEASE_UID) --regid=$(RELEASE_GID) --clear-groups make GOCACHE=/tmp/gocache local-release"

local-release: clean
	set -o errexit; \
	for OS in darwin linux windows; do \
		EXT=; \
		ARCHS=; \
		case $$OS in \
			darwin) \
				ARCHS='amd64 arm64'; \
				;; \
			linux) \
				ARCHS='386 amd64 arm arm64'; \
				;; \
			windows) \
				ARCHS='386 amd64 arm64'; \
				EXT=".exe"; \
				;; \
		esac; \
		for ARCH in $$ARCHS; do \
			echo Building release binary for $$OS/$$ARCH...; \
			test -d release/$$OS/$$ARCH|| mkdir -p release/$$OS/$$ARCH; \
			env GOOS=$$OS GOARCH=$$ARCH $(GO_BUILD) $(if $(GO_TAGS),-tags $(GO_TAGS)) -ldflags "-w -s -X 'github.com/khulnasoft/triangle/pkg.Version=${VERSION}'" -o release/$$OS/$$ARCH/$(TARGET)$$EXT; \
			tar -czf release/$(TARGET)-$$OS-$$ARCH.tar.gz -C release/$$OS/$$ARCH $(TARGET)$$EXT; \
			(cd release && sha256sum $(TARGET)-$$OS-$$ARCH.tar.gz > $(TARGET)-$$OS-$$ARCH.tar.gz.sha256sum); \
		done; \
		rm -r release/$$OS; \
	done;

install: triangle
	$(INSTALL) -m 0755 -d $(DESTDIR)$(BINDIR)
	$(INSTALL) -m 0755 ./triangle $(DESTDIR)$(BINDIR)

clean:
	rm -f $(TARGET)
	rm -rf ./release

test:
	$(GO_TEST) -race -cover $$($(GO) list ./...)

bench: TEST_TIMEOUT=30s
bench:
	$(GO_TEST) -bench=. $$($(GO) list ./...)

ifneq (,$(findstring $(GOLANGCILINT_WANT_VERSION:v%=%),$(GOLANGCILINT_VERSION)))
check:
	golangci-lint run
else
check:
	$(CONTAINER_ENGINE) run --rm -v `pwd`:/app -w /app docker.io/golangci/golangci-lint:$(GOLANGCILINT_WANT_VERSION)@$(GOLANGCILINT_IMAGE_SHA) golangci-lint run
endif

image:
	$(CONTAINER_ENGINE) build $(DOCKER_FLAGS) -t $(IMAGE_REPOSITORY)$(if $(IMAGE_TAG),:$(IMAGE_TAG)) .

.PHONY: all triangle release install clean test bench check image
