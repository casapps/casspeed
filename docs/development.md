# Development Guide

Guide for developers contributing to casspeed.

## Prerequisites

- Docker (required for building)
- Git
- Make
- (Optional) Go 1.21+ for local development

## Getting Started

### Clone Repository

```bash
git clone https://github.com/casapps/casspeed
cd casspeed
```

### Quick Build

```bash
# Build all platforms (uses Docker)
make build

# Outputs to binaries/
```

### Development Build

```bash
# Quick dev build to temp directory
make dev

# Or using Docker Compose
docker-compose -f docker/docker-compose.dev.yml up
```

## Project Structure

```
casspeed/
├── src/                    # Source code
│   ├── config/            # Configuration handling
│   ├── mode/              # App mode detection
│   ├── paths/             # OS-specific paths
│   ├── server/            # HTTP server
│   │   ├── handler/       # Request handlers
│   │   ├── model/         # Data models
│   │   ├── service/       # Business logic
│   │   └── store/         # Database layer
│   └── client/            # CLI client
├── docker/                 # Docker files
├── docs/                   # Documentation (MkDocs)
├── tests/                  # Test files
└── web/                    # Frontend templates
```

## Building

### Using Make (Recommended)

```bash
# Build all platforms
make build

# Build Docker image
make docker

# Run tests
make test

# Clean build artifacts
make clean
```

### Direct Docker Build

```bash
# Build using Docker directly
docker run --rm \
  -v $(pwd):/build \
  -w /build \
  -e CGO_ENABLED=0 \
  golang:alpine \
  go build -o binaries/casspeed ./src
```

## Testing

### Run All Tests

```bash
make test
```

### Docker Testing

```bash
# Test with Docker Compose
docker-compose -f docker/docker-compose.test.yml up

# Manual Docker test
docker build -f docker/Dockerfile -t casspeed:test .
docker run --rm -p 8080:80 casspeed:test
```

### Incus Testing

```bash
# Test in Incus container
./tests/incus.sh
```

## Code Style

### Go Code

- Use `gofmt` for formatting
- Follow standard Go conventions
- Comments above code (never inline)
- Descriptive variable names

### File Naming

- Go files: `lowercase_snake.go`
- Packages: `lowercase` (single word)
- Functions: `PascalCase` (exported), `camelCase` (private)

## Configuration

Development config auto-generates on first run:

```yaml
server:
  mode: development
  port: 8080
  address: "127.0.0.1"

test:
  max_threads: 8
  default_duration: 10
```

## Debugging

### Enable Debug Mode

```bash
# Via flag
./casspeed --debug

# Via environment
export DEBUG=true
./casspeed

# Via config
# server.yml: debug: true
```

### Debug Endpoints

When debug mode enabled:

- `/debug/pprof/` - Go profiling
- `/debug/vars` - Runtime variables

## Contributing

### Workflow

1. Fork repository
2. Create feature branch
3. Make changes
4. Test thoroughly
5. Submit pull request

### Commit Messages

Use emoji prefixes:

- ✨ `feat:` New feature
- 🐛 `fix:` Bug fix
- 📝 `docs:` Documentation
- ♻️ `refactor:` Code refactoring
- ✅ `test:` Add tests
- 🔧 `chore:` Maintenance

### Pull Request

- Clear description
- Link related issues
- All tests passing
- Documentation updated

## Resources

- Repository: https://github.com/casapps/casspeed
- Issues: https://github.com/casapps/casspeed/issues
- Documentation: https://casapps-casspeed.readthedocs.io
