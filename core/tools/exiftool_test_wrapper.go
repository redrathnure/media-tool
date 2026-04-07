package tools

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// testExifToolWrapper is a test implementation of exifToolWrapper that tracks exec calls
type testExifToolWrapper struct {
	ExifToolWrapper
	Calls []string
}

type ExifRunArgs struct {
	Expected   []string
	Unexpected []string
}

type ExifSingleCallCaseData struct {
	Name       string
	Args       []string
	Recursive  bool
	Expected   []string
	Unexpected []string
}

type ExifMultiCallCaseData struct {
	Name      string
	Args      []string
	Recursive bool
	CallArgs  []ExifRunArgs
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

func (testTool *testExifToolWrapper) AssertCalls(t *testing.T, expectedCalls int) {
	assert.Len(t, testTool.Calls, expectedCalls, "exiftool exec should be called in %s", t.Name())
}

func (testTool *testExifToolWrapper) AssertCallArgs(t *testing.T, callIndex int, sourcePath string, caseData ExifRunArgs) {

	callArg := testTool.Calls[callIndex]

	assert.Contains(t, callArg, sourcePath, "source path not set in %s (call: %s)", t.Name(), callIndex)
	for _, tag := range caseData.Expected {
		assert.Contains(t, callArg, tag, "missing expected tag in %s (call: %s)", t.Name(), callIndex)
	}

	for _, tag := range caseData.Unexpected {
		assert.NotContains(t, callArg, tag, "found unexpected tag in %s (call: %s)", t.Name(), callIndex)
	}
}

func (caseData *ExifSingleCallCaseData) CallArgs() ExifRunArgs {
	return ExifRunArgs{Expected: caseData.Expected, Unexpected: caseData.Unexpected}
}
