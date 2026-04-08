package fix

import (
	root "github.com/redrathnure/media-tool/cmd"
	"github.com/redrathnure/media-tool/core/tools"
	"github.com/spf13/cobra"
)

const (
	whatsAppImgFileNameTemplate = "IMG_%Y%m%d_%H%M%S_WhatsApp%%-c.%%e"
	whatsAppVidFileNameTemplate = "VID_%Y%m%d_%H%M%S_WhatsApp%%-c.%%e"
)

var whatsAppCmd = &cobra.Command{
	Use:   "whatsapp [dir_or_files]",
	Short: "Fix WhatsApp file names and Exif dates",
	Long: `Fix 'IMG_' and 'VID_' name prefixes for photo and video files downloaded from WhatsApp as well as updating Exif/QuickTime dates.
	It process only files with 'WhatsApp' marker in the name.
	May be useful to correct naming for home archive. 
	'dir_or_files' argument may be dir (to process all files in it) or wildcards file names (process only matched files).
	A current dir ('.' value) will be used by default.`,
	Args:    cobra.RangeArgs(0, 1),
	Aliases: []string{"whatsapp", "whatsapp_names", "whatsappNames", "wa", "waNames", "wa_names"},
	Run:     runFixWhatsAppFiles,
}

func runFixWhatsAppFiles(cmd *cobra.Command, args []string) {
	tools.PrintCommandArgs(cmd, args, log)
	log.Infof("Fix WhatsApp files")

	files := tools.ExtractPath(args, 0, ".")
	log.Infof("files to process: '%s'", files)

	log.Infof("recursively: %v", recursively)

	bu := tools.NewBackuper(getConf().GetBackupLocation(), dryRun, recursively, cmd.CommandPath())
	defer bu.CleanupWorkDir(files)

	exifTool := tools.GetExifTool()
	exifTool.DeleteOriginals(bu.ShouldExifDeleteOriginal())

	//Adjust metadata
	log.Infof("Processing WhatsApp image and video metadata...")
	exifArgs := exifTool.NewArgs()
	// For image and video formats
	exifArgs.ForImages()
	exifArgs.ForVideoMp4()

	if dryRun {
		exifArgs.CopyTag("WriteNothing", "filename")
	} else {
		exifArgs.ChangeFileDate("filename")
		exifArgs.ChangeExifDate("filename")
		exifArgs.ChangeMp4Date("filename")

		exifArgs.SetTag("Exif:ProcessingSoftware", "WhatsApp")
		exifArgs.SetTag("Exif:Software", "WhatsApp")
		exifArgs.SetTag("Exif:MetadataEditingSoftware", "WhatsApp")
	}

	exifArgs.IncludeWhatsAppFiles()
	//Just in case somebody decide to compress already recompressed files
	exifArgs.ExcludeTranscodedVideoFiles()

	exifArgs.Recursively(recursively)
	exifArgs.Src(files)

	exifTool.Exec(root.Context.Verbose)

	// Rename files
	tagName := exifTool.GetFileNameTag(dryRun)

	//Images
	log.Infof("Renaming WhatsApp image files...")
	imgArgs := exifTool.NewArgs()
	// For images only
	imgArgs.ForImages()

	imgArgs.CopyTag(tagName, "CreateDate")
	imgArgs.ForDateFormat(whatsAppImgFileNameTemplate)

	imgArgs.IncludeWhatsAppFiles()
	//Just in case somebody decide to compress already recompressed files
	imgArgs.ExcludeTranscodedVideoFiles()

	imgArgs.Recursively(recursively)
	imgArgs.Src(files)

	exifTool.Exec(root.Context.Verbose)

	//Video
	log.Infof("Renaming WhatsApp video files...")
	vidArgs := exifTool.NewArgs()
	// For video only
	vidArgs.ForVideoMp4()

	vidArgs.CopyTag(tagName, "CreateDate")
	vidArgs.ForDateFormat(whatsAppVidFileNameTemplate)

	vidArgs.IncludeWhatsAppFiles()
	//Just in case somebody decide to compress already recompressed files
	vidArgs.ExcludeTranscodedVideoFiles()

	vidArgs.Recursively(recursively)
	vidArgs.Src(files)

	exifTool.Exec(root.Context.Verbose)
}

func init() {
	fixCmd.AddCommand(whatsAppCmd)
}
