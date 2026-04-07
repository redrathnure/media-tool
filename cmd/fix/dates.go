package fix

import (
	root "github.com/redrathnure/media-tool/cmd"
	"github.com/redrathnure/media-tool/core/tools"
	"github.com/spf13/cobra"
)

var datesCmd = &cobra.Command{
	Use:   "dates [dir_or_files]",
	Short: "Fix Exif/QuickTime and file creation dates",
	Long: `Reads dates from file name and put into Exif and QuickTime metadata attributes s well as updating file creation/modification dates. 
	'dir_or_files' argument may be dir (to process all files in it) or wildcards file names (process only matched files).
	A current dir ('.' value) will be used by default.`,
	Args:    cobra.RangeArgs(0, 1),
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
	// For all supported files
	//imgArgs.ForImages()
	//imgArgs.ForVideoMp4()

	if dryRun {
		imgArgs.CopyTag("WriteNothing", "filename")
	} else {
		imgArgs.ChangeFileDate("filename")
		imgArgs.ChangeExifDate("filename")
		imgArgs.ChangeMp4Date("filename")
	}

	imgArgs.Recursively(recursively)
	imgArgs.Src(files)

	exifTool.Exec(root.Context.Verbose)
}

func init() {
	fixCmd.AddCommand(datesCmd)
}
