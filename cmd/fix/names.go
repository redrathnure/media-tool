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
	Short: "Fix file names (excluding WhatsApp and '*_x265' files)",
	Long: `Fix 'IMG_' and 'VID_' name prefixes for photo and video files (excluding WhatsApp and '*_x265' files). 
	May be useful to correct weird Pixel/GCam file names.  
	'dir_or_files' argument may be dir (to process all files in it) or wildcards file names (process only matched files).
	A current dir ('.' value) will be used by default.`,
	Args:    cobra.RangeArgs(0, 1),
	Aliases: []string{"fixNames"},
	Run:     runFixNames,
}

func runFixNames(cmd *cobra.Command, args []string) {
	tools.PrintCommandArgs(cmd, args, log)
	log.Infof("Fix file names (excluding WhatsApp and '*_x265' files)")

	files := tools.ExtractPath(args, 0, ".")
	log.Infof("files to process: '%s'", files)

	log.Infof("recursively: %v", recursively)

	bu := tools.NewBackuper(getConf().GetBackupLocation(), dryRun, recursively, cmd.CommandPath())
	defer bu.CleanupWorkDir(files)

	fixFileNames(files, recursively, dryRun, bu.ShouldExifDeleteOriginal())
}

func fixFileNames(files string, recursively, dryRun, shouldExifDeleteOriginal bool) {
	exifTool := tools.GetExifTool()
	exifTool.DeleteOriginals(shouldExifDeleteOriginal)

	tagName := exifTool.GetFileNameTag(dryRun)

	//Images
	log.Infof("Processing image files (excluding WhatsApp ones)...")
	imgArgs := exifTool.NewArgs()
	imgArgs.ForImages()

	imgArgs.CopyTag(tagName, "CreateDate")
	imgArgs.ForDateFormat(imgFileNameTemplate)

	imgArgs.ExcludeWhatsAppFiles()
	imgArgs.ExcludeTranscodedVideoFiles()

	imgArgs.Recursively(recursively)
	imgArgs.Src(files)

	exifTool.Exec(root.Context.Verbose)

	//Video
	log.Infof("Processing video files (excluding WhatsApp and transcoded ones)...")
	vidArgs := exifTool.NewArgs()
	vidArgs.ForVideoMp4()

	vidArgs.CopyTag(tagName, "CreateDate")
	vidArgs.ForDateFormat(vidFileNameTemplate)

	vidArgs.ExcludeWhatsAppFiles()
	vidArgs.ExcludeTranscodedVideoFiles()

	vidArgs.Recursively(recursively)
	vidArgs.Src(files)

	exifTool.Exec(root.Context.Verbose)
}

func init() {
	fixCmd.AddCommand(namesCmd)
}
