package mtp

import (
	"path"
	"strings"

	"github.com/tobwithu/gowpd"
)

type MtpDeviceFilter interface {
	accept(deviceId int, device *gowpd.Device, deviceLabel string) bool
}

type HasFileFilter struct {
	fileName string
}

func (filter HasFileFilter) accept(deviceId int, device *gowpd.Device, deviceLabel string) bool {
	deviceFileName := gowpd.PathSeparator + filter.fileName
	obj := device.FindObject(deviceFileName)
	return obj != nil
}

type HasDeviceNameFilter struct {
	deviceName string
}

func (filter HasDeviceNameFilter) accept(deviceId int, device *gowpd.Device, deviceLabel string) bool {
	return strings.Contains(deviceLabel, filter.deviceName)
}

type NotFilter struct {
	filter MtpDeviceFilter
}

func (filter NotFilter) accept(deviceId int, device *gowpd.Device, deviceLabel string) bool {
	return !filter.filter.accept(deviceId, device, deviceLabel)
}

type AndFilter struct {
	filters []MtpDeviceFilter
}

func (filter AndFilter) accept(deviceId int, device *gowpd.Device, deviceLabel string) bool {
	result := true

	for _, f := range filter.filters {
		result = result && f.accept(deviceId, device, deviceLabel)
	}
	return result
}

type OrFilter struct {
	filters []MtpDeviceFilter
}

func (filter OrFilter) accept(deviceId int, device *gowpd.Device, deviceLabel string) bool {
	result := false

	for _, f := range filter.filters {
		result = result || f.accept(deviceId, device, deviceLabel)
	}
	return result
}

var DCIM_DIR string = "DCIM"
var CAM_FILES_DIR string = path.Join("PRIVATE", "AVCHD", "BDMV", "STREAM")
var GOPRO_DIR string = path.Join("DCIM", "100GOPRO")

var heroDeviceNameFilter MtpDeviceFilter = HasDeviceNameFilter{deviceName: "HERO"}
var goProDeviceNameFilter MtpDeviceFilter = HasDeviceNameFilter{deviceName: "GoPro"}
var _100GOPROFolderFilter MtpDeviceFilter = HasFileFilter{fileName: GOPRO_DIR}
var GoProFilter MtpDeviceFilter = OrFilter{[]MtpDeviceFilter{heroDeviceNameFilter, goProDeviceNameFilter, _100GOPROFolderFilter}}

var camDeviceNameFilter MtpDeviceFilter = HasDeviceNameFilter{deviceName: "CAM"}
var streamFolderFilter MtpDeviceFilter = HasFileFilter{fileName: CAM_FILES_DIR}
var CamFilter MtpDeviceFilter = AndFilter{[]MtpDeviceFilter{camDeviceNameFilter, streamFolderFilter}}

var dcmiFilderFilter MtpDeviceFilter = HasFileFilter{fileName: DCIM_DIR}
var SdPhotosFilter MtpDeviceFilter = AndFilter{[]MtpDeviceFilter{dcmiFilderFilter, NotFilter{GoProFilter}, NotFilter{CamFilter}}}
