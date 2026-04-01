package akali

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xpfo-go/akali/internal/command/version"
	"strings"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print build metadata of akali",
	Long:  `Print detailed build metadata for akali.`,
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
