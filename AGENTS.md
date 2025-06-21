# Agent Guidelines for OpenCode Trees

## Build/Test Commands
- `make build` - Build the opentree binary
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

## Project Structure
- This is a CLI tool for managing AI development worktrees
- Main module: `github.com/jesses-code-adventures/opentree`
- No existing source files yet - follow standard Go project layout when creating