package tpl

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenServerTemplateFS(t *testing.T) {
	base := filepath.Join(t.TempDir(), "akali_gen_test")
	err := GenServerTemplateFS(ServerTemplateFSData{
		BasePath: base,
		TplData: struct {
			ProjectName   string
			ModulePath    string
			GoVersion     string
			Profile       string
			EnableMySQL   bool
			EnableRedis   bool
			EnableSwagger bool
			EnableMetrics bool
			EnableAuth    bool
			EnableRate    bool
			EnableMigrate bool
		}{
			ProjectName:   "akali_gen_test",
			ModulePath:    "akali_gen_test",
			GoVersion:     "1.21",
			Profile:       "full",
			EnableMySQL:   true,
			EnableRedis:   true,
			EnableSwagger: true,
			EnableMetrics: true,
			EnableAuth:    false,
			EnableRate:    false,
			EnableMigrate: false,
		},
	})
	if err != nil {
		t.Fatalf("GenServerTemplateFS() error = %v", err)
	}
	expectedFiles := []string{
		filepath.Join(base, "main.go"),
		filepath.Join(base, "go.mod"),
		filepath.Join(base, "internal", "server", "server.go"),
	}
	for _, file := range expectedFiles {
		if _, err := os.Stat(file); err != nil {
			t.Fatalf("expected generated file %s to exist: %v", file, err)
		}
	}
}

func TestGenServerTemplateFS_ReturnsErrorWhenTargetExists(t *testing.T) {
	base := filepath.Join(t.TempDir(), "akali_gen_test")
	data := ServerTemplateFSData{
		BasePath: base,
		TplData: struct {
			ProjectName   string
			ModulePath    string
			GoVersion     string
			Profile       string
			EnableMySQL   bool
			EnableRedis   bool
			EnableSwagger bool
			EnableMetrics bool
			EnableAuth    bool
			EnableRate    bool
			EnableMigrate bool
		}{
			ProjectName:   "akali_gen_test",
			ModulePath:    "akali_gen_test",
			GoVersion:     "1.21",
			Profile:       "full",
			EnableMySQL:   true,
			EnableRedis:   true,
			EnableSwagger: true,
			EnableMetrics: true,
			EnableAuth:    false,
			EnableRate:    false,
			EnableMigrate: false,
		},
	}
	if err := GenServerTemplateFS(data); err != nil {
		t.Fatalf("first generation failed: %v", err)
	}
	if err := GenServerTemplateFS(data); err == nil {
		t.Fatalf("second generation expected error when files already exist")
	}
}

func TestGenServerTemplateFS_ProductionTemplateIncludesOpsFiles(t *testing.T) {
	base := filepath.Join(t.TempDir(), "akali_prod_ops")
	data := ServerTemplateFSData{
		BasePath: base,
		TplData: struct {
			ProjectName   string
			ModulePath    string
			GoVersion     string
			Profile       string
			EnableMySQL   bool
			EnableRedis   bool
			EnableSwagger bool
			EnableMetrics bool
			EnableAuth    bool
			EnableRate    bool
			EnableMigrate bool
		}{
			ProjectName:   "akali_prod_ops",
			ModulePath:    "github.com/acme/akali_prod_ops",
			GoVersion:     "1.21",
			Profile:       "production",
			EnableMySQL:   true,
			EnableRedis:   true,
			EnableSwagger: false,
			EnableMetrics: true,
			EnableAuth:    true,
			EnableRate:    true,
			EnableMigrate: true,
		},
	}
	if err := GenServerTemplateFS(data); err != nil {
		t.Fatalf("GenServerTemplateFS() error = %v", err)
	}

	expectedFiles := []string{
		filepath.Join(base, ".github", "workflows", "ci.yml"),
		filepath.Join(base, "Dockerfile"),
		filepath.Join(base, ".dockerignore"),
		filepath.Join(base, "README.md"),
	}
	for _, file := range expectedFiles {
		if _, err := os.Stat(file); err != nil {
			t.Fatalf("expected generated file %s to exist: %v", file, err)
		}
	}
}

func TestGenServerTemplateFS_InitCommandAvoidsPanic(t *testing.T) {
	base := filepath.Join(t.TempDir(), "akali_init_error")
	data := ServerTemplateFSData{
		BasePath: base,
		TplData: struct {
			ProjectName   string
			ModulePath    string
			GoVersion     string
			Profile       string
			EnableMySQL   bool
			EnableRedis   bool
			EnableSwagger bool
			EnableMetrics bool
			EnableAuth    bool
			EnableRate    bool
			EnableMigrate bool
		}{
			ProjectName:   "akali_init_error",
			ModulePath:    "github.com/acme/akali_init_error",
			GoVersion:     "1.21",
			Profile:       "full",
			EnableMySQL:   true,
			EnableRedis:   true,
			EnableSwagger: true,
			EnableMetrics: true,
			EnableAuth:    false,
			EnableRate:    false,
			EnableMigrate: false,
		},
	}
	if err := GenServerTemplateFS(data); err != nil {
		t.Fatalf("GenServerTemplateFS() error = %v", err)
	}

	content, err := os.ReadFile(filepath.Join(base, "cmd", "init.go"))
	if err != nil {
		t.Fatalf("os.ReadFile(init.go) error = %v", err)
	}
	if strings.Contains(string(content), "panic(") {
		t.Fatalf("generated init.go should avoid panic-based control flow")
	}
}

