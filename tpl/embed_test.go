package tpl

import (
	"fmt"
	"github.com/xfpo-go/akali/internal/pkg/system"
	"html/template"
	"os"
	"testing"
)

func TestEmbedServer(t *testing.T) {
	type Data struct {
		ProjectName string
	}
	a := new(Data)
	a.ProjectName = "xpfo"

	tmp, err := template.ParseFS(ServerTemplateFS, "server/main.go.tpl")
	if err != nil {
		t.Error(err.Error())
		return
	}

	f := system.CreateFile("./", "main.go")
	defer f.Close()

	_ = tmp.Execute(f, a)
}

func TestCreateFile(t *testing.T) {
	path := "./a/b/c"
	system.CreateFile(path, "t.txt")
}

func TestNotParse(t *testing.T) {
	data := struct {
		Name string
		Sex  string
	}{Name: "name", Sex: "sex<"}

	tpl := template.Must(template.New("test.tpl").Delims("<<", ">>").ParseFS(ServerTemplateFS, "server/test.tpl"))

	//template.HTMLEscape()
	fmt.Println(*tpl)

	//tpl.ExecuteTemplate(os.Stdout, "test.tpl", data)
	_ = tpl.Execute(os.Stdout, data)
	fmt.Println()

	//tplName := "docs.go.tpl"
	//tplPath := "server/docs.go.tpl"
	//
	//tplFile := template.Must(template.New(tplName).Delims("<a,./s222d<", ">dsa2>").ParseFS(ServerTemplateFS, tplPath))
	//
	//_ = tplFile.Execute(os.Stdout, nil)
}
