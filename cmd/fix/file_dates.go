package fix

import (
	root "github.com/redrathnure/media-tool/cmd"
	"github.com/redrathnure/media-tool/core/tools"
	"github.com/spf13/cobra"
)

var fileDatesCmd = &cobra.Command{
	Use:   "file_dates [dir_or_files]",
	Short: "Fix file dates",
	Long: `Fix file creation and last update date attribute for photo and video files.
	Similar to 'fix dates' command but takes date from Exif data instead file name.
	May be useful to restore original file attributes after manual copy of media library or correct time shift for video files.  
	'dir_or_files' argument may be dir (to process all files in it) or wildcards file names (process only matched files).
	A current dir ('.' value) will be used by default.`,
	Args:    cobra.RangeArgs(0, 1),
	Aliases: []string{"fileDates", "fixFileDates"},
	Run:     runFixFileDates,
}

func runFixFileDates(cmd *cobra.Command, args []string) {
	tools.PrintCommandArgs(cmd, args, log)
	log.Infof("Fix file dates")

	files := tools.ExtractPath(args, 0, ".")
	log.Infof("files to process: '%s'", files)

	log.Infof("recursively: %v", recursively)

	bu := tools.NewBackuper(getConf().GetBackupLocation(), dryRun, recursively, cmd.CommandPath())
	defer bu.CleanupWorkDir(files)

	exifTool := tools.GetExifTool()
	exifTool.DeleteOriginals(bu.ShouldExifDeleteOriginal())

	log.Infof("Processing image and video files...")

	//Images + Video
	imgArgs := exifTool.NewArgs()
	// For all supported files
	//imgArgs.ForImages()
	//imgArgs.ForVideoMp4()

	if dryRun {
		imgArgs.CopyTag("WriteNothing", "CreateDate")
	} else {
		imgArgs.ChangeFileDate("CreateDate")
	}

	imgArgs.Recursively(recursively)
	imgArgs.Src(files)

	exifTool.Exec(root.Context.Verbose)
}

func init() {
	fixCmd.AddCommand(fileDatesCmd)
}
