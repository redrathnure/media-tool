package import_

import (
	"os"

	"github.com/spf13/cobra"

	root "github.com/redrathnure/media-tool/cmd"
	"github.com/redrathnure/media-tool/core/removable"
	"github.com/redrathnure/media-tool/core/tools"
)

var videoCamCmd = &cobra.Command{
	Use:   "videocam targetDir",
	Short: "Import media from Panasonic camcoder",
	Long: `Copy video from Panasonic camcoder (WPD) to disk. 
	By default creates subdirectories by dates and rename files 
	according to creation data and content type.
	If no targetDir was specified application will try to read 
	'import.camvideo.default.targetDir' configuration property`,
	Args:    cobra.RangeArgs(0, 1),
	Aliases: []string{"camVideo", "CamVideo", "camvideo", "videoCam", "VideoCam"},
	Run: func(cmd *cobra.Command, args []string) {
		tools.PrintCommandArgs(cmd, args, log)

		log.Infof("src: 'CamSD' media")

		dstDir := tools.ExtractPath(args, 0, "")
		if dstDir == "" {
			dstDir = getConf().GetImportCamVideoDefaultDst()
			log.Infof("No args for targetDir was specified. Using '%s' from configuration", dstDir)

			if dstDir == "" {
				log.Errorf("No target dir was specified")
				os.Exit(1)
			}
		}
		log.Infof("dst: '%s'", dstDir)

		log.Infof("dry ryn: %v", dryRun)

		src, err := removable.LoadCamVideos(dstDir, dryRun)
		if err != nil {
			log.Errorf("Unable to copy camcoder files: %v", err)
			os.Exit(1)
		}
		defer tools.RemoveDir(src, dryRun)
		log.Infof("Files were downloaded to: %v. Moving to target folder...", src)

		exifTool := tools.GetExifTool()

		tagName := exifTool.GetFileNameTag(dryRun)

		//Video
		videoArgs := exifTool.NewArgs()
		if !dryRun {
			videoArgs.ChangeFileDate("DateTimeOriginal")
		}
		videoArgs.ChangeTag(tagName, "DateTimeOriginal")
		videoArgs.ForDateFormat(dstDir + "\\%Y.%m.%d\\VID_%Y%m%d_%H%M%S%%-c.%%e")
		videoArgs.ForVideoAvchd()
		videoArgs.Recursively(true)
		videoArgs.Src(src)

		exifTool.Exec(root.Context.Verbose)
	},
}

func init() {
	importCmd.AddCommand(videoCamCmd)
}
