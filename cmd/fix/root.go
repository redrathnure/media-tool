package fix

import (
	root "github.com/redrathnure/media-tool/cmd"
	"github.com/spf13/cobra"
)

var recursively bool

var dryRun bool

// importCmd represents the import command
var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Fix media files",
	Long:  `Fix metadata of video and photos.`,
}

func init() {
	root.RootCmd.AddCommand(fixCmd)

	fixCmd.PersistentFlags().BoolVarP(&dryRun, "dry", "d", false, "Dry run")
	fixCmd.PersistentFlags().BoolVarP(&recursively, "recursively", "r", false, "also analyze child directories")
}
