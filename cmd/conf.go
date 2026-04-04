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
	"path"
	"path/filepath"

	"github.com/redrathnure/media-tool/core/config"
)

var conf = config.Empty()

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	builder := config.NewConfig()
	if cfgFile != "" {
		log.Infof("Loading configuration from '%s' file only. The rest locations will be ignored.", cfgFile)
		builder.SetConfigFile(cfgFile)
	} else {
		builder.AddConfigLocation("/", "etc", "media-tool")
		builder.AddConfigLocation("$HOME", "media-tool")
		builder.AddConfigLocation("$HOME", ".config", "media-tool")
		if ex, err := os.Executable(); err == nil {
			rootConfigDir := path.Join(filepath.Dir(ex), "conf")
			builder.AddConfigLocation(rootConfigDir)
		}
		builder.AddConfigLocation(".", "conf")
	}

	newConf, err := builder.Build()
	if err != nil {
		log.Errorf("%v", err)
		os.Exit(1)
	}

	conf = newConf
}
