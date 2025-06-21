# Agent Guidelines for treeai

## Build/Test Commands
- `make build` - Build the treeai binary
- `make install` - Build and install to XDG_BIN_HOME or ~/.local/bin
- `make test` - Run all tests
- `make check` - Run fmt, vet, and test
- `make clean` - Remove build artifacts
- `go test -run TestFunctionName` - Run single test function
- `go mod tidy` - Clean up dependencies

## Code Style Guidelines
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
- In the makefile, do not echo what is happening - instead, just allow the command to print
- ensure the binary is always built into the `build/` directory
- **NEVER use emojis** in code, comments, documentation, or any other files

## Project Structure
- This is a CLI tool for managing AI development worktrees
- Git module: `github.com/jesses-code-adventures/git` (git cli wrapper)
- Tmux module: `github.com/jesses-code-adventures/tmux` (tmux cli wrapper)
- Opentree module: `github.com/jesses-code-adventures/treeai` (all application specific logic)
- Cmd module: `github.com/jesses-code-adventures/treeai/cmd` (command line interface)

## Testing
- Use `go test` for unit tests
- Use `go test -race` for race tests

## CI/CD
- Use GitHub Actions for CI/CD
