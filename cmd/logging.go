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
	if verbose {
		logging.SetLevel(logging.DEBUG, "")
	} else {
		logging.SetLevel(logging.INFO, "")
	}
}
