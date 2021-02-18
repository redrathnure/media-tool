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
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	cfgExifToolPath = "exiftool.path"
)

type exifToolWrapper struct {
	cmd         string
	defaultArgs []string
}

type exifToolArgs struct {
	args []string
}

var exifToolObj *exifToolWrapper

func newExifTool() *exifToolWrapper {
	result := exifToolWrapper{
		cmd:         "exiftool",
		defaultArgs: []string{"-v0", "-progress"},
	}
	result.initCmd()
	return &result
}

func getExifTool() *exifToolWrapper {
	if exifToolObj == nil {
		exifToolObj = newExifTool()
	}
	return exifToolObj
}

func (tool *exifToolWrapper) initCmd() {
	customPath := viper.GetString(cfgExifToolPath)
	if customPath != "" {

		if strings.Contains(customPath, "$APP_DIR") {
			ex, err := os.Executable()
			if err != nil {
				log.Infof("Unable to find custom exiftool: '%s'. Trying to use '%s' from $PATH", err, tool.cmd)
				return
			}
			customPath = strings.ReplaceAll(customPath, "$APP_DIR", filepath.Dir(ex))
			customPath, err = filepath.Abs(customPath)
			if err != nil {
				log.Infof("Unable to find custom exiftool: '%s'. Trying to use '%s' from $PATH", err, tool.cmd)
				return
			}
			if _, err := os.Stat(customPath); os.IsNotExist(err) {
				log.Infof("Unable to find custom exiftool: '%s'. Trying to use '%s' from $PATH", err, tool.cmd)
				return
			}
		}

		tool.cmd = customPath
	}
}

func (tool *exifToolWrapper) exec(args *exifToolArgs) {
	cmd := exec.Command(tool.cmd, args.args...)

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

func (tool *exifToolWrapper) newArgs() *exifToolArgs {
	return &exifToolArgs{args: tool.defaultArgs}
}

func (toolArgs *exifToolArgs) add(args ...string) {
	toolArgs.args = append(toolArgs.args, args...)
}

func (toolArgs *exifToolArgs) recursively() {
	toolArgs.add("-r")
}

func (toolArgs *exifToolArgs) src(dirOrFilepath string) {
	toolArgs.add(dirOrFilepath)
}

func (toolArgs *exifToolArgs) forImages() {
	toolArgs.add("-ext", "jpg")
	toolArgs.add("-ext", "nef")
}

func (toolArgs *exifToolArgs) forVideoMp4() {
	toolArgs.add("-ext", "mp4")
}

func (toolArgs *exifToolArgs) forVideoLrv() {
	toolArgs.add("-ext", "LRV")
}

func (toolArgs *exifToolArgs) forVideoAvchd() {
	toolArgs.add("-ext", "mts")
}

func (toolArgs *exifToolArgs) forDateFormat(dateFormat string) {
	toolArgs.add("-d", dateFormat)
}

func (toolArgs *exifToolArgs) changeTag(tagName string, tagValue string) {
	toolArgs.add(fmt.Sprintf("-%s<%s", tagName, tagValue))
}

func (toolArgs *exifToolArgs) changeFileDate(tagValue string) {
	//File:
	toolArgs.changeTag("FileModifyDate", tagValue)
	toolArgs.changeTag("FileCreateDate", tagValue)
}

func (toolArgs *exifToolArgs) changeExifDate(tagValue string) {
	//'EXIF:
	toolArgs.changeTag("CreateDate", tagValue)
	toolArgs.changeTag("DateTimeOriginal", tagValue)
}

func (toolArgs *exifToolArgs) changeMp4Date(tagValue string) {
	//quicktime:
	toolArgs.changeTag("CreateDate", tagValue)
	toolArgs.changeTag("ModifyDate", tagValue)
	toolArgs.changeTag("TrackCreateDate", tagValue)
	toolArgs.changeTag("TrackModifyDate", tagValue)
	toolArgs.changeTag("MediaCreateDate", tagValue)
	toolArgs.changeTag("MediaModifyDate", tagValue)
}

func init() {
	viper.SetDefault(cfgExifToolPath, "")
}
