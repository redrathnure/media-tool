//go:build !windows
// +build !windows

package removable

import (
	"github.com/redrathnure/media-tool/core/removable/core"
	kioclient "github.com/redrathnure/media-tool/core/removable/linux_kioclient"
	udisks "github.com/redrathnure/media-tool/core/removable/linux_udisks"
)

func (d *MtpDownloader) findDevices(deviceFilter core.DeviceFilter) []core.RemovableDevice {
	result := udisks.FindDevices(deviceFilter)
	result = append(result, kioclient.FindDevices(deviceFilter)...)
	return result
}

func (d *MtpDownloader) close() {
}
