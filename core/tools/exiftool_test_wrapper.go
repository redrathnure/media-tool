package tools

import "os/exec"

// testExifToolWrapper is a test implementation of exifToolWrapper that tracks exec calls
type testExifToolWrapper struct {
	ExifToolWrapper
	ExecCalled bool
}

func NewTestExifTool() *testExifToolWrapper {
	tool := &testExifToolWrapper{
		ExifToolWrapper: *newExifTool(),
	}
	mockExecCommand := func(name string, args ...string) *exec.Cmd {
		tool.ExecCalled = true
		return exec.Command(name, args...)
	}

	tool.ExifToolWrapper.execCommand = mockExecCommand

	exifToolObj = &tool.ExifToolWrapper

	return tool
}

func (testTool *testExifToolWrapper) Clear() {
	exifToolObj = nil
}
