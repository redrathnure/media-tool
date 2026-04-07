package tools

import (
	"os/exec"
	"strings"
)

// testExifToolWrapper is a test implementation of exifToolWrapper that tracks exec calls
type testExifToolWrapper struct {
	ExifToolWrapper
	Calls []string
}

func NewTestExifTool() *testExifToolWrapper {
	tool := &testExifToolWrapper{
		ExifToolWrapper: *newExifTool(),
	}
	mockExecCommand := func(name string, args ...string) *exec.Cmd {
		tool.Calls = append(tool.Calls, strings.Join(args, " "))
		return exec.Command(name, args...)
	}

	tool.ExifToolWrapper.execCommand = mockExecCommand

	exifToolObj = &tool.ExifToolWrapper

	return tool
}

func (testTool *testExifToolWrapper) Clear() {
	exifToolObj = nil
	testTool.Calls = []string{}
}
