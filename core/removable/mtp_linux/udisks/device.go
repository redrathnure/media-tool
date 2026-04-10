//go:build !windows
// +build !windows

package udisks

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/cheggaaa/pb/v3"

	_udisks "github.com/sandbankdisperser/go-udisks"

	"github.com/redrathnure/media-tool/core/removable/mtp_linux/core"
)

type UdiskDevice struct {
	dev *_udisks.BlockDevice
}

func NewUdiskDev(dev *_udisks.BlockDevice) core.RemovableDevice {
	return &UdiskDevice{dev: dev}
}

func (d *UdiskDevice) Name() string {
	return d.dev.IdLabel
}

func (d *UdiskDevice) HasFile(deviceFile string) bool {
	rootDir, err := d.getMount()
	if err != nil {
		return false
	}

	_, err = os.Stat(path.Join(rootDir, deviceFile))
	return err == nil
}

func (d *UdiskDevice) CopyFile(srcDeviceFile string, dstFile string, progressBar *pb.ProgressBar) (copyBytes int64, err error) {
	rootDir, err := d.getMount()
	if err != nil {
		return 0, err
	}

	srcFile := path.Join(rootDir, srcDeviceFile)

	return d.copyFileManual(srcFile, dstFile, progressBar)
}

func (d *UdiskDevice) DeleteFile(deviceFile string) error {
	rootDir, err := d.getMount()
	if err != nil {
		return err
	}

	file := path.Join(rootDir, deviceFile)

	if err := os.Remove(file); err != nil {
		return err
	}
	return nil
}

func (d *UdiskDevice) GetChildren(deviceFile string) (children []core.FileDescriptor, err error) {
	rootDir, err := d.getMount()
	if err != nil {
		return nil, err
	}

	baseDir := path.Join(rootDir, deviceFile)
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}

	result := []core.FileDescriptor{}

	for _, entry := range entries {
		isDir := entry.IsDir()
		name := entry.Name()

		size := int64(0)
		if !isDir {
			fi, err := os.Stat(path.Join(baseDir, name))
			if err != nil {
				return nil, err
			}
			size = fi.Size()
		}

		result = append(result, core.FileDescriptor{Name: path.Join(deviceFile, name), IsDir: isDir, Size: size})
	}
	return result, nil
}

func (d *UdiskDevice) getMount() (rootDir string, err error) {
	if !d.dev.IsMounted() {
		return "", fmt.Errorf("Is not mounted")
	}

	for _, fs := range d.dev.Filesystems {
		if fs.IsMounted() && len(fs.MountPoints) > 0 {
			return fs.MountPoints[0], nil
		}
	}
	return "", fmt.Errorf("Is not mounted")
}

func (*UdiskDevice) copyFileManual(srcFile string, dstFile string, progressBar *pb.ProgressBar) (copyBytes int64, err error) {
	//TODO extract to tools and reuse in backups
	src, err := os.Open(srcFile)
	if err != nil {
		return 0, err
	}
	defer src.Close()

	dst, err := os.Create(dstFile)
	if err != nil {
		return 0, err
	}
	defer dst.Close()

	proxyWriter := progressBar.NewProxyWriter(dst)
	defer proxyWriter.Close()

	result, err := io.Copy(proxyWriter, src)
	if err != nil {
		return result, err
	}
	src.Close()
	proxyWriter.Close()
	dst.Close()

	fi, err := os.Stat(srcFile)
	if err != nil {
		os.Remove(dstFile)
		return result, err
	}
	if err := os.Chmod(dstFile, fi.Mode()); err != nil {
		os.Remove(dstFile)
		return result, err
	}
	return result, nil
}
