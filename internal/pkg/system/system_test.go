package system

import (
	"fmt"
	"os/exec"
	"regexp"
	"testing"
)

func TestPwd(t *testing.T) {
	t.Log(Pwd())
}

func TestCreateFolder(t *testing.T) {
	CreateFolder("123")
}

func TestGetGoVersion(t *testing.T) {
	cmd := exec.Command("go", "version")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 提取 Go 版本
	version := string(out)
	version = "adssadgo1.18.2asdasdasd"
	// 使用正则表达式匹配 Go 版本号
	r := regexp.MustCompile(`go(\d+.\d+)`)
	matches := r.FindStringSubmatch(version)
	if len(matches) != 2 {
		fmt.Println("无法提取 Go 版本号")
		return
	}

	// 输出 Go 版本号
	fmt.Println(matches[1])
}
