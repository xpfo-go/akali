package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"<xpfo{ .ProjectName }xpfo>/internal/version"
	"strings"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of <xpfo{ .ProjectName }xpfo>",
	Long:  `All software has versions. This is <xpfo{ .ProjectName }xpfo>`,
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
	rootCmd.AddCommand(versionCmd)
}
