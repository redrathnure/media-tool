/*
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

		srcDir, err := filepath.Abs(args[0])
		if err != nil {
			panic(err)
		}
		fmt.Printf("src dir: '%s'\n", srcDir)

		dstDir := "d:\\tmp\\test"

		if len(args) > 1 {
			dstDir, err = filepath.Abs(args[1])
			if err != nil {
				panic(err)
			}
		}
		fmt.Printf("dst dir: '%s'\n", dstDir)
		fmt.Printf("dry ryn: %v\n", DryRun)

		tagName := "FileName"
		if DryRun {
			tagName = "TestName"
		}
		cmdToExec := exec.Command("exiftool", "-"+tagName+"<CreateDate", "-d", dstDir+"\\%Y.%m.%d\\%%f%%-c.%%e", srcDir)
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
