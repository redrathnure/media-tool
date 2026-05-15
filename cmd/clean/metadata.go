package clean

import (
	root "github.com/redrathnure/media-tool/cmd"
	"github.com/redrathnure/media-tool/core/tools"
	"github.com/spf13/cobra"
)

var includingLocation bool
var includingVendor bool
var includingCamera bool

// metadataCmd represents the fixNames command
var metadataCmd = &cobra.Command{
	Use:   "metadata [files]",
	Short: "Cleanup image metadata",
	Long: `Remove vendor metadata from media files. 
	files argument may be dir (process all files) or wildcards file names (process only matched files)`,
	Args: cobra.RangeArgs(1, 1),
	Run:  runMetadata,
}

func runMetadata(cmd *cobra.Command, args []string) {
	tools.PrintCommandArgs(cmd, args, log)
	log.Infof("Clean metadata")

	files := tools.ExtractPath(args, 0, ".")
	log.Infof("files to process: '%s'", files)

	log.Infof("recursively: %v", recursively)
	log.Infof("includingLocation: %v", includingLocation)
	log.Infof("includingVendor: %v", includingVendor)

	log.Infof("dry ryn: %v", dryRun)

	bu := tools.NewBackuper(getConf().GetBackupLocation(), dryRun, recursively, cmd.CommandPath())
	defer bu.CleanupWorkDir(files)

	exifTool := tools.GetExifTool()
	exifTool.DeleteOriginals(bu.ShouldExifDeleteOriginal())

	imgArgs := exifTool.NewArgs()
	if includingLocation {
		imgArgs.CleanLocationTags()
	}
	if includingVendor {
		imgArgs.CleanVendorTags()
	}
	if includingCamera {
		imgArgs.CleanCameraTags()
	}

	//Images and video
	//imgArgs.forImages()
	//imgArgs.forVideoMp4()

	imgArgs.Recursively(recursively)

	imgArgs.Src(files)

	if !dryRun {
		exifTool.Exec(root.Context.Verbose)
	}
}

func init() {
	cleanCmd.AddCommand(metadataCmd)

	metadataCmd.Flags().BoolVarP(&includingLocation, "includingLocation", "l", false, "Remove GPS data too")
	metadataCmd.Flags().BoolVarP(&includingVendor, "includingVendor", "s", true, "Remove vendor specific tags")
	metadataCmd.Flags().BoolVarP(&includingCamera, "includingCamera", "p", false, "Remove photo/video camera info too")
}
