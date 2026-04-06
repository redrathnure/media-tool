package clean

import (
	root "github.com/redrathnure/media-tool/cmd"
	"github.com/spf13/cobra"
)

var recursively bool

var dryRun bool

// importCmd represents the import command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleanup media files",
	Long:  `Cleanup names adn attributes of video and photos.`,
}

func init() {
	root.RootCmd.AddCommand(cleanCmd)

	cleanCmd.PersistentFlags().BoolVarP(&dryRun, "dry", "d", false, "Dry run")
	cleanCmd.PersistentFlags().BoolVarP(&recursively, "recursively", "r", false, "also analyze child directories")
}
