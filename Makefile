.PHONY: all clean build test install cross-build

all: build test

CMD = license-header-checker

clean:
	@echo "=> Cleaning project"
	@rm -rf bin
	@go clean -testcache -cache

build:
	@echo "=> Building project"
	@go mod tidy
	@mkdir -p bin
	@cd bin
	@go build -v -o bin/$(CMD) ./cmd/$(CMD)

test:
	@echo "=> Executing unit tests"
	@go test ./...

install:
	@echo "=> Installing in go/bin"
	@go install ./cmd/$(CMD)

cross-build: clean
	@for OS in darwin linux windows; do									\
		for ARCH in 386 amd64; do										\
			echo "=> Building $$OS-$$ARCH";								\
			env GOOS=$$OS GOARCH=$$ARCH									\
			go build -o bin/targets/$$OS/$$ARCH/$(CMD) ./cmd/$(CMD);    \
		done															\
	done
