package import_

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/redrathnure/media-tool/core/removable"

	root "github.com/redrathnure/media-tool/cmd"
	"github.com/redrathnure/media-tool/core/tools"
)

var sdPhotosCmd = &cobra.Command{
	Use:   "sdphotos targetDir",
	Short: "Import photos from SD card(s)",
	Long: `Copy images and video from SD card(s)to disk. 
	By default creates subdirectories by dates and rename files 
	according to creation data and content type.
	If no targetDir was specified application will try to read 
	'import.sdPhotos.default.targetDir' configuration property`,
	Args:    cobra.RangeArgs(0, 1),
	Aliases: []string{"sd", "sdPhotos"},
	Run: func(cmd *cobra.Command, args []string) {
		tools.PrintCommandArgs(cmd, args, log)

		log.Infof("src: USB flash with 'DCIM' dir")

		dstDir := tools.ExtractPath(args, 0, "")
		if dstDir == "" {
			dstDir = getConf().GetImportSdPhotosDefaultDst()
			log.Infof("No args for targetDir was specified. Using '%s' from configuration", dstDir)

			if dstDir == "" {
				log.Errorf("No target dir was specified")
				os.Exit(1)
			}
		}
		log.Infof("dst: '%s'", dstDir)

		log.Infof("dry ryn: %v", dryRun)

		src, err := removable.LoadSdPhotos(dstDir, dryRun)
		if err != nil {
			log.Errorf("Unable to copy photos files: %v", err)
			os.Exit(1)
		}
		defer tools.RemoveDir(src, dryRun)
		log.Infof("Files were downloaded to: %v. Moving to target folder...", src)

		tagName := "FileName"
		if dryRun {
			tagName = "TestName"
		}

		exifTool := tools.GetExifTool()

		//Images and video
		imgArgs := exifTool.NewArgs()
		if !dryRun {
			imgArgs.ChangeFileDate("CreateDate")
		}
		imgArgs.ChangeTag(tagName, "CreateDate")
		imgArgs.ForDateFormat(dstDir + "\\%Y.%m.%d\\%%f%%-c.%%e")
		imgArgs.ForImages()
		imgArgs.ForVideoMp4()
		imgArgs.Recursively(true)
		imgArgs.Src(src)

		exifTool.Exec(root.Context.Verbose)
	},
}

func init() {
	importCmd.AddCommand(sdPhotosCmd)
}
