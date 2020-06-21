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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/redrathnure/media-tool/cmd/mtp"
)

const (
	cfgImportCamVideoDefaultDst = "import.camvideo.default.targetDir"
)

// goproCmd represents the gopro command
var camVideoCmd = &cobra.Command{
	Use:   "camVideo targetDir",
	Short: "Import media from Panasonic camcoder",
	Long: `Copy video from Panasonic camcoder (WPD) to disk. 
	By default creates subdirectories by dates and rename files 
	according to creation data and content type.
	If no targetDir was specified application will try to read 
	'import.camvideo.default.targetDir' configuration property`,
	Args:    cobra.RangeArgs(0, 1),
	Aliases: []string{"camvideo", "CamVideo"},
	Run: func(cmd *cobra.Command, args []string) {
		printCommandArgs(cmd, args)

		log.Infof("src: 'CamSD' media")

		dstDir := extractPath(args, 0, "")
		if dstDir == "" {
			log.Infof("No args for targetDir was specified. Reading '%s' configuration", cfgImportCamVideoDefaultDst)
			dstDir = viper.GetString(cfgImportCamVideoDefaultDst)
			if dstDir == "" {
				log.Errorf("No target dir was specified")
				os.Exit(1)
			}
		}
		log.Infof("dst: '%s'", dstDir)

		log.Infof("dry ryn: %v", DryRun)

		src, err := mtp.LoadFromWpd("CAM", path.Join("PRIVATE", "AVCHD", "BDMV", "STREAM"), dstDir, !DryRun)
		if err != nil {
			log.Errorf("Unable to copy camcoder files: %v", err)
			os.Exit(1)
		}
		defer removeDir(src, DryRun)
		log.Infof("Files were downloaded to: %v", src)

		tagName := "FileName"
		if DryRun {
			tagName = "TestName"
		}

		//Video
		exifTool := getExifTool()

		//Images
		videoArgs := exifTool.newArgs()
		if !DryRun {
			videoArgs.changeFileDate("DateTimeOriginal")
		}
		videoArgs.changeTag(tagName, "DateTimeOriginal")
		videoArgs.forDateFormat(dstDir + "\\%Y.%m.%d\\VID_%Y%m%d_%H%M%S%%-c.%%e")
		videoArgs.forVideoAvchd()
		videoArgs.recursively()
		videoArgs.src(src)

		exifTool.exec(videoArgs)
	},
}

func init() {
	importCmd.AddCommand(camVideoCmd)

	viper.SetDefault(cfgImportCamVideoDefaultDst, "")
}
