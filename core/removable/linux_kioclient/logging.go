//go:build !windows
// +build !windows

package linux_kioclient

import (
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("kioclient")
