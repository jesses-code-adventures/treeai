# treeai

A cli application & tmux plugin for tight integration between tmux, git worktrees and [Opencode](https://github.com/sst/OpenCode).

## Installation

```bash
curl -sSL https://raw.githubusercontent.com/jesses-code-adventures/treeai/main/install.sh | bash
```

Or manually:
```bash
git clone https://github.com/jesses-code-adventures/treeai.git
cd treeai && make install
```

Add to `~/.tmux.conf`:
```tmux
bind-key o command-prompt -p "worktree name:" "run-shell 'treeai %%'"
```

## Usage

- `treeai branch-name` - Create worktree + tmux session with opencode
- `treeai branch-name --merge` - Merge worktree and cleanup
- `--silent` - Suppress output
- `--window "cmd"` - Add tmux windows with custom commands

### Development Commands

- `make build` - Build the treeai binary
- `make install` - Build and install to XDG_BIN_HOME or ~/.local/bin
- `make test` - Run all tests
- `make check` - Run fmt, vet, and test
- `make clean` - Remove build artifacts
- `go test -run TestFunctionName` - Run single test function
- `go mod tidy` - Clean up dependencies


### Submitting Changes

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes following the code style guidelines
4. Run tests: `make check`
5. Push to your fork and submit a pull request

## License

MIT License - see [LICENSE](LICENSE) for details.
