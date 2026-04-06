package import_

import (
	root "github.com/redrathnure/media-tool/cmd"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import media data",
	Long:  `Import video and photos from different medias.`,
}

// dryRun just test instead real file manimupations
var dryRun bool

func init() {
	root.RootCmd.AddCommand(importCmd)

	importCmd.PersistentFlags().BoolVarP(&dryRun, "dry", "d", false, "Dry run")
}
