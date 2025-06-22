# treeai

A cli application & tmux plugin for tight integration between tmux, git worktrees and [Opencode](https://github.com/sst/OpenCode).

## Installation

```bash
# Clone the repository
git clone https://github.com/jesses-code-adventures/treeai.git
cd treeai
make install
```

### ~/.tmux.conf

After installing `treeai`, add the following to your `~/.tmux.conf` file, then reload your configuration:

```tmux
# treeai
bind-key o command-prompt -p "worktree name:" "run-shell 'treeai %%'" # binds creation of a new worktree to `<prefix>o`
```

## Usage

### Create a worktree & tmux session

Call `treeai branch-name` to create a new branch & worktree called `branch-name` in the `.treeai` directory of your project. A tmux session will then be created and switched to, with `opencode` running in the default window (window 0). By default, you can toggle between your main tmux session and the `opencode` session using `<prefix>L` to alternate between two tmux sessions. This allows you to assign opencode some work and quickly switch back to what you were doing.

### Merge your worktree and clean up git environment

When you're satisfied with what `opencode` has implemented, merge the worktree from your `main` directory by calling `treeai branch-name --merge`. This will check out the worktree, rebase it against `main`, merge it into `main`, prune worktrees and delete merged git branches.

### Silent mode

Use the `--silent` flag to suppress all output from treeai operations:

```bash
# Create worktree silently
treeai branch-name --silent

# Merge worktree silently
treeai branch-name --merge --silent
```

This is useful for scripting or when you want to minimize output during automated workflows.

### Custom window layouts

Use the `--window` flag to create additional tmux windows with custom commands:

```bash
# Create a worktree with additional windows
treeai branch-name --window "npm run dev" --window "git log --oneline"

# Multiple windows with different commands
treeai branch-name --window "make build" --window "make test" --window "htop"
```

Each `--window` flag creates a new tmux window that runs the specified bash command. The `opencode` window (window 0) remains the default focused window, and you can navigate between windows using standard tmux commands (`<prefix>0`, `<prefix>1`, etc.).

### Development Commands

- `make build` - Build the treeai binary
- `make install` - Build and install to XDG_BIN_HOME or ~/.local/bin
- `make test` - Run all tests
- `make check` - Run fmt, vet, and test
- `make clean` - Remove build artifacts
- `go test -run TestFunctionName` - Run single test function
- `go mod tidy` - Clean up dependencies

### Code Style Guidelines

- Use `gofmt` for consistent formatting
- Follow standard Go naming conventions (PascalCase for exported, camelCase for unexported)
- Use meaningful variable names, avoid abbreviations
- Keep functions small and focused
- Use early returns to reduce nesting
- Handle errors explicitly, don't ignore them
- Use `context.Context` for cancellation and timeouts
- Group imports: standard library, third-party, local packages
- Add comments for exported functions and types
- Use interfaces for testability and flexibility

### Project Structure

- **Git module**: `github.com/jesses-code-adventures/git` (git cli wrapper)
- **Tmux module**: `github.com/jesses-code-adventures/tmux` (tmux cli wrapper)
- **Opentree module**: `github.com/jesses-code-adventures/treeai` (application logic)
- **Cmd module**: `github.com/jesses-code-adventures/treeai/cmd` (command line interface)

### Testing

Run tests before submitting changes:

```bash
make test
```

For race condition testing:
```bash
go test -race ./...
```

### Submitting Changes

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes following the code style guidelines
4. Run tests: `make check`
5. Commit your changes with clear, descriptive messages
6. Push to your fork and submit a pull request

### Reporting Issues

Please use GitHub Issues to report bugs or request features. Include:
- Go version
- Operating system
- Steps to reproduce
- Expected vs actual behavior

## License

MIT License - see [LICENSE](LICENSE) for details.
