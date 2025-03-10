.PHONY: build run test clean examples

# Build settings
BINARY_NAME=flashblock
BUILD_DIR=./build
MAIN_FILE=./cmd/server/main.go
CLIENT_FILE=./cmd/client/main.go
EXAMPLE_DIR=./examples

# Get Go version from go.mod
GO_VERSION=$(shell grep -E "^go [0-9]+\.[0-9]+(\.[0-9]+)?" go.mod | cut -d " " -f 2)

build:
	@echo "Building ${BINARY_NAME}..."
	@mkdir -p ${BUILD_DIR}
	go build -o ${BUILD_DIR}/${BINARY_NAME} ${MAIN_FILE}
	go build -o ${BUILD_DIR}/client ${CLIENT_FILE}
	@echo "Build complete: ${BUILD_DIR}/${BINARY_NAME}"

run: build
	@echo "Running ${BINARY_NAME}..."
	${BUILD_DIR}/${BINARY_NAME}

run-custom: build
	@echo "Running ${BINARY_NAME} with custom flags..."
	${BUILD_DIR}/${BINARY_NAME} --rpc-addr=:8888 --api-addr=:8889 --block-interval=500ms

run-client:
	@echo "Running client..."
	${BUILD_DIR}/client

test:
	@echo "Running tests..."
	go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -cover ./...

bench:
	@echo "Running benchmarks..."
	go test -v -bench=. ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...

lint:
	@echo "Linting code..."
	go vet ./...
	@if command -v golint > /dev/null; then \
		golint ./...; \
	else \
		echo "golint not installed. Installing..."; \
		go install golang.org/x/lint/golint@latest; \
		golint ./...; \
	fi

clean:
	@echo "Cleaning..."
	rm -rf ${BUILD_DIR}
	go clean

build-examples:
	@echo "Building examples..."
	go build -o ${BUILD_DIR}/client ${EXAMPLE_DIR}/client.go
	go build -o ${BUILD_DIR}/ws_client ${EXAMPLE_DIR}/ws_client/main.go


# Default target
all: build 