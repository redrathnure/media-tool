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
package cmd

import (
	"github.com/spf13/cobra"
)

var includingLocation bool
var includingVendor bool
var includingCamera bool

// cleanMetadataCmd represents the fixNames command
var cleanMetadataCmd = &cobra.Command{
	Use:   "metadata [files]",
	Short: "Cleanup image metadata",
	Long: `Remove vendor metadata from media files. 
	files argument may be dir (process all files) or wildcards file names (process only matched files)`,
	Args: cobra.RangeArgs(1, 1),
	Run:  runCleanMetadata,
}

func runCleanMetadata(cmd *cobra.Command, args []string) {
	printCommandArgs(cmd, args)

	files := extractPath(args, 0, ".")
	log.Infof("files to process: '%s'", files)

	log.Infof("recursively: %v", recursively)
	log.Infof("includingLocation: %v", includingLocation)
	log.Infof("includingVendor: %v", includingVendor)

	log.Infof("dry ryn: %v", DryRun)

	exifTool := getExifTool()

	imgArgs := exifTool.newArgs()
	if includingLocation {
		imgArgs.cleanLocationTags()
	}
	if includingVendor {
		imgArgs.cleanVendorTags()
	}
	if includingCamera {
		imgArgs.cleanCameraTags()
	}

	//Images and video
	//imgArgs.forImages()
	//imgArgs.forVideoMp4()

	if recursively {
		imgArgs.recursively()
	}

	imgArgs.src(files)

	if !DryRun {
		exifTool.exec()
	}
}

func init() {
	cleanCmd.AddCommand(cleanMetadataCmd)

	cleanMetadataCmd.Flags().BoolVarP(&includingLocation, "includingLocation", "l", false, "Remove GPS data too")
	cleanMetadataCmd.Flags().BoolVarP(&includingVendor, "includingVendor", "s", true, "Remove vendor specific tags")
	cleanMetadataCmd.Flags().BoolVarP(&includingCamera, "includingCamera", "p", false, "Remove photo/video camera info too")
}
