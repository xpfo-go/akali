package create

import "github.com/spf13/cobra"

var CmdCreate = &cobra.Command{
	Use:     "create [project-name]",
	Short:   "Create a new project",
	Example: "akali create [project-name] --module github.com/org/service --go 1.22",
	Args:    cobra.ExactArgs(1),
	RunE:    create,
}

func init() {
	CmdCreate.Flags().String("module", "", "go module path for generated project (default: project-name)")
	CmdCreate.Flags().String("go", "", "go language version for generated project (default: system go version)")
	CmdCreate.Flags().Bool("skip-tidy", false, "skip running `go mod tidy` after generation")
	CmdCreate.Flags().String("output", ".", "output directory for generated project")
	CmdCreate.Flags().Bool("force", false, "remove existing target directory before scaffold generation")
	CmdCreate.Flags().Bool("dry-run", false, "preview scaffold generation without writing files")
	CmdCreate.Flags().String("profile", "full", "scaffold profile: minimal | api | full")
	CmdCreate.Flags().Bool("with-mysql", false, "enable mysql integration (overrides profile default when set)")
	CmdCreate.Flags().Bool("with-redis", false, "enable redis integration (overrides profile default when set)")
	CmdCreate.Flags().Bool("with-swagger", false, "enable swagger docs endpoint (overrides profile default when set)")
	CmdCreate.Flags().Bool("with-metrics", false, "enable metrics endpoint (overrides profile default when set)")
}
