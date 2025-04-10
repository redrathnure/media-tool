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
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleanup media files",
	Long:  `Cleanup names adn attributes of video and photos.`,
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	cleanCmd.PersistentFlags().BoolVarP(&DryRun, "dry", "d", false, "Dry run")
	cleanCmd.PersistentFlags().BoolVarP(&recursively, "recursively", "r", false, "also analyze child directories")

}
