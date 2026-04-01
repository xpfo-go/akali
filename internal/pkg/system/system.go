package system

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

func Pwd() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

func CreateFolder(dirPath string) error {
	if err := os.MkdirAll(dirPath, 0o750); err != nil {
		return fmt.Errorf("failed to create dir %s: %w", dirPath, err)
	}
	return nil
}

func CreateFile(dirPath string, filename string) (*os.File, error) {
	filePath := filepath.Join(dirPath, filename)

	if err := os.MkdirAll(dirPath, 0o750); err != nil {
		return nil, fmt.Errorf("failed to create dir %s: %w", dirPath, err)
	}
	// #nosec G304 -- filePath is built from scaffold-controlled template paths.
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil, fmt.Errorf("file already exists: %s", filePath)
		}
		return nil, fmt.Errorf("failed to create file %s: %w", filePath, err)
	}

	return file, nil
}

// GetSystemGoVersion 获取当前系统的go的版本号 out：1.21
func GetSystemGoVersion() (string, error) {
	cmd := exec.Command("go", "version")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get go version: %w", err)
	}
	r := regexp.MustCompile(`go(\d+.\d+)`)
	matches := r.FindStringSubmatch(string(out))
	if len(matches) != 2 {
		return "", fmt.Errorf("failed to parse go version from output: %q", string(out))
	}

	// 输出 Go 版本号
	return matches[1], nil
}
