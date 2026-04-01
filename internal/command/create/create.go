package create

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xpfo-go/akali/internal/pkg/system"
	"github.com/xpfo-go/akali/tpl"
)

var runGoModTidy = func(projectName string) error {
	c := exec.Command("go", "mod", "tidy")
	c.Dir = projectName
	output, err := c.CombinedOutput()
	if err != nil {
		if len(output) > 0 {
			return fmt.Errorf("go mod tidy failed: %w: %s", err, strings.TrimSpace(string(output)))
		}
		return fmt.Errorf("go mod tidy failed: %w", err)
	}
	return nil
}

func create(cmd *cobra.Command, args []string) error {
	projectName := args[0]
	if strings.TrimSpace(projectName) == "" {
		return fmt.Errorf("project name cannot be empty")
	}
	if _, err := os.Stat(projectName); err == nil {
		return fmt.Errorf("target path already exists: %s", projectName)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to access target path %s: %w", projectName, err)
	}

	goVersion, err := system.GetSystemGoVersion()
	if err != nil {
		return err
	}

	projectData := &Create{
		ProjectName: projectName,
		GoVersion:   goVersion,
	}

	if err := tpl.GenServerTemplateFS(tpl.ServerTemplateFSData{
		BasePath: projectName,
		TplData:  projectData,
	}); err != nil {
		return fmt.Errorf("failed to generate project template: %w", err)
	}

	if err := runGoModTidy(projectName); err != nil {
		return err
	}
	return nil
}
