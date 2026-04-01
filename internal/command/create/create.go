package create

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

type scaffoldOptions struct {
	projectName string
	targetPath  string
	modulePath  string
	goVersion   string
	profile     string

	force    bool
	dryRun   bool
	skipTidy bool

	enableMySQL   bool
	enableRedis   bool
	enableSwagger bool
	enableMetrics bool
}

func resolveFeatureEnabled(cmd *cobra.Command, flagName string, profileDefault bool) (bool, error) {
	if !cmd.Flags().Changed(flagName) {
		return profileDefault, nil
	}
	return cmd.Flags().GetBool(flagName)
}

func resolveProfileDefaults(profile string) (enableMySQL, enableRedis, enableSwagger, enableMetrics bool, err error) {
	switch profile {
	case "minimal":
		return false, false, false, false, nil
	case "api":
		return false, false, true, true, nil
	case "full":
		return true, true, true, true, nil
	default:
		return false, false, false, false, fmt.Errorf("invalid profile %q: allowed values are minimal, api, full", profile)
	}
}

func resolveTargetPath(projectName, output string) (string, error) {
	output = strings.TrimSpace(output)
	if output == "" {
		return "", fmt.Errorf("output path cannot be empty")
	}
	if filepath.IsAbs(projectName) {
		return filepath.Clean(projectName), nil
	}
	return filepath.Clean(filepath.Join(output, projectName)), nil
}

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

func resolveScaffoldOptions(cmd *cobra.Command, args []string) (*scaffoldOptions, error) {
	projectName := args[0]
	if strings.TrimSpace(projectName) == "" {
		return nil, fmt.Errorf("project name cannot be empty")
	}

	modulePath, err := cmd.Flags().GetString("module")
	if err != nil {
		return nil, err
	}
	modulePath = strings.TrimSpace(modulePath)
	if modulePath == "" {
		modulePath = projectName
	}
	if err := validateModulePath(modulePath); err != nil {
		return nil, err
	}

	goVersionFlag, err := cmd.Flags().GetString("go")
	if err != nil {
		return nil, err
	}
	goVersion, err := resolveGoVersion(goVersionFlag)
	if err != nil {
		return nil, err
	}

	outputPath, err := cmd.Flags().GetString("output")
	if err != nil {
		return nil, err
	}
	targetPath, err := resolveTargetPath(projectName, outputPath)
	if err != nil {
		return nil, err
	}

	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		return nil, err
	}

	profile, err := cmd.Flags().GetString("profile")
	if err != nil {
		return nil, err
	}
	profile = strings.ToLower(strings.TrimSpace(profile))
	enableMySQLDefault, enableRedisDefault, enableSwaggerDefault, enableMetricsDefault, err := resolveProfileDefaults(profile)
	if err != nil {
		return nil, err
	}
	enableMySQL, err := resolveFeatureEnabled(cmd, "with-mysql", enableMySQLDefault)
	if err != nil {
		return nil, err
	}
	enableRedis, err := resolveFeatureEnabled(cmd, "with-redis", enableRedisDefault)
	if err != nil {
		return nil, err
	}
	enableSwagger, err := resolveFeatureEnabled(cmd, "with-swagger", enableSwaggerDefault)
	if err != nil {
		return nil, err
	}
	enableMetrics, err := resolveFeatureEnabled(cmd, "with-metrics", enableMetricsDefault)
	if err != nil {
		return nil, err
	}

	dryRun, err := cmd.Flags().GetBool("dry-run")
	if err != nil {
		return nil, err
	}
	skipTidy, err := cmd.Flags().GetBool("skip-tidy")
	if err != nil {
		return nil, err
	}

	return &scaffoldOptions{
		projectName:   projectName,
		targetPath:    targetPath,
		modulePath:    modulePath,
		goVersion:     goVersion,
		profile:       profile,
		force:         force,
		dryRun:        dryRun,
		skipTidy:      skipTidy,
		enableMySQL:   enableMySQL,
		enableRedis:   enableRedis,
		enableSwagger: enableSwagger,
		enableMetrics: enableMetrics,
	}, nil
}

func create(cmd *cobra.Command, args []string) error {
	options, err := resolveScaffoldOptions(cmd, args)
	if err != nil {
		return err
	}

	if _, err := os.Stat(options.targetPath); err == nil {
		if !options.force {
			return fmt.Errorf("target path already exists: %s (use --force to overwrite or --output to change directory)", options.targetPath)
		}
		if err := os.RemoveAll(options.targetPath); err != nil {
			return fmt.Errorf("failed to remove existing target path %s: %w", options.targetPath, err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to access target path %s: %w", options.targetPath, err)
	}

	projectData := &Create{
		ProjectName:   options.projectName,
		ModulePath:    options.modulePath,
		GoVersion:     options.goVersion,
		Profile:       options.profile,
		EnableMySQL:   options.enableMySQL,
		EnableRedis:   options.enableRedis,
		EnableSwagger: options.enableSwagger,
		EnableMetrics: options.enableMetrics,
	}

	skipPaths := []string{
		"server/go.sum.tpl",
	}
	if !options.enableMySQL {
		skipPaths = append(skipPaths, "server/internal/database")
	}
	if !options.enableSwagger {
		skipPaths = append(skipPaths, "server/docs")
	}
	if !options.enableRedis {
		skipPaths = append(skipPaths, "server/internal/cache")
	}

	if err := tpl.GenServerTemplateFS(tpl.ServerTemplateFSData{
		BasePath:  options.targetPath,
		TplData:   projectData,
		DryRun:    options.dryRun,
		SkipPaths: skipPaths,
	}); err != nil {
		return fmt.Errorf("failed to generate project template: %w", err)
	}

	if options.dryRun {
		cmd.Printf(
			"dry-run: scaffold=%s profile=%s mysql=%t redis=%t swagger=%t metrics=%t\n",
			options.targetPath,
			options.profile,
			options.enableMySQL,
			options.enableRedis,
			options.enableSwagger,
			options.enableMetrics,
		)
		return nil
	}

	if !options.skipTidy {
		if err := runGoModTidy(options.targetPath); err != nil {
			return err
		}
	}
	return nil
}
