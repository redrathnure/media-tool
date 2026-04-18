//go:build !windows
// +build !windows

package linux_udisks

import (
	_udisks "github.com/sandbankdisperser/go-udisks"

	"github.com/redrathnure/media-tool/core/removable/core"
)

func FindDevices(deviceFilter core.DeviceFilter) []*core.RemovableDevice {
	result := []*core.RemovableDevice{}

	client, err := _udisks.NewClient()
	if err != nil {
		log.Warning("Unable to initialize udisks client: ", err)
		return result
	}

	devs, err := client.BlockDevices()
	if err != nil {
		log.Warning("Unable to list udisks devices: ", err)
		return result
	}

	log.Infof("Scanning udisks devices...")

	for _, dev := range devs {
		if dev.Drive != nil && dev.Drive.Removable && dev.Drive.Ejectable {
			udev := NewUdiskDev(dev)
			log.Infof("Found removable device: '%s' (vendor: %s, model: %s)", udev.Name(), dev.Drive.Vendor, dev.Drive.Model)
			if !dev.IsMounted() {
				log.Infof("Device is not mounted. Skipping...")
				continue
			}
			if deviceFilter.Accept(&udev) {
				result = append(result, &udev)
			} else {
				log.Infof("Device was rejected by filters. Skipping...")
			}
		}
	}
	return result
}
