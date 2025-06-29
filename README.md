# treeai

A cli application & tmux plugin for tight integration between tmux, git worktrees and [Opencode](https://github.com/sst/OpenCode).

To quickly create isolated environments for coding agents, this cli application automates the following workflow:

- Create a git worktree with a new branch (stored in `$HOME/.local/share/treeai`)
- Open the worktree in a new tmux session or window, with `opencode` open
- Provide the agent with a prompt (or use the --prompt flag to pass one in without focusing the new session/window)
- Use `<leader>L` in tmux to flick back to your main session while the agent works (you might even create some more trees at this point)
- When happy with the agent's work, run `treeai <branch> --merge` to merge the worktree back to main and clean up the tmux sessions/prune worktrees

## Installation

Note that you must have go>=1.24.2 installed - will be fixed in future, when I'm less lazy.

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
- `--command "cmd"` - Add tmux windows with custom commands
- `--window` - Open tmux window instead of session
- `--prompt "prompt"` - Send a prompt to opencode in the new session/window
- `--bin "bin"` - Binary to launch in the tmux session/window, if not `opencode`
- `--data "path"` - Path to data directory, if not `$HOME/.local/share/treeai`
- `--gitignore` - Use .gitignore instead of .git/info/exclude to exclude worktrees from git
- `--debug` - Enable debug logging
- `--copy "file"` - Copy a gitignored file to the worktree

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
