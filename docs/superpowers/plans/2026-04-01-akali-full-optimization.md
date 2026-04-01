# Akali Full Optimization Plan

Goal: deliver production-grade scaffold ergonomics and release automation in one iteration.

Scope:
1. Release automation (GoReleaser + tag workflow)
2. Generated-project E2E tests (default and flag combinations)
3. Template dependency governance (minimal + conditional)
4. Scaffold profiles (`minimal|api|full`) and feature toggles
5. CLI ergonomics (`--output`, `--force`, `--dry-run`) with explicit errors
6. CI quality gates (`golangci-lint`, `gosec`, `govulncheck`)
7. Documentation upgrade

Execution order:
- Implement CLI/data-model first
- Add template conditions and file filtering
- Add E2E tests and wire CI
- Add release automation and security/lint gates
- Run final verification and push
