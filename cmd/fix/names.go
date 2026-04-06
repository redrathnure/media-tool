package fix

import (
	root "github.com/redrathnure/media-tool/cmd"
	"github.com/redrathnure/media-tool/core/tools"
	"github.com/spf13/cobra"
)

const (
	imgFileNameTemplate = "IMG_%Y%m%d_%H%M%S%%-c.%%e"
	vidFileNameTemplate = "VID_%Y%m%d_%H%M%S%%-c.%%e"
)

var namesCmd = &cobra.Command{
	Use:   "names [dir_or_files]",
	Short: "Fix file names",
	Long: `Fix 'IMG_' and 'VID_' file name prefixes for photo and video files. 
	May be useful to correct weird Pixel/GCam file names.  
	'dir_or_files' argument may be dir (to process all files in it) or wildcards file names (process only matched files).
	A current dir ('.' value) will be used by default.`,
	Args:    cobra.RangeArgs(0, 1),
	Aliases: []string{"fixNames"},
	Run:     runFixNames,
}

func runFixNames(cmd *cobra.Command, args []string) {
	tools.PrintCommandArgs(cmd, args, log)

	files := tools.ExtractPath(args, 0, ".")
	log.Infof("files to process: '%s'", files)

	log.Infof("recursively: %v", recursively)

	tagName := "FileName"
	if dryRun {
		tagName = "TestName"
	}

	exifTool := tools.GetExifTool()

	log.Infof("Processing image files...")

	//Images
	imgArgs := exifTool.NewArgs()
	imgArgs.ChangeTag(tagName, "CreateDate")
	imgArgs.ForDateFormat(imgFileNameTemplate)
	imgArgs.ForImages()
	imgArgs.Recursively(recursively)
	imgArgs.Src(files)

	exifTool.Exec(root.Context.Verbose)

	//Video
	log.Infof("Processing video files...")
	vidArgs := exifTool.NewArgs()
	vidArgs.ChangeTag(tagName, "CreateDate")
	vidArgs.ForDateFormat(vidFileNameTemplate)
	vidArgs.ForVideoMp4()
	vidArgs.Recursively(recursively)
	vidArgs.Src(files)

	exifTool.Exec(root.Context.Verbose)
}

func init() {
	fixCmd.AddCommand(namesCmd)
}
