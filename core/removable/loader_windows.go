//go:build windows
// +build windows

package removable

import (
	"github.com/redrathnure/media-tool/core/removable/core"
	removable "github.com/redrathnure/media-tool/core/removable/windows_removable"
	wdp "github.com/redrathnure/media-tool/core/removable/windows_wdp"
)

func (d *MtpDownloader) findDevices(deviceFilter core.DeviceFilter) []*core.RemovableDevice {
	result := removable.FindDevices(deviceFilter)
	result = append(result, wdp.FindDevices(deviceFilter)...)
	return result
}

func (d *MtpDownloader) close() {
	wdp.Close()
}
