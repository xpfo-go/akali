package tpl

import (
	"os"
	"path/filepath"
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
