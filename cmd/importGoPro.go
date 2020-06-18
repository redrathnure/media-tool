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
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/redrathnure/media-tool/cmd/mtp"
)

// goproCmd represents the gopro command
var goproCmd = &cobra.Command{
	Use:   "gopro targetDir",
	Short: "Import GoPro media",
	Long: `Copy images and video from GoPro card (WPD) to disk. 
	By default creates subdirectories by dates and rename files 
	according to creation data and content type.`,
	Args:    cobra.RangeArgs(1, 1),
	Aliases: []string{"GoPro"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gopro called " + strings.Join(args, " "))

		fmt.Printf("src: 'GoPro' media\n")

		dstDir := extractPath(args, 0, ".")
		fmt.Printf("dst: '%s'\n", dstDir)

		fmt.Printf("dry ryn: %v\n", DryRun)

		src, err := mtp.LoadFromWpd("HERO", "DCIM", dstDir, !DryRun)
		if err != nil {
			panic(err)
		}
		defer removeDir(src, DryRun)
		fmt.Printf("Files were downloaded to: %v\n", src)

		tagName := "FileName"
		if DryRun {
			tagName = "TestName"
		}

		//Images
		//TODO exclude file date for dry run
		exifToolArgs := []string{"-FileDate<CreateDate", "-" + tagName + "<CreateDate", "-d", dstDir + "\\%Y.%m.%d\\src\\IMG_%Y%m%d_%H%M%S%%-c.%%e", "-ext", "jpg", "-r", src}
		execExifTool(exifToolArgs)

		//Video
		//TODO exclude file date for dry run
		exifToolArgs = []string{"-FileDate<CreateDate", "-" + tagName + "<CreateDate", "-d", dstDir + "\\%Y.%m.%d\\src\\VID_%Y%m%d_%H%M%S%%-c.%%e", "-ext", "mp4", "-r", src}
		execExifTool(exifToolArgs)
	},
}

func init() {
	importCmd.AddCommand(goproCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// goproCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// goproCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
