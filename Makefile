.PHONY: all clean build test install cross-build test-e2e

all: build test

BIN_PATH = bin
VERSION = $(shell git describe --tags)
CMD = license-header-checker
BUILD_PATH = ./cmd/$(CMD)
LD_FLAGS = -ldflags="-X 'github.com/lsm-dev/license-header-checker/internal/config.version=${VERSION}'"

clean:
	@echo "=> Cleaning project"
	@rm -rf $(BIN_PATH)
	@go clean -testcache -cache

build:
	@echo "=> Building project"
	@go mod tidy
	@go build -v $(LD_FLAGS) -o $(BIN_PATH)/$(CMD) $(BUILD_PATH)

test:
	@echo "=> Running unit tests"
	@go test -cover ./...

test-e2e: build
	@echo "=> Executing end to end tests"
	@cd test && bash test.sh ../$(BIN_PATH)/$(CMD)

install:
	@echo "=> Installing ${CMD} ${VERSION} in go/bin"
	@go install $(LD_FLAGS) $(BUILD_PATH)

cross-build: clean
	@for OS in darwin linux windows; do														\
		for ARCH in 386 amd64; do															\
			echo "=> Building $$OS-$$ARCH";													\
			env GOOS=$$OS GOARCH=$$ARCH														\
			go build $(LD_FLAGS) -o $(BIN_PATH)/targets/$$OS/$$ARCH/$(CMD) $(BUILD_PATH);	\
		done																				\
	done
