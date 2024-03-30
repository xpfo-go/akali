package system

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

func Pwd() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

func CreateFolder(dirPath string) {
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		log.Fatalf("Failed to create dir %s: %v", dirPath, err)
	}
}

func CreateFile(dirPath string, filename string) *os.File {
	filePath := filepath.Join(dirPath, filename)

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		log.Fatalf("Failed to create dir %s: %v", dirPath, err)
	}
	stat, _ := os.Stat(filePath)
	if stat != nil {
		return nil
	}
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", filePath, err)
	}

	return file
}

func ExecCommand(name string, args ...string) {
	cmd := exec.Command(name, args...) // 拼接参数与命令
	if err := cmd.Run(); err != nil {  // 执行命令，若命令出错则打印错误到 stderr
		log.Println(err)
	}
}

// GetSystemGoVersion 获取当前系统的go的版本号 out：1.21
func GetSystemGoVersion() string {
	cmd := exec.Command("go", "version")
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to get go version. err: %v", err)
	}
	r := regexp.MustCompile(`go(\d+.\d+)`)
	matches := r.FindStringSubmatch(string(out))
	if len(matches) != 2 {
		log.Fatalf("Failed to get go version.")
	}

	// 输出 Go 版本号
	return matches[1]
}
