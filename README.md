# akali

`akali` is a Go CLI scaffold tool for generating backend service projects.
The generated project includes a practical baseline around Gin, sqlx, Swagger,
logging, metrics, and common server wiring.

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

## Commands

```bash
akali create [project-name]
akali version
akali --version
```

## Development

Run checks locally:

```bash
go test ./...
go vet ./...
go build ./...
```

Cross-platform builds:

```bash
make build-all VERSION=v1.0.3
```

## Generated Project Stack

- Gin
- sqlx
- Zap
- Swagger (swaggo)
- Prometheus metrics

## CI

This repository uses GitHub Actions to run:

- `go test ./...`
- `go vet ./...`
- `go build ./...`

across Linux, macOS, and Windows.

## License

MIT. See [LICENSE](LICENSE).
