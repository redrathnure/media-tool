package import_

import (
	"os"
	"path"

	"github.com/spf13/cobra"

	root "github.com/redrathnure/media-tool/cmd"
	"github.com/redrathnure/media-tool/core/removable"
	"github.com/redrathnure/media-tool/core/tools"
)

var goProCmd = &cobra.Command{
	Use:   "gopro targetDir",
	Short: "Import GoPro media",
	Long: `Copy images and video from GoPro card (WPD) to disk. 
	By default creates subdirectories by dates and rename files 
	according to creation data and content type. 
	If no targetDir was specified application will try to read 
	'import.gopro.default.targetDir' configuration property`,
	Args:    cobra.RangeArgs(0, 1),
	Aliases: []string{"GoPro"},
	Run: func(cmd *cobra.Command, args []string) {
		tools.PrintCommandArgs(cmd, args, log)

		log.Infof("src: 'GoPro' media")

		dstDir := tools.ExtractPath(args, 0, "")
		if dstDir == "" {
			dstDir = getConf().GetImportGoProDefaultDst()
			log.Infof("No args for targetDir was specified. Using '%s' from configuration", dstDir)

			if dstDir == "" {
				log.Errorf("No target dir was specified")
				os.Exit(1)
			}
		}
		dstDir = tools.ExpandPath(dstDir)
		log.Infof("dst: '%s'", dstDir)

		log.Infof("dry ryn: %v", dryRun)

		src, err := removable.LoadGoProVideos(dstDir, dryRun)
		if err != nil {
			log.Errorf("Unable to copy GoPro files: %v", err)
			os.Exit(1)
		}
		defer tools.RemoveDir(src, dryRun)
		log.Infof("Files were downloaded to: %v. Moving to target folder...", src)

		exifTool := tools.GetExifTool()

		tagName := exifTool.GetFileNameTag(dryRun)

		//Images
		log.Infof("Processing image files...")
		imgArgs := exifTool.NewArgs()
		if !dryRun {
			imgArgs.ChangeFileDate("CreateDate")
		}
		imgArgs.CopyTag(tagName, "CreateDate")
		imgArgs.ForDateFormat(path.Join(dstDir, "%Y.%m.%d", "src", "IMG_%Y%m%d_%H%M%S%%-c.%%e"))
		imgArgs.ForImages()
		imgArgs.Recursively(true)
		imgArgs.Src(src)

		exifTool.Exec(root.Context.Verbose)

		//Video
		log.Infof("Processing video files...")
		vidArgs := exifTool.NewArgs()
		if !dryRun {
			vidArgs.ChangeFileDate("CreateDate")
		}
		vidArgs.CopyTag(tagName, "CreateDate")
		vidArgs.ForDateFormat(path.Join(dstDir, "%Y.%m.%d", "src", "VID_%Y%m%d_%H%M%S%%-c.%%e"))
		vidArgs.ForVideoMp4()
		vidArgs.Recursively(true)
		vidArgs.Src(src)

		exifTool.Exec(root.Context.Verbose)

		//Video Preview
		log.Infof("Processing video preview files...")
		vidPreviewArgs := exifTool.NewArgs()
		if !dryRun {
			vidPreviewArgs.ChangeFileDate("CreateDate")
		}
		vidPreviewArgs.CopyTag(tagName, "CreateDate")
		vidPreviewArgs.ForDateFormat(path.Join(dstDir, "%Y.%m.%d", "src", "VID_%Y%m%d_%H%M%S%%-c.preview.mp4"))
		vidPreviewArgs.ForVideoLrv()
		vidPreviewArgs.Recursively(true)
		vidPreviewArgs.Src(src)

		exifTool.Exec(root.Context.Verbose)

		tools.RemoveFiles(src, "leinfo.sav")
		tools.RemoveFiles(src, path.Join("**", "*.THM"))
	},
}

func init() {
	importCmd.AddCommand(goProCmd)
}
