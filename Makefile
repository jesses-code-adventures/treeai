# OpenCode Trees Makefile

BINARY_NAME=opentree
BUILD_DIR=build
INSTALL_DIR=$(if $(XDG_BIN_HOME),$(XDG_BIN_HOME),$(HOME)/.local/bin)

.PHONY: build install test clean help

# Default target
all: build

# Build the binary
build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

# Install the binary to XDG_BIN_HOME or ~/.local/bin
install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)

# Format code
fmt:
	go fmt ./...

# Lint code
vet:
	go vet ./...

# Run all checks (format, vet, test)
check: fmt vet test

# Show help
help:
	@echo "Available targets:"
	@echo "  build   - Build the binary"
	@echo "  install - Build and install to XDG_BIN_HOME or ~/.local/bin"
	@echo "  test    - Run tests"
	@echo "  clean   - Remove build artifacts"
	@echo "  fmt     - Format code"
	@echo "  vet     - Run go vet"
	@echo "  check   - Run fmt, vet, and test"
	@echo "  help    - Show this help"
