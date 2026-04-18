//go:build windows
// +build windows

package windows_wdp

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"

	"github.com/redrathnure/media-tool/core/removable/core"
	"github.com/redrathnure/media-tool/core/tools"
	"github.com/tobwithu/gowpd"
)

var ignoreFiles = []string{"System Volume Information", "$RECYCLE.BIN"}

type WpdDevice struct {
	dev        *gowpd.Device
	deviceId   int
	wdpObjects map[string]*gowpd.Object
	name       string
}

func NewWpdDevice(dev *gowpd.Device, deviceId int) core.RemovableDevice {
	result := WpdDevice{dev: dev, deviceId: deviceId, wdpObjects: map[string]*gowpd.Object{}}
	result.name = fmt.Sprintf("MTP#%v - '%v (%v)'", deviceId, gowpd.GetDeviceName(deviceId), gowpd.GetDeviceDescription(deviceId))
	return &result
}

func (d *WpdDevice) Name() string {
	return d.name
}

func (d *WpdDevice) HasFile(deviceFile string) bool {
	wpdObj := d.findObject(deviceFile)
	return wpdObj != nil
}

func (d *WpdDevice) findObject(deviceFile string) *gowpd.Object {
	wpdObj, ok := d.wdpObjects[deviceFile]

	if !ok {
		wpdFileName := gowpd.PathSeparator + deviceFile
		wpdObj = d.dev.FindObject(wpdFileName)
		if wpdObj == nil {
			log.Debugf("'%v' was not found.", deviceFile)
		} else {
			d.wdpObjects[deviceFile] = wpdObj
			ok = true
		}
	}

	return wpdObj
}

func (d *WpdDevice) CopyFile(srcDeviceFile string, dstFile string, progressBar *pb.ProgressBar) (copyBytes int64, err error) {
	obj := d.findObject(srcDeviceFile)

	id := obj.Id

	reader, err := d.dev.GetReader(id)
	if err != nil {
		return 0, err
	}
	defer reader.Close()

	f, err := os.Create(dstFile)
	if err != nil {
		return 0, err
	}
	writer := gowpd.NewBufWriteCloser(f, tools.CopyBufferSize)
	defer writer.Close()

	proxyWriter := progressBar.NewProxyWriter(writer)

	written, err := io.Copy(proxyWriter, reader)

	if err != nil {
		return 0, err
	}
	return written, gowpd.SetFileTime(dstFile, obj.ModTime)
}

func (d *WpdDevice) DeleteFile(deviceFile string) error {
	wpdObj := d.findObject(deviceFile)
	if wpdObj != nil && !wpdObj.IsDir {
		return d.dev.Delete(wpdObj.Id)
	}

	return nil
}

func (d *WpdDevice) GetChildren(deviceFile string) (children []*core.FileDescriptor, err error) {
	wpdObj := d.findObject(deviceFile)
	if wpdObj == nil {
		return nil, fmt.Errorf("'%s' was not found", deviceFile)
	}

	objs, err := d.dev.GetChildObjects(wpdObj.Id)
	if err != nil {
		return nil, err
	}

	result := []*core.FileDescriptor{}
	for _, o := range objs {

		if d.isIgnored(o.Name) {
			log.Debugf("Skipping '%v' file", o.Name)
			continue
		}

		name := filepath.Join(deviceFile, o.Name)
		isDir := o.IsDir
		size := int64(0)
		if !isDir {
			size = o.Size
		}
		log.Debugf("Found: '%s' file", name)

		result = append(result, &core.FileDescriptor{Name: name, IsDir: isDir, Size: size})
	}

	return result, nil
}

func (*WpdDevice) isIgnored(fileName string) bool {
	for _, ignore := range ignoreFiles {
		if fileName == ignore {
			return true
		}
	}
	return false
}
