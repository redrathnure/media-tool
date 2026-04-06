package main

import (
	"github.com/redrathnure/media-tool/cmd"

	_ "github.com/redrathnure/media-tool/cmd/clean"
	_ "github.com/redrathnure/media-tool/cmd/config"
	_ "github.com/redrathnure/media-tool/cmd/fix"
	_ "github.com/redrathnure/media-tool/cmd/import_"
	_ "github.com/redrathnure/media-tool/cmd/version"
)

func main() {
	cmd.Execute()
}
