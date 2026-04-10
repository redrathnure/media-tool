package tools

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/op/go-logging"
)

var specialPaths = map[string]string{}

func extractAbsPath(args []string, argPosition int, defaultValue string) string {
	if len(args) > argPosition {
		return getAbsPath(args[argPosition])
	}

	return getAbsPath(defaultValue)
}

func ExtractPath(args []string, argPosition int, defaultValue string) string {
	if len(args) > argPosition {
		return args[argPosition]
	}

	return defaultValue
}

func getAbsPath(path string) string {
	result, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return result
}

func RemoveDir(dirName string, removeNonEmpty bool) {
	if removeNonEmpty || checkDirEmpty(dirName) {
		os.RemoveAll(dirName)
	}
}

func RemoveFiles(baseDir string, glob string) {
	files, err := filepath.Glob(path.Join(baseDir, glob))
	if err != nil {
		log.Warningf("Unable to open scan '%v' with '%v' pattern", baseDir, glob)
	}
	for _, file := range files {

		err := os.Remove(file)
		if err != nil {
			log.Warningf("'%v' unable to remove", file)
		}
	}
}

func checkDirEmpty(dirName string) bool {
	d, err := os.Open(dirName)
	if err != nil {
		log.Debugf("'%v' unable to open", dirName)
		return false
	}
	defer d.Close()

	stat, err := d.Stat()
	if err != nil || !stat.IsDir() {
		log.Debugf("'%v' is file and cannot be deleted", dirName)
		return false
	}

	names, err := d.Readdirnames(-1)
	if err != nil {
		log.Debugf("Unable to list '%v' children", dirName)
		return false
	}

	for _, name := range names {
		childIsEmpty := checkDirEmpty(path.Join(dirName, name))
		if !childIsEmpty {
			return false
		}
	}

	return true
}

func PrintCommandArgs(cmd *cobra.Command, args []string, logToUse *logging.Logger) {
	logToUse.Debugf("%s called with '%v' args", cmd.CommandPath(), strings.Join(args, " "))
}

func initSpecialPaths() {
	if len(specialPaths) == 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Warningf("Unable to resolve '~' path value. A current dir will be used instead.")
			home = "."
		}

		specialPaths["~"] = home

		exePath, err := os.Executable()
		if err != nil {
			log.Warningf("Unable to resolve '$APP_DIR' path value. A current dir will be used instead.")
			exePath = "."
		}
		exePath = filepath.Dir(exePath)
		specialPaths["$APP_DIR"] = exePath
	}

}

func ExpandPath(path string) string {
	initSpecialPaths()

	for key, value := range specialPaths {
		path = strings.ReplaceAll(path, key, value)
	}

	return filepath.Clean(os.ExpandEnv(path))
}

func HumanizeFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(size)/float64(div), "KMGTPE"[exp])
}
