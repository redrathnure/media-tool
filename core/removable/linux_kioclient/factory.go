//go:build !windows
// +build !windows

package linux_kioclient

import (
	"github.com/redrathnure/media-tool/core/removable/core"
)

func FindDevices(deviceFilter core.DeviceFilter) []core.RemovableDevice {
	result := []core.RemovableDevice{}

	client := NewKioClient()

	if !client.CanWork() {
		log.Debugf("kioclient was not found. Unable to use 'mtp:/' devices")
		return result
	}

	log.Infof("Scanning mtp:/ devices (using 'kioclient')...")

	devices := client.Ls("mtp:/")

	for _, deviceName := range devices {
		deviceUrl := "mtp:/" + deviceName
		log.Infof("Found mtp device: '%s'", deviceUrl)

		disks := client.Ls("mtp:/" + deviceName)
		for _, diskName := range disks {
			dev := NewKioDev(deviceName, diskName, client)
			log.Infof("Found drives on mtp device: '%s'", diskName)

			if deviceFilter.Accept(&dev) {
				result = append(result, dev)
			} else {
				log.Infof("Device was rejected by filters. Skipping...")
			}
		}
	}

	return result
}
