.PHONY: all clean build test test-coverage test-e2e install build-all

all: build test-e2e test

BIN_PATH = bin
VERSION = $(shell git describe --tags)
CMD = license-header-checker
BUILD_PATH = ./cmd/$(CMD)
LD_FLAGS = -ldflags="-X 'main.version=${VERSION}'"

clean:
	@echo ">> Cleaning project"
	rm -rf ${BIN_PATH}

build:
	@echo ">> Building project"
	go mod tidy
	go build -v ${LD_FLAGS} -o ${BIN_PATH}/${CMD} ${BUILD_PATH}

test:
	@echo ">> Running unit tests"
	go test -cover ./...

test-coverage:
	@echo ">> Running unit tests and checking coverage"
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

test-e2e: build
	@echo ">> Executing end to end tests"
	cd test && bash test.sh ../${BIN_PATH}/${CMD}

install:
	@echo ">> Installing ${CMD} ${VERSION} in ${GOPATH}/bin"
	go install ${LD_FLAGS} ${BUILD_PATH}

TOOLS := $(CURDIR)/.tools

install-tools:
	mkdir -p ${TOOLS}
	GOPATH=${TOOLS} go install github.com/caarlos0/svu@latest
	GOPATH=${TOOLS} go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	GOPATH=${TOOLS} go install github.com/goreleaser/goreleaser@latest

lint:
	${TOOLS}/bin/golangci-lint run

build-one:
	${TOOLS}/bin/goreleaser build --single-target

# Generate the binaries for all the build targets
build-all: install-tools
	${TOOLS}/bin/goreleaser release --snapshot --rm-dist

# Verify that the project builds without errors for all build targets
build-ci:
	goreleaser build


ARG ?= next
VERSION := $(shell ${TOOLS}/bin/svu $(ARG))
release: install-tools
	echo "$(VERSION)"
	#git tag "$(VERSION)"
	#git push origin "$(VERSION)"

release-major: ARG=major
release-major: release

release-minor: ARG=minor
release-minor: release

release-patch: ARG=patch
release-patch: release
