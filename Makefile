# Git Tools Makefile

# Binary name
BINARY_NAME=git-tools

# Build the application
build:
	go build -o $(BINARY_NAME) .

# Build and run
run: build
	./$(BINARY_NAME)

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)

# Install to GOPATH/bin
install:
	go install .

# Build for multiple platforms
build-all:
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 .
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe .

# Test the application
test:
	go test ./...

# Format the code
fmt:
	go fmt ./...

# Get dependencies
deps:
	go mod tidy

# Help
help:
	@echo "Available targets:"
	@echo "  build     - Build the application"
	@echo "  run       - Build and run the application"
	@echo "  clean     - Clean build artifacts"
	@echo "  install   - Install to GOPATH/bin"
	@echo "  build-all - Build for multiple platforms"
	@echo "  test      - Run tests"
	@echo "  fmt       - Format code"
	@echo "  deps      - Get dependencies"
	@echo "  help      - Show this help message"

.PHONY: build run clean install build-all test fmt deps help 