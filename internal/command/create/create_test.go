package create

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestCmdCreateUsesRunE(t *testing.T) {
	if CmdCreate.RunE == nil {
		t.Fatalf("CmdCreate.RunE must be set")
	}
}

func TestCreateCommandRejectsExistingDirectory(t *testing.T) {
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
	if !strings.Contains(content, "\ngo 1.22\n") {
		t.Fatalf("go.mod go version mismatch, content=%q", content)
	}
}

func TestCreateCommandSkipTidyDoesNotRunTidy(t *testing.T) {
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
