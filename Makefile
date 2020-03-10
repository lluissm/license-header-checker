.PHONY: all clean build test test-e2e install cross-build cross-pack

all: build test

BIN_PATH = bin
VERSION = $(shell git describe --tags)
CMD = license-header-checker
BUILD_PATH = ./cmd/$(CMD)
LD_FLAGS = -ldflags="-X 'main.version=${VERSION}'"

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
	@go test ./...

test-cover:
	@echo "=> Running unit tests and checking coverage"
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

test-e2e: build
	@echo "=> Executing end to end tests"
	@cd test && bash test.sh ../$(BIN_PATH)/$(CMD)

install:
	@echo "=> Installing ${CMD} ${VERSION} in go/bin"
	@go install $(LD_FLAGS) $(BUILD_PATH)

cross-build: clean
	@for OS in darwin linux windows; do								\
		for ARCH in 386 amd64; do									\
			fos=$$OS;												\
			if [ $$OS = "darwin" ]; then							\
				fos=mac;											\
			fi;														\
			farch=32bit;											\
			if [ $$ARCH = "amd64" ]; then							\
				farch=64bit;										\
			fi;														\
			fname=$(BIN_PATH)/targets/$(CMD)-$$fos/$$farch/$(CMD);	\
			if [ $$OS = "windows" ]; then							\
				fname=$$fname.exe;									\
			fi;														\
			echo "=> Building $$fos $$farch";						\
			env GOOS=$$OS GOARCH=$$ARCH								\
			go build $(LD_FLAGS) -o $$fname $(BUILD_PATH);			\
		done														\
	done

cross-pack: cross-build
	@echo "=> Zipping folders"
	@cd $(BIN_PATH)/targets/ && for OS in mac linux windows; do		\
		zip -q -r ../$(CMD)-$$OS.zip $(CMD)-$$OS -x "*/\*";			\
	done
	@rm -rf $(BIN_PATH)/targets
