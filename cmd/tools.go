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
	"path"
	"path/filepath"
)

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

func execExifTool(agrs []string) {
	cmdArgs := append([]string{"-v0", "-progress"}, agrs...)

	cmd := exec.Command(getExifTools(), cmdArgs...)
	//TODO print command for verbose mode
	fmt.Printf("command: '%s'\n", cmd.String())

	/*out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("exec error: '%s'\n", err)
	}
	fmt.Printf("exec out:\n%s", string(out[:]))
	*/
	cmd.Stdout = os.Stdout
	//TODO print Stderr for verbose mode
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("exec error: '%s'\n", err)
	}
}

var exifTool string = ""

func getExifTools() string {
	if exifTool != "" {
		return exifTool
	}

	ex, err := os.Executable()
	if err != nil {
		fmt.Printf("Unable to use custom exiftool: '%s'\n", err)
		exifTool = "exiftool"
	} else {
		exifTool, err = filepath.Abs(path.Join(filepath.Dir(ex), "exiftool", "exiftool.exe"))
		if err != nil {
			fmt.Printf("Unable to use custom exiftool: '%s'\n", err)
			exifTool = "exiftool"
		}
		if _, err := os.Stat(exifTool); os.IsNotExist(err) {
			fmt.Printf("Unable to use custom exiftool: '%s'\n", err)
			exifTool = "exiftool"
		}
	}
	return exifTool
}

func removeDir(dirName string, removeNonEmpty bool) {
	if removeNonEmpty || checkDirEmpty(dirName) {
		os.RemoveAll(dirName)
	}
}

func checkDirEmpty(dirName string) bool {
	d, err := os.Open(dirName)
	if err != nil {
		fmt.Printf("'%v' unable to open", dirName)
		return false
	}
	defer d.Close()

	stat, err := d.Stat()
	if err != nil || !stat.IsDir() {
		return false
	}

	names, err := d.Readdirnames(-1)
	if err != nil {
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
