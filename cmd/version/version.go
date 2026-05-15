package version

import (
	root "github.com/redrathnure/media-tool/cmd"
	"github.com/redrathnure/media-tool/core/tools"
	"github.com/spf13/cobra"
)

// It may be overridden during build. And it will be overridden during release
var version = "1.13.0"
var goVersion = "???"
var goBuildPlatform = "???"

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Show version",
	Long:    `Display version information.`,
	Args:    cobra.NoArgs,
	Aliases: []string{"ver", "v"},
	Run: func(cmd *cobra.Command, args []string) {
		tools.PrintCommandArgs(cmd, args, log)

		log.Infof("App version: %s", version)
		log.Infof("Go version: %s", goVersion)
		log.Infof("Build platform: %s", goBuildPlatform)
	},
}

func init() {
	root.RootCmd.AddCommand(versionCmd)
}
