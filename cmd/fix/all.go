package fix

import (
	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use:   "all [dir_or_files]",
	Short: "Fix file WhatsApp files, then dates and names for non WhatsApp files",
	Long: `It's shortcut for whatsapp + dates + names commands  
	'dir_or_files' argument may be dir (to process all files in it) or wildcards file names (process only matched files).
	A current dir ('.' value) will be used by default.`,
	Args:    cobra.RangeArgs(0, 1),
	Aliases: []string{"fixNames"},
	Run:     runFixAll,
}

func runFixAll(cmd *cobra.Command, args []string) {
	runFixWhatsAppFiles(cmd, args)
	runFixDates(cmd, args)
	runFixNames(cmd, args)
}

func init() {
	fixCmd.AddCommand(allCmd)
}
