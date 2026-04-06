package cmd

import (
	"os"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("cmd")

func initLogger() {
	var format = logging.MustStringFormatter(
		`%{time:15:04:05} %{level:.1s} %{message}`,
	)

	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatted := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatted)
}

func initLoggerLevel() {
	if Context.Verbose {
		logging.SetLevel(logging.DEBUG, "")
	} else {
		logging.SetLevel(logging.INFO, "")
	}
}
