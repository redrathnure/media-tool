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
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// goproCmd represents the gopro command
var goproCmd = &cobra.Command{
	Use:   "gopro sourceDir [targetDir]",
	Short: "Import GoPro media",
	Long: `Copy images and video from GoPro card to disk. 
	By default creates subdirectories by dates, creates src directory 
	and rename files according to creation data and content type.`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gopro called " + strings.Join(args, " "))

		src := extractPath(args, 0, ".")
		fmt.Printf("src: '%s'\n", src)

		dstDir := extractPath(args, 1, src+"\\..")
		fmt.Printf("dst: '%s'\n", dstDir)

		fmt.Printf("dry ryn: %v\n", DryRun)

		tagName := "FileName"
		if DryRun {
			tagName = "TestName"
		}

		//Images
		cmdToExec := exec.Command("exiftool", "-v0", "-progress", "-FileDate<CreateDate", "-"+tagName+"<CreateDate", "-d", dstDir+"\\%Y.%m.%d\\src\\IMG_%Y%m%d_%H%M%S%%-c.%%e", "-ext", "jpg", src)
		/*fmt.Printf("command: '%s'\n", cmdToExec.String())*/

		/*out, err := cmdToExec.CombinedOutput()
		if err != nil {
			fmt.Printf("exec error: '%s'\n", err)
		}
		fmt.Printf("exec out:\n%s", string(out[:]))
		*/
		cmdToExec.Stdout = os.Stdout
		cmdToExec.Stderr = os.Stderr

		err := cmdToExec.Run()
		if err != nil {
			fmt.Printf("exec error: '%s'\n", err)
		}

		//Video
		cmdToExec = exec.Command("exiftool", "-v0", "-progress", "-FileDate<CreateDate", "-"+tagName+"<CreateDate", "-d", dstDir+"\\%Y.%m.%d\\src\\VID_%Y%m%d_%H%M%S%%-c.%%e", "-ext", "mp4", src)
		/*fmt.Printf("command: '%s'\n", cmdToExec.String())*/

		/*out, err = cmdToExec.CombinedOutput()
		if err != nil {
			fmt.Printf("exec error: '%s'\n", err)
		}
		fmt.Printf("exec out:\n%s", string(out[:]))*/

		cmdToExec.Stdout = os.Stdout
		cmdToExec.Stderr = os.Stderr

		err = cmdToExec.Run()
		if err != nil {
			fmt.Printf("exec error: '%s'\n", err)
		}
	},
}

func init() {
	importCmd.AddCommand(goproCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// goproCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// goproCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}