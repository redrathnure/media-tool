/*
Package cmd provides command handlers

Copyright © 2020 Maksym Medvedev <redrathnure@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

const (
	version = "1.6.1"
)

var cfgFile string

var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "media-tool",
	Short: "Tooling to handle video and photo content",
	Long: `Application for importing and correction of video and photo 
	materials from digital video about photo cameras.`,
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func init() {
	initLogger()

	cobra.OnInitialize(initLoggerLevel)
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.media-tool/media-tool.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print debug messages")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("media-tool")
		viper.SetConfigType("yml")

		if ex, err := os.Executable(); err == nil {
			rootConfigDir := path.Join(filepath.Dir(ex), "conf")
			viper.AddConfigPath(rootConfigDir)
		}

		viper.AddConfigPath("/etc/media-tool")
		viper.AddConfigPath("$HOME/.media-tool")
		viper.AddConfigPath("./conf")
	}

	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debugf("Using config file: %s", viper.ConfigFileUsed())
	}
}
