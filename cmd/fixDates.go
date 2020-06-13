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
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var recursively bool

// fixDatesCmd represents the fixDates command
var fixDatesCmd = &cobra.Command{
	Use:   "fixDates [files]",
	Short: "Fix EXIF/Quicktime dates",
	Long: `Reads dates from file name and put into EXIF and Quicktime metadata attributes. 
	files argument may be dir (proces all files) or wildcards file names (process only matched files)`,
	Args: cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fixDates called " + strings.Join(args, " "))

		files := extractFiles(args)
		fmt.Printf("files to process: '%s'\n", files)

		fmt.Printf("recursively: %v\n", recursively)

		exifToolArgs := []string{"-v2", "-ImageDate<filename", "-VideoDate<filename", "-FileDate<filename", files}

		if recursively {
			exifToolArgs = append(exifToolArgs, "-r")
		}

		execExifTool(exifToolArgs)
	},
}

func init() {
	rootCmd.AddCommand(fixDatesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fixDatesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fixDatesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	fixDatesCmd.Flags().BoolVarP(&recursively, "recursively", "r", false, "also analyze child directories")
}

func extractFiles(args []string) string {
	if len(args) > 0 {
		dstDir, err := filepath.Abs(args[0])
		if err != nil {
			panic(err)
		}
		return dstDir
	}
	return "."
}
