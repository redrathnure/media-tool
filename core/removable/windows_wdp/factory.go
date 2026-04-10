//go:build windows
// +build windows

package windows_wdp

import (
	"github.com/redrathnure/media-tool/core/removable/core"
	"github.com/tobwithu/gowpd"
)

func FindDevices(deviceFilter core.DeviceFilter) []core.RemovableDevice {
	result := []core.RemovableDevice{}

	if err := gowpd.Init(); err != nil {
		log.Warningf("Unable to init WDP devices")
		return result
	}

	mtpDeviceCount := gowpd.GetDeviceCount()

	log.Infof("Checking %d WDP devices...", mtpDeviceCount)

	for i := 0; i < mtpDeviceCount; i++ {
		log.Infof("Open WDP  device#%d...", i)

		wdpDev, err := gowpd.ChooseDevice(i)
		if err != nil {
			log.Warningf("Unable to open device: %s", err)

			return result
		}

		udev := NewWpdDevice(wdpDev, i)
		log.Infof("Found removable device: '%s'", udev.Name())
		if deviceFilter.Accept(&udev) {
			result = append(result, udev)
		} else {
			log.Infof("Device was rejected by filters. Skipping...")
		}
	}
	return result
}

func Close() {
	gowpd.Destroy()
}
