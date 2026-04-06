package clean

import (
	"github.com/redrathnure/media-tool/core/tools"
	"github.com/spf13/cobra"

	root "github.com/redrathnure/media-tool/cmd"
)

// namesCmd represents the fixNames command
var namesCmd = &cobra.Command{
	Use:   "names [files]",
	Short: "Normalize image names and remove -copy suffix",
	Long: `Renaming files with the -copy suffix to shorten variations. 
	files argument may be dir (process all files) or wildcards file names (process only matched files)`,
	Args: cobra.RangeArgs(1, 1),
	Run:  runNames,
}

func runNames(cmd *cobra.Command, args []string) {
	tools.PrintCommandArgs(cmd, args, log)

	files := tools.ExtractPath(args, 0, ".")
	log.Infof("files to process: '%s'", files)

	log.Infof("recursively: %v", recursively)

	log.Infof("dry ryn: %v", dryRun)

	exifTool := tools.GetExifTool()

	imgArgs := exifTool.NewArgs()
	tagName := "filename"

	if dryRun {
		tagName = "testname"
	}
	imgArgs.ChangeTag(tagName, "${filename;s/ - Copy/%-c/gi;s/ Copy/%-c/gi}")

	//Images and video
	//imgArgs.forImages()
	//imgArgs.forVideoMp4()

	imgArgs.Recursively(recursively)

	imgArgs.Src(files)

	//tagName = "testname" should avoid real renaming
	exifTool.Exec(root.Context.Verbose)
}

func init() {
	cleanCmd.AddCommand(namesCmd)
}
