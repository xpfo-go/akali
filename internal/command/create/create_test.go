package create

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func resetCreateCommandFlags() {
	CmdCreate.Flags().VisitAll(func(f *pflag.Flag) {
		_ = f.Value.Set(f.DefValue)
		f.Changed = false
	})
}

func TestCmdCreateUsesRunE(t *testing.T) {
	resetCreateCommandFlags()
	if CmdCreate.RunE == nil {
		t.Fatalf("CmdCreate.RunE must be set")
	}
}

func TestCreateCommandRejectsExistingDirectory(t *testing.T) {
	resetCreateCommandFlags()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("os.Chdir() error = %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })

	if err := os.Mkdir(filepath.Join(tmp, "demo"), 0o755); err != nil {
		t.Fatalf("os.Mkdir() error = %v", err)
	}

	root := &cobra.Command{Use: "akali"}
	root.SilenceErrors = true
	root.SilenceUsage = true
	root.SetOut(io.Discard)
	root.SetErr(io.Discard)
	root.AddCommand(CmdCreate)
	root.SetArgs([]string{"create", "demo"})

	if err := root.Execute(); err == nil {
		t.Fatalf("expected create command to fail for existing directory")
	}
}

func TestCreateCommandGeneratesWithCustomModuleAndGoVersion(t *testing.T) {
	resetCreateCommandFlags()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("os.Chdir() error = %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })

	origRunGoModTidy := runGoModTidy
	runGoModTidy = func(projectName string) error { return nil }
	t.Cleanup(func() { runGoModTidy = origRunGoModTidy })

	root := &cobra.Command{Use: "akali"}
	root.SilenceErrors = true
	root.SilenceUsage = true
	root.SetOut(io.Discard)
	root.SetErr(io.Discard)
	root.AddCommand(CmdCreate)
	root.SetArgs([]string{
		"create", "demo",
		"--module", "github.com/acme/demo",
		"--go", "1.22",
		"--skip-tidy",
	})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	goMod, err := os.ReadFile(filepath.Join(tmp, "demo", "go.mod"))
	if err != nil {
		t.Fatalf("os.ReadFile(go.mod) error = %v", err)
	}
	content := string(goMod)
	if !strings.Contains(content, "module github.com/acme/demo") {
		t.Fatalf("go.mod module mismatch, content=%q", content)
	}
	normalized := strings.ReplaceAll(content, "\r\n", "\n")
	if !strings.Contains(normalized, "\ngo 1.22\n") {
		t.Fatalf("go.mod go version mismatch, content=%q", content)
	}
}

func TestCreateCommandSkipTidyDoesNotRunTidy(t *testing.T) {
	resetCreateCommandFlags()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("os.Chdir() error = %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })

	callCount := 0
	origRunGoModTidy := runGoModTidy
	runGoModTidy = func(projectName string) error {
		callCount++
		return nil
	}
	t.Cleanup(func() { runGoModTidy = origRunGoModTidy })

	root := &cobra.Command{Use: "akali"}
	root.SilenceErrors = true
	root.SilenceUsage = true
	root.SetOut(io.Discard)
	root.SetErr(io.Discard)
	root.AddCommand(CmdCreate)
	root.SetArgs([]string{"create", "demo", "--skip-tidy"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}
	if callCount != 0 {
		t.Fatalf("runGoModTidy() called %d times, want 0", callCount)
	}
}

func TestCreateCommandRejectsInvalidModulePath(t *testing.T) {
	resetCreateCommandFlags()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("os.Chdir() error = %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })

	root := &cobra.Command{Use: "akali"}
	root.SilenceErrors = true
	root.SilenceUsage = true
	root.SetOut(io.Discard)
	root.SetErr(io.Discard)
	root.AddCommand(CmdCreate)
	root.SetArgs([]string{"create", "demo", "--module", "invalid module"})

	if err := root.Execute(); err == nil {
		t.Fatalf("expected error for invalid module path")
	}
}

func TestCreateCommandDryRunCreatesNoFiles(t *testing.T) {
	resetCreateCommandFlags()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("os.Chdir() error = %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })

	var out bytes.Buffer
	root := &cobra.Command{Use: "akali"}
	root.SilenceErrors = true
	root.SilenceUsage = true
	root.SetOut(&out)
	root.SetErr(io.Discard)
	root.AddCommand(CmdCreate)
	root.SetArgs([]string{"create", "demo", "--dry-run"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "demo")); !os.IsNotExist(err) {
		t.Fatalf("dry-run should not create scaffold directory")
	}
	if !strings.Contains(out.String(), "dry-run") {
		t.Fatalf("expected dry-run output, got: %q", out.String())
	}
}

func TestCreateCommandOutputFlagCreatesInTargetDir(t *testing.T) {
	resetCreateCommandFlags()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("os.Chdir() error = %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })

	root := &cobra.Command{Use: "akali"}
	root.SilenceErrors = true
	root.SilenceUsage = true
	root.SetOut(io.Discard)
	root.SetErr(io.Discard)
	root.AddCommand(CmdCreate)
	root.SetArgs([]string{"create", "demo", "--output", "work", "--skip-tidy"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "work", "demo", "go.mod")); err != nil {
		t.Fatalf("expected generated go.mod in output path: %v", err)
	}
}

func TestCreateCommandForceAllowsOverwrite(t *testing.T) {
	resetCreateCommandFlags()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("os.Chdir() error = %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })

	target := filepath.Join(tmp, "demo")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("os.MkdirAll() error = %v", err)
	}
	legacy := filepath.Join(target, "legacy.txt")
	if err := os.WriteFile(legacy, []byte("legacy"), 0o644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	root := &cobra.Command{Use: "akali"}
	root.SilenceErrors = true
	root.SilenceUsage = true
	root.SetOut(io.Discard)
	root.SetErr(io.Discard)
	root.AddCommand(CmdCreate)
	root.SetArgs([]string{"create", "demo", "--force", "--skip-tidy"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(target, "go.mod")); err != nil {
		t.Fatalf("expected scaffold go.mod after force overwrite: %v", err)
	}
	if _, err := os.Stat(legacy); !os.IsNotExist(err) {
		t.Fatalf("expected legacy file removed by force overwrite")
	}
}

func TestCreateCommandProfileMinimalDisablesSwaggerAndMetrics(t *testing.T) {
	resetCreateCommandFlags()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("os.Chdir() error = %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })

	root := &cobra.Command{Use: "akali"}
	root.SilenceErrors = true
	root.SilenceUsage = true
	root.SetOut(io.Discard)
	root.SetErr(io.Discard)
	root.AddCommand(CmdCreate)
	root.SetArgs([]string{"create", "demo", "--profile", "minimal", "--skip-tidy"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "demo", "docs")); !os.IsNotExist(err) {
		t.Fatalf("minimal profile should not generate docs directory")
	}
	routerFile := filepath.Join(tmp, "demo", "internal", "api", "basic", "router.go")
	content, err := os.ReadFile(routerFile)
	if err != nil {
		t.Fatalf("os.ReadFile(router.go) error = %v", err)
	}
	if strings.Contains(string(content), "/metrics") || strings.Contains(string(content), "/swagger/") {
		t.Fatalf("minimal profile router should not include metrics/swagger routes")
	}

	goModContent, err := os.ReadFile(filepath.Join(tmp, "demo", "go.mod"))
	if err != nil {
		t.Fatalf("os.ReadFile(go.mod) error = %v", err)
	}
	mod := string(goModContent)
	if strings.Contains(mod, "swaggo") || strings.Contains(mod, "prometheus") || strings.Contains(mod, "sqlx") {
		t.Fatalf("minimal profile go.mod should not include swagger/metrics/mysql deps: %s", mod)
	}
}

func TestCreateCommandProfileProductionEnablesHardeningModules(t *testing.T) {
	resetCreateCommandFlags()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("os.Chdir() error = %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })

	root := &cobra.Command{Use: "akali"}
	root.SilenceErrors = true
	root.SilenceUsage = true
	root.SetOut(io.Discard)
	root.SetErr(io.Discard)
	root.AddCommand(CmdCreate)
	root.SetArgs([]string{"create", "demo", "--profile", "production", "--skip-tidy"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmp, "demo", "cmd", "migrate.go")); err != nil {
		t.Fatalf("production profile should generate migrate command: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "demo", "internal", "middleware", "auth.go")); err != nil {
		t.Fatalf("production profile should generate auth middleware: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "demo", "internal", "middleware", "rate_limit.go")); err != nil {
		t.Fatalf("production profile should generate rate limit middleware: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "demo", "migrations", "000001_init.up.sql")); err != nil {
		t.Fatalf("production profile should generate migration file: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "demo", "docs")); !os.IsNotExist(err) {
		t.Fatalf("production profile should disable swagger docs by default")
	}

	goModContent, err := os.ReadFile(filepath.Join(tmp, "demo", "go.mod"))
	if err != nil {
		t.Fatalf("os.ReadFile(go.mod) error = %v", err)
	}
	mod := string(goModContent)
	if !strings.Contains(mod, "github.com/golang-jwt/jwt/v5") ||
		!strings.Contains(mod, "github.com/golang-migrate/migrate/v4") ||
		!strings.Contains(mod, "golang.org/x/time") {
		t.Fatalf("production profile go.mod missing hardening deps: %s", mod)
	}
}

func TestCreateCommandRejectsMigrationWithoutMySQL(t *testing.T) {
	resetCreateCommandFlags()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("os.Chdir() error = %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })

	root := &cobra.Command{Use: "akali"}
	root.SilenceErrors = true
	root.SilenceUsage = true
	root.SetOut(io.Discard)
	root.SetErr(io.Discard)
	root.AddCommand(CmdCreate)
	root.SetArgs([]string{"create", "demo", "--profile", "production", "--with-mysql=false"})

	if err := root.Execute(); err == nil {
		t.Fatalf("expected error when migration is enabled without mysql")
	}
}
