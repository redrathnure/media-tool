package cmd

// testExifToolWrapper is a test implementation of exifToolWrapper that tracks exec calls
type testExifToolWrapper struct {
	exifToolWrapper
	execCalled bool
}

func newTestExifTool() *testExifToolWrapper {
	tool := &testExifToolWrapper{
		exifToolWrapper: *newExifTool(),
	}
	return tool
}

func (tool *testExifToolWrapper) exec() {
	tool.execCalled = true
}
