package core

import "github.com/cheggaaa/pb/v3"

type FileDescriptor struct {
	Name  string
	IsDir bool
	Size  int64
}

type RemovableDevice interface {
	Name() string

	HasFile(deviceFile string) bool

	CopyFile(srcDeviceFile string, dstFile string, progressBar *pb.ProgressBar) (bytesCount int64, err error)
	DeleteFile(deviceFile string) error

	GetChildren(deviceFile string) (children []FileDescriptor, err error)
}
