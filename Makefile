.PHONY: all clean build test install cross-build

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
	@mkdir -p $(BIN_PATH)
	@cd $(BIN_PATH)
	@go build -v $(LD_FLAGS) -o $(BIN_PATH)/$(CMD) $(BUILD_PATH)

test:
	@echo "=> Executing unit tests"
	@go test ./...

install:
	@echo "=> Installing ${VERSION} in go/bin"
	@go install $(LD_FLAGS) $(BUILD_PATH)

cross-build: clean
	@for OS in darwin linux windows; do														\
		for ARCH in 386 amd64; do															\
			echo "=> Building $$OS-$$ARCH";													\
			env GOOS=$$OS GOARCH=$$ARCH														\
			go build $(LD_FLAGS) -o $(BIN_PATH)/targets/$$OS/$$ARCH/$(CMD) $(BUILD_PATH);	\
		done																				\
	done
