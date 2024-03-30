package create

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xfpo-go/akali/internal/pkg/system"
	"github.com/xfpo-go/akali/tpl"
	"html/template"
	"os/exec"
)

func create(cmd *cobra.Command, args []string) {
	projectName := args[0]
	projectData := &Create{
		ProjectName: projectName,
		GoVersion:   system.GetSystemGoVersion(),
		Lt:          "<",
	}

	// bin
	system.CreateFile(fmt.Sprintf("%s/bin", projectName), GitKeep)

	// cmd
	genFile(fmt.Sprintf("%s/cmd", projectName), "admin.go", "server/cmd.admin.go.tpl", projectData)
	genFile(fmt.Sprintf("%s/cmd", projectName), "init.go", "server/cmd.init.go.tpl", projectData)

	// config
	genFile(fmt.Sprintf("%s/config", projectName), "config.go", "server/config.config.go.tpl", projectData)

	// docs
	genFileNoData(fmt.Sprintf("%s/docs", projectName), "docs.go", "docs.go.tpl", "server/docs.go.tpl")
	genFile(fmt.Sprintf("%s/docs", projectName), "swagger.json", "server/swagger.json.tpl", projectData)
	genFile(fmt.Sprintf("%s/docs", projectName), "swagger.yaml", "server/swagger.yaml.tpl", projectData)

	// internal
	genFile(fmt.Sprintf("%s/internal/api", projectName), "router.go", "server/api.router.go.tpl", projectData)
	genFile(fmt.Sprintf("%s/internal/api/basic", projectName), "router.go", "server/api.basic.router.go.tpl", projectData)
	genFile(fmt.Sprintf("%s/internal/controller/basic", projectName), "basic.go", "server/ctrl.basic.go.tpl", projectData)
	genFile(fmt.Sprintf("%s/internal/controller/basic", projectName), "health.go", "server/ctrl.health.go.tpl", projectData)
	system.CreateFile(fmt.Sprintf("%s/internal/database/dao", projectName), GitKeep)
	system.CreateFile(fmt.Sprintf("%s/internal/database/do", projectName), GitKeep)
	system.CreateFile(fmt.Sprintf("%s/internal/database/entity", projectName), GitKeep)
	genFile(fmt.Sprintf("%s/internal/database", projectName), "init.go", "server/internal.database.init.go.tpl", projectData)
	genFile(fmt.Sprintf("%s/internal/database", projectName), "mysql.go", "server/internal.database.mysql.go.tpl", projectData)
	genFile(fmt.Sprintf("%s/internal/version", projectName), "version.go", "server/internal.version.go.tpl", projectData)

	// pkg
	genFile(fmt.Sprintf("%s/pkg/limiter", projectName), "dcs_leaky_bucket.go", "server/pkg.limiter.dcs_leaky_bucket.go.tpl", projectData)
	genFile(fmt.Sprintf("%s/pkg/limiter", projectName), "dcs_token_bucket.go", "server/pkg.limiter.dcs_token_bucket.go.tpl", projectData)
	genFile(fmt.Sprintf("%s/pkg/limiter", projectName), "leaky_bucket.go", "server/pkg.limiter.leaky_bucket.go.tpl", projectData)
	genFile(fmt.Sprintf("%s/pkg/limiter", projectName), "token_bucket.go", "server/pkg.limiter.token_bucket.go.tpl", projectData)
	genFile(fmt.Sprintf("%s/pkg/logs", projectName), "logs.go", "server/pkg.logs.go.tpl", projectData)
	genFile(fmt.Sprintf("%s/pkg/retry", projectName), "retry.go", "server/pkg.retry.go.tpl", projectData)
	genFile(fmt.Sprintf("%s/pkg/server", projectName), "router.go", "server/pkg.router.go.tpl", projectData)
	genFile(fmt.Sprintf("%s/pkg/server", projectName), "server.go", "server/pkg.server.go.tpl", projectData)

	// main.go
	genFile(fmt.Sprintf("%s", projectName), "main.go", "server/main.go.tpl", projectData)

	// go.mod
	genFile(fmt.Sprintf("%s", projectName), "go.mod", "server/go.mod.tpl", projectData)
	genFile(fmt.Sprintf("%s", projectName), "go.sum", "server/go.sum.tpl", projectData)

	// config.yaml
	genFile(fmt.Sprintf("%s", projectName), "config.yaml", "server/config.yaml.tpl", projectData)

	// makefile
	genFile(fmt.Sprintf("%s", projectName), "makefile", "server/makefile.tpl", projectData)

	// exec go mod tidy
	c := exec.Command("go", "mod", "tidy")
	c.Dir = projectName
	_ = c.Run()
}

func genFile(filePath, fileName, tplPath string, tplData any) {
	tplFile, _ := template.ParseFS(tpl.ServerTemplateFS, tplPath)
	f := system.CreateFile(filePath, fileName)
	defer func() {
		_ = f.Close()
	}()
	_ = tplFile.Execute(f, tplData)
}

func genFileNoData(filePath, fileName, tplName, tplPath string) {

	tplFile := template.Must(template.New(tplName).Delims("<a,./s222d<", ">dsa2>").ParseFS(tpl.ServerTemplateFS, tplPath))

	f := system.CreateFile(filePath, fileName)
	defer func() {
		_ = f.Close()
	}()
	_ = tplFile.Execute(f, nil)
}
