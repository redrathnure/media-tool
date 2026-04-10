//go:build !windows
// +build !windows

package core

import (
	"path"
	"strings"
)

type DeviceFilter interface {
	Accept(device *RemovableDevice) bool
}

func HasFileFilter(fileName string) DeviceFilter {
	return fileFilter{fileName: fileName}
}

type fileFilter struct {
	fileName string
}

func (f fileFilter) Accept(device *RemovableDevice) bool {
	return device != nil && (*device).HasFile(f.fileName)
}

func HasDeviceNameFilter(deviceName string) DeviceFilter {
	return deviceNameFilter{deviceName: deviceName}
}

type deviceNameFilter struct {
	deviceName string
}

func (filter deviceNameFilter) Accept(device *RemovableDevice) bool {
	return device != nil && strings.Contains((*device).Name(), filter.deviceName)
}

func NotFilter(filter DeviceFilter) DeviceFilter {
	return notFilter{filter: filter}
}

type notFilter struct {
	filter DeviceFilter
}

func (filter notFilter) Accept(device *RemovableDevice) bool {
	return !filter.filter.Accept(device)
}

func AndFilter(filters ...DeviceFilter) DeviceFilter {
	return andFilter{filters: filters}
}

type andFilter struct {
	filters []DeviceFilter
}

func (filter andFilter) Accept(device *RemovableDevice) bool {
	result := true

	for _, f := range filter.filters {
		result = result && f.Accept(device)
	}
	return result
}

func OrFilter(filters ...DeviceFilter) DeviceFilter {
	return orFilter{filters: filters}
}

type orFilter struct {
	filters []DeviceFilter
}

func (filter orFilter) Accept(device *RemovableDevice) bool {
	result := false

	for _, f := range filter.filters {
		result = result || f.Accept(device)
	}
	return result
}

var DCIM_DIR string = "DCIM"
var CAM_FILES_DIR string = path.Join("PRIVATE", "AVCHD", "BDMV", "STREAM")
var GOPRO_DIR string = path.Join("DCIM", "100GOPRO")

var heroDeviceNameFilter DeviceFilter = HasDeviceNameFilter("HERO")
var goProDeviceNameFilter DeviceFilter = HasDeviceNameFilter("GoPro")
var goProFolderFilter DeviceFilter = HasFileFilter(GOPRO_DIR)
var GoProFilter DeviceFilter = OrFilter(heroDeviceNameFilter, goProDeviceNameFilter, goProFolderFilter)

var camDeviceNameFilter DeviceFilter = HasDeviceNameFilter("CAM")
var streamFolderFilter DeviceFilter = HasFileFilter(CAM_FILES_DIR)
var CamFilter DeviceFilter = AndFilter(camDeviceNameFilter, streamFolderFilter)

var dcmiFilderFilter DeviceFilter = HasFileFilter(DCIM_DIR)
var SdPhotosFilter DeviceFilter = AndFilter(dcmiFilderFilter, NotFilter(GoProFilter), NotFilter(CamFilter))
