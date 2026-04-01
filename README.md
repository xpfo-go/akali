# akali

`akali` is a Go CLI scaffold tool for generating backend service projects with
production-oriented defaults.

## Install

```bash
go install github.com/xpfo-go/akali@latest
```

## Quick Start

```bash
akali create demo-service
cd demo-service
go test ./...
go run . --help
```

## Create Command

```bash
akali create [project-name] [flags]
```

Common usage:

```bash
akali create order-service \
  --output ./workspaces \
  --module github.com/your-org/order-service \
  --go 1.22 \
  --profile api
```

### Profiles

- `minimal`: HTTP baseline only (no MySQL, no Redis, no Swagger, no Metrics)
- `api`: API-facing defaults (Swagger + Metrics, no MySQL/Redis)
- `full`: full stack baseline (MySQL + Redis + Swagger + Metrics)

### Flags

- `--module`: set generated `go.mod` module path
- `--go`: set generated Go version
- `--profile`: `minimal | api | full`
- `--with-mysql`: override profile default and enable/disable MySQL
- `--with-redis`: override profile default and enable/disable Redis
- `--with-swagger`: override profile default and enable/disable Swagger
- `--with-metrics`: override profile default and enable/disable Metrics
- `--output`: output directory for generated project
- `--force`: overwrite existing target directory
- `--skip-tidy`: skip `go mod tidy` after generation
- `--dry-run`: preview generation without writing files

Examples:

```bash
# preview without writing files
akali create demo --profile minimal --dry-run

# overwrite an existing target directory
akali create demo --force --profile full

# custom feature matrix on top of profile
akali create demo --profile api --with-mysql --with-redis
```

## Development

```bash
go test ./...
go vet ./...
go build ./...
make test-e2e
```

## CI/CD

`akali` uses GitHub Actions for:

- cross-platform verification (`test`, `vet`, `build`) on Linux/macOS/Windows
- `golangci-lint`
- `gosec` + `govulncheck`
- E2E scaffold generation tests

Tag-based release:

- push tag `vX.Y.Z`
- GitHub Actions runs GoReleaser
- multi-platform binaries and checksums are published to GitHub Releases

## License

MIT. See [LICENSE](LICENSE).
