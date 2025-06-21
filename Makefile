# OpenCode Trees Makefile

BINARY_NAME=opentree
BUILD_DIR=build
INSTALL_DIR=$(if $(XDG_BIN_HOME),$(XDG_BIN_HOME),$(HOME)/.local/bin)

.PHONY: build install test clean help

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ Built $(BINARY_NAME)"

# Install the binary to XDG_BIN_HOME or ~/.local/bin
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✅ Installed $(BINARY_NAME) to $(INSTALL_DIR)"
	@echo "💡 Make sure $(INSTALL_DIR) is in your PATH"

# Run tests
test:
	@echo "Running tests..."
	go test ./...
	@echo "✅ Tests completed"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	@echo "✅ Cleaned"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "✅ Code formatted"

# Lint code
vet:
	@echo "Running go vet..."
	go vet ./...
	@echo "✅ Vet completed"

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
