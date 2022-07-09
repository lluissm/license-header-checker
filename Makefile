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

build-all: clean
	@echo ">> Building MacOS intel"
	env GOOS=darwin GOARCH=amd64 go build ${LD_FLAGS} -o ${BIN_PATH}/targets/${CMD}_macos_intel $(BUILD_PATH)

	@echo ">> Building MacOS arm"
	env GOOS=darwin GOARCH=arm64 go build ${LD_FLAGS} -o ${BIN_PATH}/targets/${CMD}_macos_arm64 $(BUILD_PATH)

	@echo ">> Building Linux amd64"
	env GOOS=linux GOARCH=amd64 go build ${LD_FLAGS} -o ${BIN_PATH}/targets/${CMD}_linux_amd64 $(BUILD_PATH)

	@echo ">> Building Windows amd64"
	env GOOS=windows GOARCH=amd64 go build ${LD_FLAGS} -o ${BIN_PATH}/targets/${CMD}_windows_amd64.exe $(BUILD_PATH)
