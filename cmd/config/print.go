package config

import (
	"os"

	"github.com/redrathnure/media-tool/core/tools"
	"github.com/spf13/cobra"
)

var printCmd = &cobra.Command{
	Use:     "print",
	Short:   "Print configuration",
	Long:    `Display current configuration.`,
	Args:    cobra.RangeArgs(0, 0),
	Aliases: []string{"print", "printconfig"},
	Run: func(cmd *cobra.Command, args []string) {
		tools.PrintCommandArgs(cmd, args, log)

		dstFile := tools.ExtractPath(args, 0, "")
		log.Infof("dst: '%s'", dstFile)

		data, err := getConf().GetAsYaml()
		if err != nil {
			log.Errorf("Unable to generate config: %v", err)
			os.Exit(1)
		}
		log.Infof("%s", data)
	},
}

func init() {
	configCmd.AddCommand(printCmd)
}
