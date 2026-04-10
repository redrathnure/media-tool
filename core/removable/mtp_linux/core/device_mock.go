//go:build !windows
// +build !windows

package core

import (
	"fmt"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

type MockRemovableDevice struct {
	name  string
	files []string
}

func NewMockDev(name string, files ...string) RemovableDevice {
	return &MockRemovableDevice{name: name, files: files}
}

func (d *MockRemovableDevice) Name() string {
	return d.name
}

func (d *MockRemovableDevice) HasFile(deviceFile string) bool {
	for _, file := range d.files {
		if strings.HasPrefix(file, deviceFile) {
			return true
		}
	}
	return false
}

func (d *MockRemovableDevice) CopyFile(srcDeviceFile string, dstFile string, progressBar *pb.ProgressBar) (int64, error) {
	if d.HasFile(srcDeviceFile) {
		return 3, nil
	}
	return -1, fmt.Errorf("No files found")
}

func (d *MockRemovableDevice) DeleteFile(deviceFile string) error {
	newFiles := []string{}
	defer func() {
		d.files = newFiles
	}()

	for _, file := range d.files {
		if !strings.HasPrefix(file, deviceFile) {
			newFiles = append(newFiles, file)
		}
	}

	if len(newFiles) == len(d.files) {
		return fmt.Errorf("File '%s' was not found", deviceFile)
	}

	return nil
}

func (d *MockRemovableDevice) GetChildren(deviceFile string) (children []FileDescriptor, err error) {
	result := []FileDescriptor{}

	for _, file := range d.files {
		if strings.HasPrefix(file, deviceFile) {
			result = append(result, FileDescriptor{Name: file, IsDir: false, Size: 10})
		}
	}
	return result, nil
}
