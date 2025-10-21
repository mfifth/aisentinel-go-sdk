# Contributing to AISentinel Go SDK

Thank you for your interest in contributing to the AISentinel Go SDK! We welcome contributions from the community.

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional, for convenience scripts)

### Setup

1. Fork and clone the repository:
```bash
git clone https://github.com/your-username/aisentinel-go-sdk.git
cd aisentinel-go-sdk
```

2. Install dependencies:
```bash
go mod download
```

3. Run tests to ensure everything works:
```bash
go test ./...
```

## Development Workflow

### 1. Choose an Issue

- Check the [issue tracker](https://github.com/aisentinel/aisentinel-go-sdk/issues) for open issues
- Look for issues labeled `good first issue` or `help wanted`
- Comment on the issue to indicate you're working on it

### 2. Create a Branch

Create a feature branch from `main`:

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/issue-number-description
```

### 3. Make Changes

- Write clear, concise commit messages
- Follow Go conventions and best practices
- Add tests for new functionality
- Update documentation as needed

### 4. Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test ./pkg/client

# Run benchmarks
go test -bench=. ./...
```

### 5. Code Quality

- Run `go fmt` to format your code
- Run `go vet` to check for common errors
- Consider using `golint` or `golangci-lint` for additional checks

### 6. Documentation

- Update README.md for API changes
- Add GoDoc comments for exported functions/types
- Update examples if needed

### 7. Commit and Push

```bash
git add .
git commit -m "feat: add new feature description"
git push origin feature/your-feature-name
```

### 8. Create Pull Request

- Go to your fork on GitHub
- Click "New Pull Request"
- Fill out the pull request template
- Reference any related issues

## Code Guidelines

### Go Conventions

- Follow standard Go formatting (`go fmt`)
- Use `gofmt -s` for additional simplifications
- Write clear, idiomatic Go code
- Use meaningful variable and function names

### Error Handling

- Return errors rather than panicking
- Use error wrapping for context: `fmt.Errorf("failed to connect: %w", err)`
- Define custom error types when appropriate

### Testing

- Write unit tests for all exported functions
- Use table-driven tests for multiple test cases
- Mock external dependencies
- Aim for good test coverage (>80%)

### Documentation

- Add GoDoc comments for all exported identifiers
- Keep comments up to date with code changes
- Provide usage examples in comments

## Pull Request Process

1. **Title**: Use conventional commit format (e.g., `feat: add batch evaluation`, `fix: handle timeout errors`)

2. **Description**: Include:
   - What changes were made
   - Why the changes were needed
   - How to test the changes
   - Screenshots/videos for UI changes (if applicable)

3. **Checklist**:
   - [ ] Tests pass (`go test ./...`)
   - [ ] Code is formatted (`go fmt`)
   - [ ] No linting errors (`go vet`)
   - [ ] Documentation updated
   - [ ] Examples work correctly

4. **Review Process**:
   - Maintainers will review your PR
   - Address any feedback or requested changes
   - Once approved, your PR will be merged

## Issue Reporting

When reporting bugs or requesting features:

1. **Check existing issues** to avoid duplicates
2. **Use issue templates** when available
3. **Provide detailed information**:
   - Go version (`go version`)
   - OS and architecture
   - Steps to reproduce
   - Expected vs actual behavior
   - Code samples or test cases

## Community

- Join our [Community Forum](https://community.aisentinel.com) for discussions
- Follow [@AISentinel](https://twitter.com/AISentinel) for updates
- Read our [Blog](https://blog.aisentinel.com) for technical insights

## License

By contributing to this project, you agree that your contributions will be licensed under the Apache License 2.0.

Thank you for contributing to AISentinel Go SDK! 🚀
