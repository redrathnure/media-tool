package import_

import (
	"os"
	"path"

	"github.com/spf13/cobra"

	"github.com/redrathnure/media-tool/core/removable"

	root "github.com/redrathnure/media-tool/cmd"
	"github.com/redrathnure/media-tool/core/tools"
)

var keepFileNames bool

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
		dstDir = tools.ExpandPath(dstDir)
		log.Infof("dst: '%s'", dstDir)

		log.Infof("Keep original file names: %v", keepFileNames)

		log.Infof("dry ryn: %v", dryRun)

		src, err := removable.LoadSdPhotos(dstDir, dryRun)
		if err != nil {
			log.Errorf("Unable to copy photos files: %v", err)
			os.Exit(1)
		}
		defer tools.RemoveDir(src, dryRun)
		log.Infof("Files were downloaded to: %v. Moving to target folder...", src)

		sdPhotosMoveToDst(src, dstDir, keepFileNames, dryRun)
	},
}

func sdPhotosMoveToDst(src, dstDir string, keepOriginFileNames, dryRun bool) {
	imgFileName := "%%f%%-c.%%e"
	vidFileName := "%%f%%-c.%%e"
	if !keepOriginFileNames {
		imgFileName = "IMG_%Y%m%d_%H%M%S%%-c.%%e"
		vidFileName = "VID_%Y%m%d_%H%M%S%%-c.%%e"
	}

	exifTool := tools.GetExifTool()

	tagName := exifTool.GetFileNameTag(dryRun)

	if keepOriginFileNames {
		//Workaround for weird ExifTool file modification handling, when time component are 00:00:00
		log.Infof("Preparing file dates...")
		metaArgs := exifTool.NewArgs()
		if !dryRun {
			metaArgs.ChangeFileDate("CreateDate")
		}
		metaArgs.ForImages()
		metaArgs.ForVideoMp4()
		metaArgs.Recursively(true)
		metaArgs.Src(src)

		exifTool.Exec(root.Context.Verbose)
	}

	//Images
	log.Infof("Processing image files...")
	imgArgs := exifTool.NewArgs()
	if !dryRun && !keepOriginFileNames {
		imgArgs.ChangeFileDate("CreateDate")
	}
	imgArgs.CopyTag(tagName, "CreateDate")
	imgArgs.ForDateFormat(path.Join(dstDir, "%Y.%m.%d", imgFileName))
	imgArgs.ForImages()
	imgArgs.Recursively(true)
	imgArgs.Src(src)

	exifTool.Exec(root.Context.Verbose)

	// Video
	log.Infof("Processing video files...")
	vidArgs := exifTool.NewArgs()
	if !dryRun && !keepOriginFileNames {
		vidArgs.ChangeFileDate("CreateDate")
	}
	vidArgs.CopyTag(tagName, "CreateDate")
	vidArgs.ForDateFormat(path.Join(dstDir, "%Y.%m.%d", vidFileName))
	vidArgs.ForVideoMp4()
	vidArgs.Recursively(true)
	vidArgs.Src(src)

	exifTool.Exec(root.Context.Verbose)
}

func init() {
	importCmd.AddCommand(sdPhotosCmd)

	sdPhotosCmd.Flags().BoolVarP(&keepFileNames, "keep-names", "k", false, "Keep original file names. By default files will be renamed using 'IMG_$DATE_$TIME' and 'VID_$DATE_$TIME' patterns")
}
