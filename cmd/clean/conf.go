package clean

import (
	root "github.com/redrathnure/media-tool/cmd"
	"github.com/redrathnure/media-tool/core/config"
)

func getConf() *config.Config {
	return root.Context.Conf
}
