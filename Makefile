.PHONY: build run-server run-client run-analyze fmt lint clean

# Build settings
BINARY_NAME=flashblock
BUILD_DIR=./bin
MAIN_FILE=./cmd/server/main.go
CLIENT_FILE=./cmd/client/main.go
# Get Go version from go.mod
GO_VERSION=$(shell grep -E "^go [0-9]+\.[0-9]+(\.[0-9]+)?" go.mod | cut -d " " -f 2)

build:
	@echo "Building ${BINARY_NAME}..."
	@mkdir -p ${BUILD_DIR}
	go build -o ${BUILD_DIR}/${BINARY_NAME} ${MAIN_FILE}
	go build -o ${BUILD_DIR}/client ${CLIENT_FILE}
	@echo "Build complete: ${BUILD_DIR}/${BINARY_NAME}"

run-server:
	@echo "Running ${BINARY_NAME}..."
	sudo ${BUILD_DIR}/${BINARY_NAME}

run-client:
	@echo "Running client..."
	${BUILD_DIR}/client

run-analyze:
	@echo "Running analyze..."
	./logs/analyze.sh

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

# Default target
all: build 