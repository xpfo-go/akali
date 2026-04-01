# Akali Hardening Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `akali` production-ready baseline by fixing installability, CLI error semantics, template generation reliability, tests, and CI.

**Architecture:** Keep current Cobra + embedded templates structure, but refactor error flow end-to-end: command layer returns errors, template generator propagates failures, and system helpers avoid `log.Fatalf`. Add deterministic tests with temp dirs and add GitHub Actions verification pipelines.

**Tech Stack:** Go 1.21+, Cobra, Go testing package, GitHub Actions.

---

### Task 1: Fix module path and import consistency

**Files:**
- Modify: `go.mod`, `main.go`, `cmd/akali/root.go`, `cmd/akali/version.go`, `internal/command/create/create.go`, `tpl/embed_server.go`, `tpl/embed_server_test.go`

- [ ] Step 1: change module path to `github.com/xpfo-go/akali`
- [ ] Step 2: update all internal imports from `github.com/xfpo-go/akali/...` to `github.com/xpfo-go/akali/...`
- [ ] Step 3: run `go test ./...` to confirm compile/import integrity

### Task 2: Build reliable error handling for CLI and template generation

**Files:**
- Modify: `main.go`, `cmd/akali/root.go`, `internal/command/create/cmd.go`, `internal/command/create/create.go`, `internal/pkg/system/system.go`, `tpl/embed_server.go`
- Test: `internal/command/create/create_test.go`, `tpl/embed_server_test.go`, `internal/pkg/system/system_test.go`

- [ ] Step 1: convert `create` command from `Run` to `RunE`
- [ ] Step 2: make `create` return explicit errors for invalid args/target exists/template failure/tidy failure
- [ ] Step 3: make `main` exit with non-zero code when `Execute` fails
- [ ] Step 4: replace `log.Fatalf` helper behavior with error returns
- [ ] Step 5: make recursive template generation return and propagate errors
- [ ] Step 6: add tests first for failure and success paths, then implement minimal code to satisfy tests

### Task 3: Stabilize tests and remove side effects

**Files:**
- Modify: `internal/pkg/system/system_test.go`, `tpl/embed_server_test.go`

- [ ] Step 1: replace side-effect tests writing to repo dirs with `t.TempDir()` based tests
- [ ] Step 2: add assertions (no print-only tests)
- [ ] Step 3: ensure tests run repeatably (`go test ./...` twice)

### Task 4: Versioning/build metadata and repo hygiene

**Files:**
- Modify: `internal/command/version/version.go`, `makefile`, `.gitignore`, `README.md`
- Delete tracked artifacts: `.idea/*`, `bin/*`

- [ ] Step 1: set version defaults to dev values and wire ldflags injection in make targets
- [ ] Step 2: ignore IDE and generated binaries in `.gitignore`
- [ ] Step 3: remove tracked `.idea` and `bin` artifacts from repository
- [ ] Step 4: update README with install, usage, and maintenance sections

### Task 5: Add CI workflow and verify pipeline locally

**Files:**
- Create: `.github/workflows/ci.yml`

- [ ] Step 1: add workflow on push/PR to run `go test ./...`, `go vet ./...`, `go build ./...`
- [ ] Step 2: include multi-OS verification matrix (linux/mac/windows)
- [ ] Step 3: locally run same commands before commit

### Task 6: Final verification and commit

**Files:**
- Modify: all above touched files

- [ ] Step 1: run `go test ./...`, `go vet ./...`, `go build ./...`
- [ ] Step 2: run smoke test `go run . create demo` in temp dir
- [ ] Step 3: verify clean working tree for accidental generated files
- [ ] Step 4: commit with conventional message
