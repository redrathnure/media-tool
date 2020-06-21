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

// photosCmd represents the photos command
var photosCmd = &cobra.Command{
	Use:   "photos sourceDir [targetDir]",
	Short: "Import photos from directory",
	Long:  `Copy images from sourceDir to ditargetDirsk. By default creates subdirectories by dates.`,
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		printCommandArgs(cmd, args)

		src := extractPath(args, 0, ".")
		log.Infof("src: '%s'", src)

		dstDir := extractPath(args, 1, "d:\\tmp\\test")
		log.Infof("dst: '%s'", dstDir)

		log.Infof("dry ryn: %v", DryRun)

		tagName := "FileName"
		if DryRun {
			tagName = "TestName"
		}

		exifToolArgs := []string{"-" + tagName + "<CreateDate", "-d", dstDir + "\\%Y.%m.%d\\%%f%%-c.%%e", src}
		execExifTool(exifToolArgs)
	},
}

func init() {
	importCmd.AddCommand(photosCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// photosCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// photosCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
