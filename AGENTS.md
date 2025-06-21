# Agent Guidelines for OpenCode Trees

## Build/Test Commands
- `go build ./...` - Build all packages
- `go test ./...` - Run all tests
- `go test ./path/to/package` - Run tests for specific package
- `go test -run TestFunctionName` - Run single test function
- `go mod tidy` - Clean up dependencies
- `go fmt ./...` - Format all Go files
- `go vet ./...` - Run static analysis

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