package create

import (
	"io"
	"os"
	"path/filepath"
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
