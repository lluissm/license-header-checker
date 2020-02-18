all: build test

clean:
	@echo "=> Cleaning..."
	@rm -rf bin
	@go clean -testcache -cache

build:
	@echo "=> Building..."
	@go mod tidy
	@mkdir -p bin
	@cd bin
	@go build -v -o bin/license-header-checker ./cmd/license-header-checker

test:
	@echo "=> Testing..."
	@go test ./...

install:
	@echo "=> Installing..."
	@go install ./cmd/license-header-checker
