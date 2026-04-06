package config

import (
	root "github.com/redrathnure/media-tool/cmd"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "App configuration",
	Long:  `Opts to manage app configuration.`,
}

func init() {
	root.RootCmd.AddCommand(configCmd)
}
