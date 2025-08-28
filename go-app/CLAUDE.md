# PDFminion Development Guide

## Build Commands
- Build: `make mac` (for Apple Silicon) or `make compile-linux-amd64`, `make compile-darwin-amd64`, etc.
- Test: `gotestsum --format=testdox ./...` or `make test`
- Run single test: `go test -v ./internal/domain -run TestFunctionName`
- Lint: `golangci-lint run` or `make lint`
- Coverage: `make cover` (generates coverage.out and coverage.html)

## Code Style
- **Formatting**: Use `gofmt` standards
- **Imports**: Group standard library, third-party, and internal imports with blank lines between
- **Error Handling**: Check all errors; use zerolog for logging (`log.Error().Err(err).Msg("error message")`)
- **Naming**: Use CamelCase for exported names, camelCase for unexported; use descriptive names
- **Tests**: Use testify/assert package for assertions; place tests in `*_test.go` files
- **Comments**: Document all exported functions, types, and constants
- **Constants**: Use `const` blocks for related constants
- **File Structure**: Organize by functionality in the internal/ directory
- **Null Checks**: Always check for nil pointers before dereferencing (e.g., in MergeWith)