package tpl

import (
	"github.com/xfpo-go/akali/internal/pkg/system"
	"io/fs"
	"path"
	"testing"
)

func TestCreateFile(t *testing.T) {
	path := "./a/b/c"
	system.CreateFile(path, "t.txt")
}

func TestGen(t *testing.T) {
	err := GenServerTemplateFS(ServerTemplateFSData{
		BasePath: "akali_gen_test",
		TplData: struct {
			ProjectName string
			GoVersion   string
		}{ProjectName: "akali_gen_test",
			GoVersion: system.GetSystemGoVersion()},
	})
	if err != nil {
		t.Error(err.Error())
	}
}

func TestRev(t *testing.T) {
	fileList, err := ServerTemplateFS.ReadDir("server")
	if err != nil {
		t.Log(err.Error())
		return
	}

	var f func(prefixPath string, fs []fs.DirEntry)
	f = func(prefixPath string, fs []fs.DirEntry) {
		if len(fs) == 0 {
			return
		}

		for i := range fs {
			if fs[i].IsDir() {
				println(path.Join(prefixPath, fs[i].Name()))
				tfs, _ := ServerTemplateFS.ReadDir(path.Join(prefixPath, fs[i].Name()))
				f(path.Join(prefixPath, fs[i].Name()), tfs)
			} else {
				//println(prefixPath, fs[i].Name())
			}
		}
	}
	f("server", fileList)

}
