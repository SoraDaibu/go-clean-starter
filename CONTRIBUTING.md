# Contributing to go-clean-starter

We welcome contributions! This project follows clean architecture principles and emphasizes simplicity and maintainability.

## Quick Start

1. **Fork & Clone**
   ```bash
   git clone https://github.com/your-username/go-clean-starter.git
   cd go-clean-starter
   ```

2. **Setup Environment**
   ```bash
   cp .env.sample .env
   make build && make up
   ```

3. **Run Tests**
   ```bash
   make test
   ```

## Development Guidelines

### Code Quality
- **Keep it simple**: Follow SOLID, KISS, and DRY principles
- **Clean Architecture**: Respect the domain → service → repository → handler layers
- **Test Coverage**: Add tests for new features using `make test`
- **Code Formatting**: Run `goimports` and `golangci-lint` before committing
- **Documentation**: Update relevant docs when adding features

### Testing Requirements
- Write unit tests for business logic in `internal/service/`
- Ensure all tests pass: `make test`
- Maintain test coverage above 80%

### Commit Guidelines
- Suggest to use conventional commits: `feat:`, `fix:`, `docs:`, `refactor:`
- Keep commits atomic and focused
- Write clear commit messages explaining the "why"

## Project Structure

- `domain/` - Business entities and interfaces
- `internal/service/` - Business logic
- `internal/repository/` - Data access layer
- `internal/http/handler/` - HTTP handlers
- `internal/task/` - Background tasks

## Submitting Changes

1. **Fork to create PR**
2. **Make your changes** following the existing patterns
3. **Add tests** for new functionality
4. **Update documentation** if needed
5. **Ensure quality checks pass**:
   ```bash
   make test
   ```
6. **Submit a pull request** with:
   - Clear description of changes
   - Link to related issues
   - Screenshots for UI changes
   - Breaking change notes if applicable

## Pull Request Review Process

- All PRs require at least one approval
- Automated checks must pass (tests, linting, formatting)
- Maintainers will review
- Address feedback promptly and professionally

## Types of Contributions

### 🐛 Bug Fixes
- Include reproduction steps
- Add regression tests
- Reference the issue number

### ✨ New Features
- Discuss in an issue first for large features
- Follow existing patterns and architecture
- Include tests and documentation

### 📚 Documentation
- Fix typos, improve clarity
- Add examples and use cases
- Keep documentation up-to-date with code changes

### 🔧 Refactoring
- Maintain backward compatibility
- Explain the benefits of the refactoring
- Ensure no functionality is lost

## Need Help?

- Check existing issues for similar problems
- Review the [README.md](README.md) and [WHY.md](WHY.md) for context
- Join discussions in issues and PRs
- Open an issue for bugs or feature requests

**Remember**: Simple, maintainable code is better than clever code. 🚀
