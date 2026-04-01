package akali

import (
	"github.com/spf13/cobra"
	"github.com/xpfo-go/akali/internal/command/create"
	akaliVersion "github.com/xpfo-go/akali/internal/command/version"
)

var CmdRoot = &cobra.Command{
	Use:           "akali",
	Example:       "akali create [project-name]",
	Short:         "Scaffold production-ready Go backend services",
	Version:       akaliVersion.Version,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	CmdRoot.AddCommand(create.CmdCreate)
}

// Execute executes the root command.
func Execute() error {
	return CmdRoot.Execute()
}
