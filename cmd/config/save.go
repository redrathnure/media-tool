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
