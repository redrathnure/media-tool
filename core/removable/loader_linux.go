//go:build !windows
// +build !windows

package removable

import mtp "github.com/redrathnure/media-tool/core/removable/mtp_linux"

func LoadSdPhotos(targetDir string, dryRun bool) (string, error) {
	return mtp.LoadSdPhotos(targetDir, dryRun)
}

func LoadGoProVideos(targetDir string, dryRun bool) (string, error) {
	return mtp.LoadGoProVideos(targetDir, dryRun)
}

func LoadCamVideos(targetDir string, dryRun bool) (string, error) {
	return mtp.LoadCamVideos(targetDir, dryRun)
}
