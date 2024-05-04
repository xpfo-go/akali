package akali

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xfpo-go/akali/internal/command/version"
	"strings"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of server-tpl",
	Long:  `All software has versions. This is server-tpl`,
	Run: func(cmd *cobra.Command, args []string) {
		info := []string{
			"Version: " + version.Version,
			"Commit: " + version.Commit,
			"Build Time: " + version.BuildTime,
			"Go Version: " + version.GoVersion,
		}
		fmt.Println(strings.Join(info, "\n"))
	},
}

func init() {
	CmdRoot.AddCommand(versionCmd)
}
