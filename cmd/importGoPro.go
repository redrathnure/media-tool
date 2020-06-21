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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/redrathnure/media-tool/cmd/mtp"
)

const (
	cfgImportGoProDefaultDst = "import.gopro.default.targetDir"
)

// goproCmd represents the gopro command
var goproCmd = &cobra.Command{
	Use:   "gopro targetDir",
	Short: "Import GoPro media",
	Long: `Copy images and video from GoPro card (WPD) to disk. 
	By default creates subdirectories by dates and rename files 
	according to creation data and content type. 
	If no targetDir was specified application will try to read 
	'import.gopro.default.targetDir' configuration property`,
	Args:    cobra.RangeArgs(0, 1),
	Aliases: []string{"GoPro"},
	Run: func(cmd *cobra.Command, args []string) {
		printCommandArgs(cmd, args)

		log.Infof("src: 'GoPro' media")

		dstDir := extractPath(args, 0, "")
		if dstDir == "" {
			log.Infof("No args for targetDir was specified. Reading '%s' configuration", cfgImportGoProDefaultDst)
			dstDir = viper.GetString(cfgImportGoProDefaultDst)
			if dstDir == "" {
				log.Errorf("No target dir was specified")
				os.Exit(1)
			}
		}
		log.Infof("dst: '%s'", dstDir)

		log.Infof("dry ryn: %v", DryRun)

		src, err := mtp.LoadFromWpd("HERO", "DCIM", dstDir, !DryRun)
		if err != nil {
			log.Errorf("Unable to copy GoPro files: %v", err)
			os.Exit(1)
		}
		defer removeDir(src, DryRun)
		log.Infof("Files were downloaded to: %v", src)

		tagName := "FileName"
		if DryRun {
			tagName = "TestName"
		}

		exifTool := getExifTool()

		//Images
		imgArgs := exifTool.newArgs()
		if !DryRun {
			imgArgs.changeFileDate("CreateDate")
		}
		imgArgs.changeTag(tagName, "CreateDate")
		imgArgs.forDateFormat(dstDir + "\\%Y.%m.%d\\src\\IMG_%Y%m%d_%H%M%S%%-c.%%e")
		imgArgs.forImages()
		imgArgs.recursively()
		imgArgs.src(src)

		exifTool.exec(imgArgs)

		//Video
		videoArgs := exifTool.newArgs()
		if !DryRun {
			imgArgs.changeFileDate("CreateDate")
		}
		videoArgs.changeTag(tagName, "CreateDate")
		videoArgs.forDateFormat(dstDir + "\\%Y.%m.%d\\src\\VID_%Y%m%d_%H%M%S%%-c.%%e")
		videoArgs.forVideoMp4()
		videoArgs.recursively()
		videoArgs.src(src)

		exifTool.exec(videoArgs)
	},
}

func init() {
	importCmd.AddCommand(goproCmd)

	viper.SetDefault(cfgImportGoProDefaultDst, "")
}
