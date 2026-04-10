//go:build !windows
// +build !windows

package linux_udisks

import (
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("udisks")
