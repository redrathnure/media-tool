//go:build !windows
// +build !windows

package removable

import "fmt"

func LoadSdPhotos(targetDir string, dryRun bool) (string, error) {
	return warnUnsupported("LoadSdPhotos", targetDir)
}

func LoadGoProVideos(targetDir string, dryRun bool) (string, error) {
	return warnUnsupported("LoadGoProVideos", targetDir)
}

func LoadCamVideos(targetDir string, dryRun bool) (string, error) {
	return warnUnsupported("LoadCamVideos", targetDir)
}

func warnUnsupported(name, targetDir string) (string, error) {
	message := fmt.Sprintf("%s is not supported on this platform. MTP imports are only available on Windows. Target dir: %s", name, targetDir)
	log.Warningf(message)
	return "", fmt.Errorf("%s", message)
}
