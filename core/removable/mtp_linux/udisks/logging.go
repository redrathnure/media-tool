//go:build !windows
// +build !windows

package udisks

import (
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("udisks")
