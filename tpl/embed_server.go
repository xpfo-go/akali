package tpl

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"strings"
	"text/template"

	"github.com/xpfo-go/akali/internal/pkg/system"
)

//go:embed server
var ServerTemplateFS embed.FS

const (
	serverTemplateName = "server"
	gitKeepFile        = "gitkeep"
	templateSuffix     = ".tpl"
	delimitLeft        = "<xpfo{"
	delimitRight       = "}xpfo>"
)

type ServerTemplateFSData struct {
	BasePath  string
	TplData   any
	DryRun    bool
	SkipPaths []string
}

func GenServerTemplateFS(data ServerTemplateFSData) error {
	fileList, err := ServerTemplateFS.ReadDir(serverTemplateName)
	if err != nil {
		return err
	}

	return recursionServerTemplateFS(serverTemplateName, fileList, data)
}

func recursionServerTemplateFS(prefixPath string, fs []fs.DirEntry, data ServerTemplateFSData) error {
	if len(fs) == 0 {
		return nil
	}

	for i := range fs {
		if fs[i].IsDir() {
			nextPath := path.Join(prefixPath, fs[i].Name())
			if shouldSkipPath(nextPath, data.SkipPaths) {
				continue
			}
			tfs, err := ServerTemplateFS.ReadDir(nextPath)
			if err != nil {
				return err
			}
			if err := recursionServerTemplateFS(nextPath, tfs, data); err != nil {
				return err
			}
		} else {
			if err := genServerTemplateFSFile(prefixPath, fs[i], data); err != nil {
				return err
			}
		}
	}
	return nil
}

func genServerTemplateFSFile(prefixPath string, fs fs.DirEntry, data ServerTemplateFSData) error {
	if fs.IsDir() {
		return nil
	}

	// 1.获取tplPath
	tplPath := path.Join(prefixPath, fs.Name())
	if shouldSkipPath(tplPath, data.SkipPaths) {
		return nil
	}

	// 获取filePath
	prefixList := strings.Split(prefixPath, "/")
	prefixList[0] = data.BasePath
	filePath := path.Join(prefixList...)

	// 获取fileName
	fileName := strings.TrimSuffix(fs.Name(), templateSuffix)
	if fileName == gitKeepFile {
		fileName = "." + fileName
	}

	// 2. 生成文件
	tplFile, err := template.New(fs.Name()).Delims(delimitLeft, delimitRight).ParseFS(ServerTemplateFS, tplPath)
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", tplPath, err)
	}
	if data.DryRun {
		return nil
	}
	f, err := system.CreateFile(filePath, fileName)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	if err := tplFile.Execute(f, data.TplData); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", tplPath, err)
	}
	return nil
}

func shouldSkipPath(target string, skipPaths []string) bool {
	target = strings.TrimPrefix(path.Clean(target), "./")
	for _, skip := range skipPaths {
		skip = strings.TrimPrefix(path.Clean(skip), "./")
		if target == skip || strings.HasPrefix(target, skip+"/") {
			return true
		}
	}
	return false
}
