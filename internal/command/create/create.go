package create

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
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

var modulePathPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._/-]*$`)
var goVersionPattern = regexp.MustCompile(`^\d+\.\d+(\.\d+)?$`)

func validateModulePath(modulePath string) error {
	if strings.TrimSpace(modulePath) == "" {
		return fmt.Errorf("module path cannot be empty")
	}
	if strings.Contains(modulePath, " ") {
		return fmt.Errorf("module path cannot contain spaces: %s", modulePath)
	}
	if !modulePathPattern.MatchString(modulePath) {
		return fmt.Errorf("invalid module path: %s", modulePath)
	}
	return nil
}

func resolveGoVersion(flagValue string) (string, error) {
	flagValue = strings.TrimSpace(flagValue)
	if flagValue == "" {
		return system.GetSystemGoVersion()
	}
	if !goVersionPattern.MatchString(flagValue) {
		return "", fmt.Errorf("invalid go version: %s", flagValue)
	}
	return flagValue, nil
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

	modulePath, err := cmd.Flags().GetString("module")
	if err != nil {
		return err
	}
	modulePath = strings.TrimSpace(modulePath)
	if modulePath == "" {
		modulePath = projectName
	}
	if err := validateModulePath(modulePath); err != nil {
		return err
	}

	goVersionFlag, err := cmd.Flags().GetString("go")
	if err != nil {
		return err
	}
	goVersion, err := resolveGoVersion(goVersionFlag)
	if err != nil {
		return err
	}

	skipTidy, err := cmd.Flags().GetBool("skip-tidy")
	if err != nil {
		return err
	}

	projectData := &Create{
		ProjectName: projectName,
		ModulePath:  modulePath,
		GoVersion:   goVersion,
	}

	if err := tpl.GenServerTemplateFS(tpl.ServerTemplateFSData{
		BasePath: projectName,
		TplData:  projectData,
	}); err != nil {
		return fmt.Errorf("failed to generate project template: %w", err)
	}

	if !skipTidy {
		if err := runGoModTidy(projectName); err != nil {
			return err
		}
	}
	return nil
}
