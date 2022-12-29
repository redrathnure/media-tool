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
	cfgImportSdPhotosDefaultDst = "import.sdPhotos.default.targetDir"
)

// sdPhotos represents the gopro command
var sdPhotos = &cobra.Command{
	Use:   "sdphotos targetDir",
	Short: "Import photos from SD card(s)",
	Long: `Copy images and video from SD card(s)to disk. 
	By default creates subdirectories by dates and rename files 
	according to creation data and content type.
	If no targetDir was specified application will try to read 
	'import.sdPhotos.default.targetDir' configuration property`,
	Args:    cobra.RangeArgs(0, 1),
	Aliases: []string{"sd", "sdPhotos"},
	Run: func(cmd *cobra.Command, args []string) {
		printCommandArgs(cmd, args)

		log.Infof("src: USB flash with 'DCIM' dir")

		dstDir := extractPath(args, 0, "")
		if dstDir == "" {
			log.Infof("No args for targetDir was specified. Reading '%s' configuration", cfgImportSdPhotosDefaultDst)
			dstDir = viper.GetString(cfgImportSdPhotosDefaultDst)
			if dstDir == "" {
				log.Errorf("No target dir was specified")
				os.Exit(1)
			}
		}
		log.Infof("dst: '%s'", dstDir)

		log.Infof("dry ryn: %v", DryRun)

		src, err := mtp.LoadFromAllWpd("DCIM", dstDir, DryRun)
		if err != nil {
			log.Errorf("Unable to copy photos files: %v", err)
			os.Exit(1)
		}
		defer removeDir(src, DryRun)
		log.Infof("Files were downloaded to: %v. Moving to target folder...", src)

		tagName := "FileName"
		if DryRun {
			tagName = "TestName"
		}

		exifTool := getExifTool()

		//Images and video
		imgArgs := exifTool.newArgs()
		if !DryRun {
			imgArgs.changeFileDate("CreateDate")
		}
		imgArgs.changeTag(tagName, "CreateDate")
		imgArgs.forDateFormat(dstDir + "\\%Y.%m.%d\\%%f%%-c.%%e")
		imgArgs.forImages()
		imgArgs.forVideoMp4()
		imgArgs.recursively()
		imgArgs.src(src)

		exifTool.exec(imgArgs)
	},
}

func init() {
	importCmd.AddCommand(sdPhotos)

	viper.SetDefault(cfgImportSdPhotosDefaultDst, "")
}
