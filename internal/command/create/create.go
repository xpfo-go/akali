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
	enableAuth    bool
	enableRate    bool
	enableMigrate bool
}

func resolveFeatureEnabled(cmd *cobra.Command, flagName string, profileDefault bool) (bool, error) {
	if !cmd.Flags().Changed(flagName) {
		return profileDefault, nil
	}
	return cmd.Flags().GetBool(flagName)
}

type profileDefaults struct {
	enableMySQL   bool
	enableRedis   bool
	enableSwagger bool
	enableMetrics bool
	enableAuth    bool
	enableRate    bool
	enableMigrate bool
}

func resolveProfileDefaults(profile string) (*profileDefaults, error) {
	switch profile {
	case "minimal":
		return &profileDefaults{
			enableMySQL: false, enableRedis: false, enableSwagger: false, enableMetrics: false,
			enableAuth: false, enableRate: false, enableMigrate: false,
		}, nil
	case "api":
		return &profileDefaults{
			enableMySQL: false, enableRedis: false, enableSwagger: true, enableMetrics: true,
			enableAuth: false, enableRate: false, enableMigrate: false,
		}, nil
	case "full":
		return &profileDefaults{
			enableMySQL: true, enableRedis: true, enableSwagger: true, enableMetrics: true,
			enableAuth: false, enableRate: false, enableMigrate: false,
		}, nil
	case "production":
		return &profileDefaults{
			enableMySQL: true, enableRedis: true, enableSwagger: false, enableMetrics: true,
			enableAuth: true, enableRate: true, enableMigrate: true,
		}, nil
	default:
		return nil, fmt.Errorf("invalid profile %q: allowed values are minimal, api, full, production", profile)
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
	defaults, err := resolveProfileDefaults(profile)
	if err != nil {
		return nil, err
	}
	enableMySQL, err := resolveFeatureEnabled(cmd, "with-mysql", defaults.enableMySQL)
	if err != nil {
		return nil, err
	}
	enableRedis, err := resolveFeatureEnabled(cmd, "with-redis", defaults.enableRedis)
	if err != nil {
		return nil, err
	}
	enableSwagger, err := resolveFeatureEnabled(cmd, "with-swagger", defaults.enableSwagger)
	if err != nil {
		return nil, err
	}
	enableMetrics, err := resolveFeatureEnabled(cmd, "with-metrics", defaults.enableMetrics)
	if err != nil {
		return nil, err
	}
	enableAuth, err := resolveFeatureEnabled(cmd, "with-auth", defaults.enableAuth)
	if err != nil {
		return nil, err
	}
	enableRate, err := resolveFeatureEnabled(cmd, "with-rate-limit", defaults.enableRate)
	if err != nil {
		return nil, err
	}
	enableMigrate, err := resolveFeatureEnabled(cmd, "with-migrate", defaults.enableMigrate)
	if err != nil {
		return nil, err
	}
	if enableMigrate && !enableMySQL {
		return nil, fmt.Errorf("migration requires mysql support: enable --with-mysql or disable --with-migrate")
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
		enableAuth:    enableAuth,
		enableRate:    enableRate,
		enableMigrate: enableMigrate,
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
		EnableAuth:    options.enableAuth,
		EnableRate:    options.enableRate,
		EnableMigrate: options.enableMigrate,
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
	if !options.enableAuth {
		skipPaths = append(skipPaths, "server/internal/api/secure", "server/internal/controller/secure", "server/internal/middleware/auth.go.tpl")
	}
	if !options.enableRate {
		skipPaths = append(skipPaths, "server/internal/middleware/rate_limit.go.tpl")
	}
	if !options.enableMigrate {
		skipPaths = append(skipPaths, "server/cmd/migrate.go.tpl", "server/migrations")
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
			"dry-run: scaffold=%s profile=%s mysql=%t redis=%t swagger=%t metrics=%t auth=%t rate_limit=%t migrate=%t\n",
			options.targetPath,
			options.profile,
			options.enableMySQL,
			options.enableRedis,
			options.enableSwagger,
			options.enableMetrics,
			options.enableAuth,
			options.enableRate,
			options.enableMigrate,
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
