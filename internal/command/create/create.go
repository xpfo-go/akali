package create

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xfpo-go/akali/internal/pkg/system"
	"github.com/xfpo-go/akali/tpl"
	"os/exec"
)

func create(cmd *cobra.Command, args []string) {
	projectName := args[0]
	projectData := &Create{
		ProjectName: projectName,
		GoVersion:   system.GetSystemGoVersion(),
	}

	if err := tpl.GenServerTemplateFS(tpl.ServerTemplateFSData{
		BasePath: projectName,
		TplData:  projectData,
	}); err != nil {
		fmt.Println(err.Error())
	}

	// exec go mod tidy
	c := exec.Command("go", "mod", "tidy")
	c.Dir = projectName
	_ = c.Run()
}
