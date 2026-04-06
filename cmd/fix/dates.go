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
package fix

import (
	root "github.com/redrathnure/media-tool/cmd"
	"github.com/redrathnure/media-tool/core/tools"
	"github.com/spf13/cobra"
)

// datesCmd represents the fixDates command
var datesCmd = &cobra.Command{
	Use:   "dates [files]",
	Short: "Fix Exif/QuickTime dates",
	Long: `Reads dates from file name and put into Exif and QuickTime metadata attributes. 
	files argument may be dir (process all files) or wildcards file names (process only matched files)`,
	Args: cobra.RangeArgs(1, 1),
	Aliases: []string{"fixDates"},
	Run:     runFixDates,
}

func runFixDates(cmd *cobra.Command, args []string) {
	tools.PrintCommandArgs(cmd, args, log)

	files := tools.ExtractPath(args, 0, ".")
	log.Infof("files to process: '%s'", files)

	log.Infof("recursively: %v", recursively)

	exifTool := tools.GetExifTool()

	//Images and video
	imgArgs := exifTool.NewArgs()
	imgArgs.ChangeFileDate("filename")
	imgArgs.ChangeExifDate("filename")
	imgArgs.ChangeMp4Date("filename")
	//imgArgs.forImages()
	//imgArgs.forVideoMp4()
	imgArgs.Recursively(recursively)
	imgArgs.Src(files)

	exifTool.Exec(root.Context.Verbose)

	// TODO dryRun?
}

func init() {
	fixCmd.AddCommand(datesCmd)
}
