package config

import (
	"os"

	"github.com/redrathnure/media-tool/core/tools"
	"github.com/spf13/cobra"
)

var saveCmd = &cobra.Command{
	Use:     "save targetFile",
	Short:   "Save configuration to file",
	Long:    `Save current configuration to yaml file.`,
	Args:    cobra.RangeArgs(1, 1),
	Aliases: []string{"GenConfig", "genconfig"},
	Run: func(cmd *cobra.Command, args []string) {
		tools.PrintCommandArgs(cmd, args, log)

		dstFile := tools.ExtractPath(args, 0, "")
		log.Infof("dst: '%s'", dstFile)

		if err := getConf().SaveConfig(dstFile); err != nil {
			log.Errorf("Unable to write config: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	configCmd.AddCommand(saveCmd)
}