func TestGenServerTemplateFS_MigrateCommandUsesRunE(t *testing.T) {
	base := filepath.Join(t.TempDir(), "akali_migrate_rune")
	data := ServerTemplateFSData{
		BasePath: base,
		TplData: struct {
			ProjectName   string
			ModulePath    string
			GoVersion     string
			Profile       string
			EnableMySQL   bool
			EnableRedis   bool
			EnableSwagger bool
			EnableMetrics bool
			EnableAuth    bool
			EnableRate    bool
			EnableMigrate bool
		}{
			ProjectName:   "akali_migrate_rune",
			ModulePath:    "github.com/acme/akali_migrate_rune",
			GoVersion:     "1.21",
			Profile:       "production",
			EnableMySQL:   true,
			EnableRedis:   false,
			EnableSwagger: false,
			EnableMetrics: true,
			EnableAuth:    true,
			EnableRate:    true,
			EnableMigrate: true,
		},
	}
	if err := GenServerTemplateFS(data); err != nil {
		t.Fatalf("GenServerTemplateFS() error = %v", err)
	}

	content, err := os.ReadFile(filepath.Join(base, "cmd", "migrate.go"))
	if err != nil {
		t.Fatalf("os.ReadFile(migrate.go) error = %v", err)
	}
	text := string(content)
	if !strings.Contains(text, "RunE:") {
		t.Fatalf("generated migrate command should use RunE")
	}
	if strings.Contains(text, "os.Exit(") {
		t.Fatalf("generated migrate command should return errors instead of os.Exit")
	}
}

func TestGenServerTemplateFS_RateLimitUsesSingleCleanupRoutine(t *testing.T) {
	base := filepath.Join(t.TempDir(), "akali_rate_limit")
	data := ServerTemplateFSData{
		BasePath: base,
		TplData: struct {
			ProjectName   string
			ModulePath    string
			GoVersion     string
			Profile       string
			EnableMySQL   bool
			EnableRedis   bool
			EnableSwagger bool
			EnableMetrics bool
			EnableAuth    bool
			EnableRate    bool
			EnableMigrate bool
		}{
			ProjectName:   "akali_rate_limit",
			ModulePath:    "github.com/acme/akali_rate_limit",
			GoVersion:     "1.21",
			Profile:       "production",
			EnableMySQL:   true,
			EnableRedis:   false,
			EnableSwagger: false,
			EnableMetrics: true,
			EnableAuth:    false,
			EnableRate:    true,
			EnableMigrate: false,
		},
	}
	if err := GenServerTemplateFS(data); err != nil {
		t.Fatalf("GenServerTemplateFS() error = %v", err)
	}

	content, err := os.ReadFile(filepath.Join(base, "internal", "middleware", "rate_limit.go"))
	if err != nil {
		t.Fatalf("os.ReadFile(rate_limit.go) error = %v", err)
	}
	text := string(content)
	if !strings.Contains(text, "sync.Once") {
		t.Fatalf("generated rate limit middleware should guard cleanup goroutine with sync.Once")
	}
	if strings.Contains(text, "go cleanupVisitors()") {
		t.Fatalf("generated rate limit middleware should not start cleanup goroutine on every middleware init")
	}
}

func TestGenServerTemplateFS_ProductionRuntimeAvoidsPanicControlFlow(t *testing.T) {
	base := filepath.Join(t.TempDir(), "akali_no_panic_runtime")
	data := ServerTemplateFSData{
		BasePath: base,
		TplData: struct {
			ProjectName   string
			ModulePath    string
			GoVersion     string
			Profile       string
			EnableMySQL   bool
			EnableRedis   bool
			EnableSwagger bool
			EnableMetrics bool
			EnableAuth    bool
			EnableRate    bool
			EnableMigrate bool
		}{
			ProjectName:   "akali_no_panic_runtime",
			ModulePath:    "github.com/acme/akali_no_panic_runtime",
			GoVersion:     "1.21",
			Profile:       "production",
			EnableMySQL:   true,
			EnableRedis:   true,
			EnableSwagger: false,
			EnableMetrics: true,
			EnableAuth:    true,
			EnableRate:    true,
			EnableMigrate: true,
		},
	}
	if err := GenServerTemplateFS(data); err != nil {
		t.Fatalf("GenServerTemplateFS() error = %v", err)
	}

	targets := []string{
		filepath.Join(base, "internal", "server", "server.go"),
		filepath.Join(base, "internal", "database", "init.go"),
		filepath.Join(base, "internal", "cache", "redis.go"),
	}
	for _, file := range targets {
		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("os.ReadFile(%s) error = %v", file, err)
		}
		if strings.Contains(string(content), "panic(") {
			t.Fatalf("runtime template should avoid panic control flow: %s", file)
		}
	}
}
