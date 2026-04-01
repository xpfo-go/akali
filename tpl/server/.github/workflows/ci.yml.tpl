name: ci

on:
  push:
    branches: [main]
  pull_request:

permissions:
  contents: read

jobs:
  verify:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v6

      - name: Setup Go
        uses: actions/setup-go@v6
        with:
          go-version: "<xpfo{ .GoVersion }xpfo>"

      - name: Download modules
        run: go mod download

      - name: Test
        run: go test ./...

      - name: Vet
        run: go vet ./...

      - name: Build
        run: go build ./...

  lint:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v6

      - name: Setup Go
        uses: actions/setup-go@v6
        with:
          go-version: "<xpfo{ .GoVersion }xpfo>"

      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v9
        with:
          version: v1.59
          args: --timeout=5m

  security:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v6

      - name: Setup Go
        uses: actions/setup-go@v6
        with:
          go-version: "<xpfo{ .GoVersion }xpfo>"

      - name: Download modules
        run: go mod download

      - name: Run gosec
        run: |
          go install github.com/securego/gosec/v2/cmd/gosec@latest
          gosec ./...

      - name: Run govulncheck
        run: |
          go install golang.org/x/vuln/cmd/govulncheck@latest
          govulncheck ./...
