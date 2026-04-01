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
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create dir %s: %w", dirPath, err)
	}
	return nil
}

func CreateFile(dirPath string, filename string) (*os.File, error) {
	filePath := filepath.Join(dirPath, filename)

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create dir %s: %w", dirPath, err)
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil, fmt.Errorf("file already exists: %s", filePath)
		}
		return nil, fmt.Errorf("failed to create file %s: %w", filePath, err)
	}

	return file, nil
}

func ExecCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...) // 拼接参数与命令
	if err := cmd.Run(); err != nil {  // 执行命令，若命令出错则打印错误到 stderr
		return err
	}
	return nil
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
