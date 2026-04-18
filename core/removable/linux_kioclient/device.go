//go:build !windows
// +build !windows

package linux_kioclient

import (
	"os"
	"path"

	"github.com/cheggaaa/pb/v3"

	"github.com/redrathnure/media-tool/core/removable/core"
)

type KioDevice struct {
	client  *KioClient
	baseUrl string
	name    string
}

func NewKioDev(deviceName string, driveName string, client *KioClient) core.RemovableDevice {
	return &KioDevice{baseUrl: path.Join("mtp:", deviceName, driveName), name: path.Join(deviceName, driveName), client: client}
}

func (d *KioDevice) Name() string {
	return d.name
}

func (d *KioDevice) HasFile(deviceFile string) bool {
	fullPath := path.Join(d.baseUrl, deviceFile)
	entries := d.client.Ls(fullPath)
	return len(entries) > 0
}

func (d *KioDevice) CopyFile(srcDeviceFile string, dstFile string, progressBar *pb.ProgressBar) (copyBytes int64, err error) {
	fullPath := path.Join(d.baseUrl, srcDeviceFile)

	if err := d.client.Copy(fullPath, dstFile); err != nil {
		return 0, err
	}

	fi, err := os.Stat(dstFile)
	if err != nil {
		os.Remove(dstFile)
		return 0, err
	}

	return fi.Size(), nil
}

func (d *KioDevice) DeleteFile(deviceFile string) error {
	fullPath := path.Join(d.baseUrl, deviceFile)

	d.client.Rm(fullPath)

	return nil
}

func (d *KioDevice) GetChildren(deviceFile string) (children []*core.FileDescriptor, err error) {
	fullPath := path.Join(d.baseUrl, deviceFile)

	children = []*core.FileDescriptor{}

	entries := d.client.Ls(fullPath)

	for _, entry := range entries {
		name := entry
		isDir, size := d.client.GetState(path.Join(fullPath, name))

		children = append(children, &core.FileDescriptor{Name: path.Join(deviceFile, name), IsDir: isDir, Size: size})
	}
	return children, nil
}
