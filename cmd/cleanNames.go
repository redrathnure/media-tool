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

// cleanNamesCmd represents the fixNames command
var cleanNamesCmd = &cobra.Command{
	Use:   "names [files]",
	Short: "Normalize image names and remove -copy suffix",
	Long: `Renaming files with the -copy suffix to shorten variations. 
	files argument may be dir (process all files) or wildcards file names (process only matched files)`,
	Args: cobra.RangeArgs(1, 1),
	Run:  runCleanNames,
}

func runCleanNames(cmd *cobra.Command, args []string) {
	printCommandArgs(cmd, args)

	files := extractPath(args, 0, ".")
	log.Infof("files to process: '%s'", files)

	log.Infof("recursively: %v", recursively)

	log.Infof("dry ryn: %v", DryRun)

	exifTool := getExifTool()

	imgArgs := exifTool.newArgs()
	tagName := "filename"

	if DryRun {
		tagName = "testname"
	}
	imgArgs.changeTag(tagName, "${filename;s/ - Copy/%-c/gi;s/ Copy/%-c/gi}")

	//Images and video
	//imgArgs.forImages()
	//imgArgs.forVideoMp4()
	if recursively {
		imgArgs.recursively()
	}

	imgArgs.src(files)

	//tagName = "testname" should avoid real renaming
	exifTool.exec()
}

func init() {
	cleanCmd.AddCommand(cleanNamesCmd)
}
