package clean

import (
	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use:   "all [dir_or_files]",
	Short: "Cleanup image metadata and file names",
	Long: `It's shortcut for metadata + names commands  
	'dir_or_files' argument may be dir (to process all files in it) or wildcards file names (process only matched files).
	A current dir ('.' value) will be used by default.`,
	Args:    cobra.RangeArgs(0, 1),
	Aliases: []string{"cleanall", "clean_all"},
	Run:     runCleanAll,
}

func runCleanAll(cmd *cobra.Command, args []string) {
	runMetadata(cmd, args)
	runNames(cmd, args)
}

func init() {
	cleanCmd.AddCommand(allCmd)

	allCmd.Flags().BoolVarP(&includingLocation, "includingLocation", "l", false, "Remove GPS data too")
	allCmd.Flags().BoolVarP(&includingVendor, "includingVendor", "s", true, "Remove vendor specific tags")
	allCmd.Flags().BoolVarP(&includingCamera, "includingCamera", "p", false, "Remove photo/video camera info too")
}
