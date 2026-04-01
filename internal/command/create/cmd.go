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
}
