package akali

import (
	"github.com/spf13/cobra"
	"github.com/xfpo-go/akali/internal/command/create"
)

var CmdRoot = &cobra.Command{
	Use:     "akali",
	Example: "akali create [project-name]",
	Short:   "",
	Version: "",
}

func init() {
	CmdRoot.AddCommand(create.CmdCreate)
}

// Execute executes the root command.
func Execute() error {
	return CmdRoot.Execute()
}
