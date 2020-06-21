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
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

type extifToolWrapper struct {
	cmd         string
	defaultArgs []string
}

var exifToolObj *extifToolWrapper

func newExtifTool() *extifToolWrapper {
	result := extifToolWrapper{
		cmd:         "exiftool",
		defaultArgs: []string{"-v0", "-progress"},
	}
	result.initCmd()
	return &result
}

func (tool extifToolWrapper) initCmd() {
	ex, err := os.Executable()
	if err != nil {
		log.Infof("Unable to find custom exiftool: '%s'. Trying to use '%s' from $PATH", err, tool.cmd)
		return
	}
	customePath, err := filepath.Abs(path.Join(filepath.Dir(ex), "exiftool", "exiftool.exe"))
	if err != nil {
		log.Infof("Unable to find custom exiftool: '%s'. Trying to use '%s' from $PATH", err, tool.cmd)
		return
	}
	if _, err := os.Stat(customePath); os.IsNotExist(err) {
		log.Infof("Unable to find custom exiftool: '%s'. Trying to use '%s' from $PATH", err, tool.cmd)
		return
	}

	tool.cmd = customePath
}

func (tool extifToolWrapper) exec(agrs []string) {
	cmdArgs := append(tool.defaultArgs, agrs...)
	cmd := exec.Command(tool.cmd, cmdArgs...)

	log.Debugf("ExifTool command: '%s'\n", cmd.String())

	cmd.Stdout = os.Stdout
	if verbose {
		cmd.Stderr = os.Stderr
	}

	err := cmd.Run()
	if err != nil {
		log.Warningf("ExifTool exec error: '%s'", err)
	}
}

func execExifTool(agrs []string) {
	if exifToolObj == nil {
		exifToolObj = newExtifTool()
	}

	exifToolObj.exec(agrs)
}