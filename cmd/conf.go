package cmd

import (
	"os"
	"path"
	"path/filepath"

	"github.com/redrathnure/media-tool/core/config"

	"github.com/redrathnure/media-tool/core/tools"
)

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	builder := config.NewConfig()
	if Context.cfgFile != "" {
		log.Infof("Loading configuration from '%s' file only. The rest locations will be ignored.", Context.cfgFile)
		builder.SetConfigFile(Context.cfgFile)
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

	Context.Conf = newConf

	tools.ExifToolPath = newConf.GetExifToolPath()
}
