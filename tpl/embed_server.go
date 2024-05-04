package tpl

import (
	"embed"
	"github.com/xfpo-go/akali/internal/pkg/system"
	"io/fs"
	"path"
	"strings"
	"text/template"
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
	BasePath string
	TplData  any
}

func GenServerTemplateFS(data ServerTemplateFSData) error {
	fileList, err := ServerTemplateFS.ReadDir(serverTemplateName)
	if err != nil {
		return err
	}

	recursionServerTemplateFS(serverTemplateName, fileList, data)
	return nil
}

func recursionServerTemplateFS(prefixPath string, fs []fs.DirEntry, data ServerTemplateFSData) {
	if len(fs) == 0 {
		return
	}

	for i := range fs {
		if fs[i].IsDir() {
			tfs, _ := ServerTemplateFS.ReadDir(path.Join(prefixPath, fs[i].Name()))
			recursionServerTemplateFS(path.Join(prefixPath, fs[i].Name()), tfs, data)
		} else {
			genServerTemplateFSFile(prefixPath, fs[i], data)
		}
	}
}

func genServerTemplateFSFile(prefixPath string, fs fs.DirEntry, data ServerTemplateFSData) {
	if fs.IsDir() {
		return
	}

	// 1.获取tplPath
	tplPath := path.Join(prefixPath, fs.Name())

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
	tplFile, _ := template.New(fs.Name()).Delims(delimitLeft, delimitRight).ParseFS(ServerTemplateFS, tplPath)
	f := system.CreateFile(filePath, fileName)
	defer func() {
		_ = f.Close()
	}()
	_ = tplFile.Execute(f, data.TplData)
}
