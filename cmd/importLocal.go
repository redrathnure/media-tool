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
	"path"

	"github.com/spf13/cobra"
)

var localSourceSubFolder bool
var localRename bool
var localDateFormat string

// video represents the gopro command
var importLocal = &cobra.Command{
	Use:   "local sourceDir [targetDir]",
	Short: "Import media from local directory",
	Long: `Copy images and video from directory to disk. 
	By default creates subdirectories by dates and  keep original file name.
	Combination of -f . -r flags and same src and dst dirs may be used to corrent file names and creation date`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		printCommandArgs(cmd, args)

		src := extractPath(args, 0, ".")
		log.Infof("src: '%s'", src)

		dstDir := extractPath(args, 1, src+"\\..")
		log.Infof("dst: '%s'", dstDir)

		log.Infof("dry ryn: %v", DryRun)

		dstSubDir := localDateFormat
		if localSourceSubFolder {
			dstSubDir = path.Join(dstSubDir, "src")
		}
		log.Infof("dst dir format: '%s'", dstSubDir)
		log.Infof("renaming: %v", localRename)

		tagName := "FileName"
		if DryRun {
			tagName = "TestName"
		}

		imgFileName := "%%f%%-c.%%e"
		vidFileName := "%%f%%-c.%%e"
		if localRename {
			imgFileName = "IMG_%Y%m%d_%H%M%S%%-c.%%e"
			vidFileName = "VID_%Y%m%d_%H%M%S%%-c.%%e"
		}

		exifTool := getExifTool()

		//Images
		imgArgs := exifTool.newArgs()
		if !DryRun {
			imgArgs.changeFileDate("CreateDate")
		}
		imgArgs.changeTag(tagName, "CreateDate")
		imgArgs.forDateFormat(path.Join(dstDir, dstSubDir, imgFileName))
		imgArgs.forImages()
		imgArgs.recursively()
		imgArgs.src(src)

		exifTool.exec(imgArgs)

		//Video
		vidArgs := exifTool.newArgs()
		if !DryRun {
			imgArgs.changeFileDate("CreateDate")
		}
		vidArgs.changeTag(tagName, "CreateDate")
		vidArgs.forDateFormat(path.Join(dstDir, dstSubDir, vidFileName))
		vidArgs.forVideoMp4()
		vidArgs.recursively()
		vidArgs.src(src)

		exifTool.exec(vidArgs)
	},
}

func init() {
	importCmd.AddCommand(importLocal)

	importLocal.Flags().BoolVarP(&localSourceSubFolder, "sourceSubDir", "s", false, "Use '\\$DATE\\src' subdir instead '\\$DATE'")
	importLocal.Flags().StringVarP(&localDateFormat, "dateFormat", "f", "%Y.%m.%d", "Date format (%Y.%m.%d by default)")
	importLocal.Flags().BoolVarP(&localRename, "rename", "r", false, "Set to rename files using 'IMG_$DATE_$TIME' and 'VID_$DATE_$TIME' patterns")
}
