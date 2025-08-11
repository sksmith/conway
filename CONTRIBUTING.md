# Contributing to Conway Polyhedron Notation Library

Thank you for your interest in contributing to the Conway Polyhedron Notation Library! This document provides guidelines for contributing to the project.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
3. Set up the development environment
4. Make your changes
5. Test your changes
6. Submit a pull request

## Development Environment

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional, but recommended)

### Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/conway.git
cd conway

# Install dependencies
go mod download

# Install development tools
make dev-deps
```

## Code Style and Standards

### Go Style Guide

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `go fmt` to format your code
- Use meaningful variable and function names
- Add godoc comments for public APIs
- Keep functions focused and small

### Project-Specific Guidelines

- **Package Structure**: Keep the main library code in the `conway` package
- **Error Handling**: Use explicit error handling, avoid panics in library code
- **Thread Safety**: All public APIs should be thread-safe
- **Performance**: Be mindful of memory allocations and algorithmic complexity
- **Testing**: Write tests for new functionality and ensure existing tests pass

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run benchmarks
make bench

# Run property-based tests
make property-test

# Run concurrency tests
make concurrency-test
```

### Writing Tests

- Write unit tests for all public functions
- Use property-based testing for mathematical operations
- Test edge cases and error conditions
- Include benchmark tests for performance-critical code
- Use testify/assert for assertions

Example test structure:
```go
func TestOperationName(t *testing.T) {
    tests := []struct {
        name     string
        input    *Polyhedron
        expected *Polyhedron
        wantErr  bool
    }{
        // test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := OperationName(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

## Code Quality

### Linting

We use golangci-lint for code quality checks:

```bash
# Run linter
make lint

# Check formatting
make check-fmt
```

### Pre-commit Checklist

Before submitting a PR, ensure:

- [ ] Code is formatted with `go fmt`
- [ ] All linter warnings are addressed
- [ ] All tests pass
- [ ] Coverage doesn't decrease significantly
- [ ] Documentation is updated if needed
- [ ] Examples still work

## Submitting Changes

### Pull Request Process

1. **Branch**: Create a feature branch from `main`
   ```bash
   git checkout -b feature/my-new-feature
   ```

2. **Commit**: Make atomic commits with clear messages
   ```bash
   git commit -m "Add kis operation implementation
   
   - Implement kis (stellation) operation
   - Add comprehensive tests with edge cases
   - Update documentation and examples"
   ```

3. **Test**: Ensure all tests pass and coverage is maintained
   ```bash
   make ci
   ```

4. **Push**: Push your branch to your fork
   ```bash
   git push origin feature/my-new-feature
   ```

5. **PR**: Create a pull request with a clear description

### Pull Request Guidelines

- **Title**: Use a clear, descriptive title
- **Description**: Explain what changes you made and why
- **Tests**: Include tests for new functionality
- **Documentation**: Update documentation if needed
- **Breaking Changes**: Clearly mark any breaking changes

## Areas for Contribution

### High Priority

- **New Operations**: Implement additional Conway operations
- **Performance**: Optimize existing operations
- **Validation**: Improve polyhedron validation
- **Documentation**: Add more examples and tutorials

### Medium Priority

- **Visualization**: Add 3D rendering capabilities
- **Import/Export**: Support for common 3D file formats
- **Web Interface**: Browser-based polyhedron viewer
- **CLI Tool**: Command-line interface for operations

### Low Priority

- **GPU Acceleration**: CUDA/OpenCL support for large polyhedra
- **Distributed Computing**: Support for cluster processing
- **Machine Learning**: AI-based operation optimization

## Reporting Issues

### Bug Reports

When reporting bugs, include:

- Go version and OS
- Steps to reproduce
- Expected vs actual behavior
- Minimal code example
- Stack trace if applicable

### Feature Requests

For new features, provide:

- Clear description of the feature
- Use case and motivation
- Proposed API design
- Implementation ideas (if any)

## Documentation

### Code Documentation

- Use godoc format for public APIs
- Include examples in godoc comments
- Document complex algorithms
- Explain mathematical concepts

### User Documentation

- Update README.md for user-facing changes
- Add examples for new features
- Update CHANGES.md for releases

## Community

### Communication

- GitHub Issues for bug reports and feature requests
- GitHub Discussions for questions and general discussion
- Pull Requests for code contributions

### Code of Conduct

Be respectful, inclusive, and constructive in all interactions. We want to maintain a welcoming environment for all contributors.

## Recognition

Contributors will be recognized in:

- CHANGES.md for significant contributions
- README.md contributors section
- Git commit history

Thank you for contributing to the Conway Polyhedron Notation Library!