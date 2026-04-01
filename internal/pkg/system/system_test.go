package system

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func TestCreateFolder(t *testing.T) {
	target := filepath.Join(t.TempDir(), "nested", "dir")
	if err := CreateFolder(target); err != nil {
		t.Fatalf("CreateFolder() error = %v", err)
	}
	if info, err := os.Stat(target); err != nil || !info.IsDir() {
		t.Fatalf("CreateFolder() did not create target directory, statErr=%v", err)
	}
}

func TestCreateFile(t *testing.T) {
	dir := t.TempDir()
	file, err := CreateFile(dir, "a.txt")
	if err != nil {
		t.Fatalf("CreateFile() unexpected error = %v", err)
	}
	_ = file.Close()

	if _, err := CreateFile(dir, "a.txt"); err == nil {
		t.Fatalf("CreateFile() expected error when file exists")
	}
}

func TestGetGoVersion(t *testing.T) {
	version, err := GetSystemGoVersion()
	if err != nil {
		t.Fatalf("GetSystemGoVersion() error = %v", err)
	}
	r := regexp.MustCompile(`^\d+\.\d+`)
	if !r.MatchString(version) {
		t.Fatalf("GetSystemGoVersion() invalid version: %q", version)
	}
}
