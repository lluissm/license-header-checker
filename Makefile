.PHONY: all clean build test test-coverage test-e2e install build-all release-major release-minor release-patch

all: build test-e2e test

BIN_PATH = bin
VERSION = $(shell git describe --tags)
CMD = license-header-checker
BUILD_PATH = ./cmd/$(CMD)
LD_FLAGS = -ldflags="-X 'main.version=${VERSION}'"
TOOLS := $(CURDIR)/.tools

# remove bin folder
clean:
	@echo ">> Cleaning project"
	rm -rf ${BIN_PATH}

# Build locally in bin folder
build:
	@echo ">> Building project"
	go mod tidy
	go build -v ${LD_FLAGS} -o ${BIN_PATH}/${CMD} ${BUILD_PATH}

# Run unit tests
test:
	@echo ">> Running unit tests"
	go test -cover ./...

# Run unit tests and launch cover tool to analyze coverage
test-coverage:
	@echo ">> Running unit tests and checking coverage"
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Run tests with real binary modifying real files
test-e2e: build
	@echo ">> Executing end to end tests"
	cd test && bash test.sh ../${BIN_PATH}/${CMD}

# Install locally
install:
	@echo ">> Installing ${CMD} ${VERSION} in ${GOPATH}/bin"
	go install ${LD_FLAGS} ${BUILD_PATH}

# Install necessary tools:
# svu: for semantic versioning
# golangci-lint: go linter
# goreleaser: to automate release
install-tools:
	mkdir -p ${TOOLS}
	GOPATH=${TOOLS} go install github.com/caarlos0/svu@latest
	GOPATH=${TOOLS} go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	GOPATH=${TOOLS} go install github.com/goreleaser/goreleaser@latest

# Execute golangci-lint
lint:
	${TOOLS}/bin/golangci-lint run

# Generate the binaries for all the build targets
# --rm-dist: Remove the dist folder before building
# --snapshot: Generate an unversioned snapshot build, skipping all validations
build-all: install-tools
	${TOOLS}/bin/goreleaser release --snapshot --rm-dist

# Generate tag and push to origin
define release
RELEASE=$(shell ${TOOLS}/bin/svu $1); \
git tag $$RELEASE
git push origin $$RELEASE
endef

# Increment major version, create tag and push to origin
release-major:
	$(call release,major)
# Increment minor version, create tag and push to origin
release-minor:
	$(call release,minor)
# Increment patch version, create tag and push to origin
release-patch:
	$(call release,patch)
