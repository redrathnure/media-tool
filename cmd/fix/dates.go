package fix

import (
	root "github.com/redrathnure/media-tool/cmd"
	"github.com/redrathnure/media-tool/core/tools"
	"github.com/spf13/cobra"
)

// datesCmd represents the fixDates command
var datesCmd = &cobra.Command{
	Use:   "dates [dir_or_files]",
	Short: "Fix Exif/QuickTime dates",
	Long: `Reads dates from file name and put into Exif and QuickTime metadata attributes. 
	'dir_or_files' argument may be dir (process all files) or wildcards file names (process only matched files).
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
