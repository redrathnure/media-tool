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
package version

import (
	root "github.com/redrathnure/media-tool/cmd"
	"github.com/redrathnure/media-tool/core/tools"
	"github.com/spf13/cobra"
)

// It may be overridden during build. And it will be overridden during release
var version = "1.6.2"
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
