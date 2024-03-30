package create

import "github.com/spf13/cobra"

var CmdCreate = &cobra.Command{
	Use:     "create [project-name]",
	Short:   "Create a new project",
	Example: "akali create [project-name]",
	Args:    cobra.ExactArgs(1),
	Run:     create,
}
