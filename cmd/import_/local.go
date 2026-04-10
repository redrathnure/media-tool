package import_

import (
	"path"

	root "github.com/redrathnure/media-tool/cmd"
	"github.com/redrathnure/media-tool/core/tools"
	"github.com/spf13/cobra"
)

var localSourceSubFolder bool
var localRename bool
var localDateFormat string

var localCmd = &cobra.Command{
	Use:   "local sourceDir [targetDir]",
	Short: "Import media from local directory",
	Long: `Copy images and video from directory to disk. 
	By default creates subdirectories by dates and  keep original file name.
	Combination of -f . -r flags and same src and dst dirs may be used to corrent file names and creation date`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		tools.PrintCommandArgs(cmd, args, log)

		src := tools.ExtractPath(args, 0, ".")
		log.Infof("src: '%s'", src)

		dstDir := tools.ExtractPath(args, 1, path.Join(src, ".."))
		dstDir = tools.ExpandPath(dstDir)
		log.Infof("dst: '%s'", dstDir)

		log.Infof("dry ryn: %v", dryRun)

		dstSubDir := localDateFormat
		if localSourceSubFolder {
			dstSubDir = path.Join(dstSubDir, "src")
		}
		log.Infof("dst dir format: '%s'", dstSubDir)
		log.Infof("renaming: %v", localRename)

		exifTool := tools.GetExifTool()

		tagName := exifTool.GetFileNameTag(dryRun)

		imgFileName := "%%f%%-c.%%e"
		vidFileName := "%%f%%-c.%%e"
		if localRename {
			imgFileName = "IMG_%Y%m%d_%H%M%S%%-c.%%e"
			vidFileName = "VID_%Y%m%d_%H%M%S%%-c.%%e"
		}

		//Images
		imgArgs := exifTool.NewArgs()
		if !dryRun {
			imgArgs.ChangeFileDate("CreateDate")
		}
		imgArgs.CopyTag(tagName, "CreateDate")
		imgArgs.ForDateFormat(path.Join(dstDir, dstSubDir, imgFileName))
		imgArgs.ForImages()
		imgArgs.Recursively(true)
		imgArgs.Src(src)

		exifTool.Exec(root.Context.Verbose)

		//Video
		vidArgs := exifTool.NewArgs()
		if !dryRun {
			imgArgs.ChangeFileDate("CreateDate")
		}
		vidArgs.CopyTag(tagName, "CreateDate")
		vidArgs.ForDateFormat(path.Join(dstDir, dstSubDir, vidFileName))
		vidArgs.ForVideoMp4()
		vidArgs.Recursively(true)
		vidArgs.Src(src)

		exifTool.Exec(root.Context.Verbose)
	},
}

func init() {
	importCmd.AddCommand(localCmd)

	localCmd.Flags().BoolVarP(&localSourceSubFolder, "sourceSubDir", "s", false, "Use '\\$DATE\\src' subdir instead '\\$DATE'")
	localCmd.Flags().StringVarP(&localDateFormat, "dateFormat", "f", "%Y.%m.%d", "Date format (%Y.%m.%d by default)")
	localCmd.Flags().BoolVarP(&localRename, "rename", "r", false, "Set to rename files using 'IMG_$DATE_$TIME' and 'VID_$DATE_$TIME' patterns")
}
