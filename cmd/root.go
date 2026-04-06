package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "media-tool",
	Short: "Tooling to handle video and photo content",
	Long: `Application for importing and correction of video and photo
	materials from digital video about photo cameras.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func init() {
	initLogger()

	cobra.OnInitialize(initLoggerLevel)
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&Context.cfgFile, "config", "c", "", "config file (default is $HOME/.media-tool/media-tool.yaml)")
	RootCmd.PersistentFlags().BoolVarP(&Context.Verbose, "verbose", "v", false, "Print debug messages")
}
