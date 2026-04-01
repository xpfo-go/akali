package create

const GitKeep = ".gitkeep"

type Create struct {
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
}
