#!/bin/bash

set -e

REPO_URL="https://github.com/jesses-code-adventures/treeai.git"
BINARY_NAME="treeai"
INSTALL_DIR="${XDG_BIN_HOME:-$HOME/.local/bin}"

echo "Installing treeai..."

if ! command -v go &> /dev/null; then
    echo "Error: Go is required but not installed"
    exit 1
fi

TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

cd "$TEMP_DIR"
git clone "$REPO_URL" .

make build

mkdir -p "$INSTALL_DIR"
cp "build/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"

echo "treeai installed to $INSTALL_DIR/$BINARY_NAME"
echo ""
echo "Add the following to your ~/.tmux.conf:"
echo 'bind-key o command-prompt -p "worktree name:" "run-shell '\''treeai %%'\''"'
echo ""
echo "Make sure $INSTALL_DIR is in your PATH"