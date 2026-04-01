//go:build e2e

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}
	root := filepath.Clean(filepath.Join(wd, "..", ".."))
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Fatalf("failed to resolve repo root from %s: %v", wd, err)
	}
	return root
}

func runCmd(t *testing.T, dir string, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %s %s\n%s", name, strings.Join(args, " "), string(out))
	}
	return string(out)
}

func TestCreateFullProfileE2E(t *testing.T) {
	root := repoRoot(t)
	outDir := t.TempDir()

	runCmd(t, root, "go", "run", ".", "create", "fullsvc",
		"--output", outDir,
		"--module", "github.com/acme/fullsvc",
		"--go", "1.22",
		"--skip-tidy",
		"--profile", "full",
	)

	generated := filepath.Join(outDir, "fullsvc")
	mainContent, err := os.ReadFile(filepath.Join(generated, "main.go"))
	if err != nil {
		t.Fatalf("read generated main.go error: %v", err)
	}
	if !strings.Contains(string(mainContent), "github.com/acme/fullsvc/cmd") {
		t.Fatalf("generated imports do not use custom module path")
	}

	runCmd(t, generated, "go", "mod", "tidy")
	runCmd(t, generated, "go", "test", "./...")
	runCmd(t, generated, "go", "build", "./...")
}

func TestCreateMinimalProfileE2E(t *testing.T) {
	root := repoRoot(t)
	outDir := t.TempDir()

	runCmd(t, root, "go", "run", ".", "create", "minsvc",
		"--output", outDir,
		"--profile", "minimal",
		"--skip-tidy",
	)

	generated := filepath.Join(outDir, "minsvc")
	if _, err := os.Stat(filepath.Join(generated, "docs")); !os.IsNotExist(err) {
		t.Fatalf("minimal profile should not generate docs directory")
	}
	if _, err := os.Stat(filepath.Join(generated, "internal", "database")); !os.IsNotExist(err) {
		t.Fatalf("minimal profile should not generate mysql database package")
	}

	runCmd(t, generated, "go", "mod", "tidy")
	runCmd(t, generated, "go", "test", "./...")
	runCmd(t, generated, "go", "build", "./...")
}

func TestCreateDryRunE2E(t *testing.T) {
	root := repoRoot(t)
	outDir := t.TempDir()

	runCmd(t, root, "go", "run", ".", "create", "drysvc",
		"--output", outDir,
		"--profile", "api",
		"--dry-run",
	)

	if _, err := os.Stat(filepath.Join(outDir, "drysvc")); !os.IsNotExist(err) {
		t.Fatalf("dry-run should not create any scaffold files")
	}
}

func TestCreateProductionProfileE2E(t *testing.T) {
	root := repoRoot(t)
	outDir := t.TempDir()

	runCmd(t, root, "go", "run", ".", "create", "prodsvc",
		"--output", outDir,
		"--profile", "production",
		"--skip-tidy",
	)

	generated := filepath.Join(outDir, "prodsvc")
	if _, err := os.Stat(filepath.Join(generated, "cmd", "migrate.go")); err != nil {
		t.Fatalf("production profile should generate cmd/migrate.go: %v", err)
	}
	if _, err := os.Stat(filepath.Join(generated, "internal", "middleware", "auth.go")); err != nil {
		t.Fatalf("production profile should generate auth middleware: %v", err)
	}
	if _, err := os.Stat(filepath.Join(generated, "migrations", "000001_init.up.sql")); err != nil {
		t.Fatalf("production profile should generate migration SQL: %v", err)
	}
	if _, err := os.Stat(filepath.Join(generated, "docs")); !os.IsNotExist(err) {
		t.Fatalf("production profile should disable swagger docs by default")
	}

	runCmd(t, generated, "go", "mod", "tidy")
	runCmd(t, generated, "go", "test", "./...")
	runCmd(t, generated, "go", "build", "./...")
}
