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

var recursively bool

// fixDatesCmd represents the fixDates command
var fixDatesCmd = &cobra.Command{
	Use:   "fixDates [files]",
	Short: "Fix Exif/QuickTime dates",
	Long: `Reads dates from file name and put into Exif and QuickTime metadata attributes. 
	files argument may be dir (process all files) or wildcards file names (process only matched files)`,
	Args: cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		printCommandArgs(cmd, args)

		files := extractPath(args, 0, ".")
		log.Infof("files to process: '%s'", files)

		log.Infof("recursively: %v", recursively)

		exifTool := getExifTool()

		//Images and video
		imgArgs := exifTool.newArgs()
		imgArgs.changeFileDate("filename")
		imgArgs.changeExifDate("filename")
		imgArgs.changeMp4Date("filename")
		//imgArgs.forImages()
		//imgArgs.forVideoMp4()
		if recursively {
			imgArgs.recursively()
		}
		imgArgs.src(files)

		exifTool.exec()
	},
}

func init() {
	rootCmd.AddCommand(fixDatesCmd)

	fixDatesCmd.Flags().BoolVarP(&recursively, "recursively", "r", false, "also analyze child directories")
}
