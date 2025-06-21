BINARY_NAME=treeai
WORKTREES_DIR=.opencode-trees
BUILD_DIR=build
INSTALL_DIR=$(if $(XDG_BIN_HOME),$(XDG_BIN_HOME),$(HOME)/.local/bin)

.PHONY: build install test clean help

all: build

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)

test:
	go test ./...

clean:
	rm -rf $(WORKTREES_DIR)
	git worktree prune
	git branch --merged | grep -v "main" | xargs git branch -d
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	git branch

fmt:
	go fmt ./...

vet:
	go vet ./...

check: fmt vet test

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
