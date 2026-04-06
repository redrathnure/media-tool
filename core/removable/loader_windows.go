//go:build windows
// +build windows

/*
Package cmd provides command handlers

Copyright © 2020 Maksym Medvedev <redrathnure@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package removable

import (
	mtp "github.com/redrathnure/media-tool/core/removable/mtp_windows"
)

func LoadSdPhotos(targetDir string, dryRun bool) (string, error) {
	return mtp.LoadSdPhotos(targetDir, dryRun)
}

func LoadGoProVideos(targetDir string, dryRun bool) (string, error) {
	return mtp.LoadGoProVideos(targetDir, dryRun)
}

func LoadCamVideos(targetDir string, dryRun bool) (string, error) {
	return mtp.LoadCamVideos(targetDir, dryRun)
}
