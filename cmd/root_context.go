package cmd

import "github.com/redrathnure/media-tool/core/config"

type RootContext struct {
	cfgFile string
	Verbose bool
	Conf    *config.Config
}

var Context = RootContext{Conf: config.Empty()}
