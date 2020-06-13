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
	"os/exec"
	"strings"

	"path/filepath"

	"github.com/spf13/cobra"
)

// photosCmd represents the photos command
var photosCmd = &cobra.Command{
	Use:   "photos sourceDir [targetDir]",
	Short: "Import photos from SD card",
	Long:  `Copy images from SD card to disk. By default creates subdirectories by dates.`,
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("photos called " + strings.Join(args, " "))

		src := extractPath(args, 0, ".")
		fmt.Printf("src: '%s'\n", src)

		dstDir := extractPath(args, 1, "d:\\tmp\\test")
		fmt.Printf("dst: '%s'\n", dstDir)

		fmt.Printf("dry ryn: %v\n", DryRun)

		tagName := "FileName"
		if DryRun {
			tagName = "TestName"
		}
		cmdToExec := exec.Command("exiftool", "-"+tagName+"<CreateDate", "-d", dstDir+"\\%Y.%m.%d\\%%f%%-c.%%e", src)
		fmt.Printf("command: '%s'\n", cmdToExec.String())

		out, err := cmdToExec.CombinedOutput()
		if err != nil {
			fmt.Printf("exec error: '%s'\n", err)
		}
		fmt.Printf("exec out:\n%s", string(out[:]))
	},
}

func init() {
	importCmd.AddCommand(photosCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// photosCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// photosCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func extractAbsPath(args []string, argPosition int, defaultValue string) string {
	if len(args) > argPosition {
		return getAbsPath(args[argPosition])
	}

	return getAbsPath(defaultValue)
}

func extractPath(args []string, argPosition int, defaultValue string) string {
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
