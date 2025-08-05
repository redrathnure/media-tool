package cmd

import "os/exec"

// testExifToolWrapper is a test implementation of exifToolWrapper that tracks exec calls
type testExifToolWrapper struct {
	exifToolWrapper
	execCalled bool
}

func newTestExifTool() *testExifToolWrapper {
	tool := &testExifToolWrapper{
		exifToolWrapper: *newExifTool(),
	}
	mockExecCommand := func(name string, args ...string) *exec.Cmd {
		tool.execCalled = true
		return exec.Command(name, args...)
	}

	tool.exifToolWrapper.execCommand = mockExecCommand

	exifToolObj = &tool.exifToolWrapper

	return tool
}

func (testTool *testExifToolWrapper) clear() {
	exifToolObj = nil
}
