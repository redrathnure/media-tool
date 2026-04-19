//go:build windows
// +build windows

package windows_removable

import (
	"io"
	"os"
	"path"

	"github.com/cheggaaa/pb/v3"

	"github.com/redrathnure/media-tool/core/removable/core"
)

var ignoreFiles = []string{"System Volume Information", "$RECYCLE.BIN"}

type WinDrive struct {
	name string
}

func NewWinDevice(driveLabel string) core.RemovableDevice {
	result := WinDrive{name: driveLabel}
	return &result
}

func (d *WinDrive) Name() string {
	return d.name
}

func (d *WinDrive) HasFile(deviceFile string) bool {
	_, err := os.Stat(path.Join(d.name, deviceFile))
	return err == nil
}

func (d *WinDrive) CopyFile(srcDeviceFile string, dstFile string, progressBar *pb.ProgressBar) (copyBytes int64, err error) {
	srcFile := path.Join(d.name, srcDeviceFile)

	return d.copyFileManual(srcFile, dstFile, progressBar)
}

func (d *WinDrive) DeleteFile(deviceFile string) error {
	file := path.Join(d.name, deviceFile)

	if err := os.Remove(file); err != nil {
		return err
	}
	return nil
}

func (d *WinDrive) GetChildren(deviceFile string) (children []*core.FileDescriptor, err error) {
	baseDir := path.Join(d.name, deviceFile)
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}

	result := []*core.FileDescriptor{}

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

		result = append(result, &core.FileDescriptor{Name: path.Join(deviceFile, name), IsDir: isDir, Size: size})
	}
	return result, nil
}

func (*WinDrive) isIgnored(fileName string) bool {
	for _, ignore := range ignoreFiles {
		if fileName == ignore {
			return true
		}
	}
	return false
}

func (*WinDrive) copyFileManual(srcFile string, dstFile string, progressBar *pb.ProgressBar) (copyBytes int64, err error) {
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

	result, err := io.Copy(proxyWriter, src)
	if err != nil {
		return result, err
	}

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
