package fix

import (
	"testing"

	"github.com/redrathnure/media-tool/core/tools"
	"github.com/redrathnure/media-tool/core/tools/testutil"
	"github.com/stretchr/testify/assert"
)

func TestFileDatesCmd_Initialization(t *testing.T) {
	assert.NotNil(t, fileDatesCmd)
	assert.NotNil(t, fileDatesCmd.Use)
	assert.NotNil(t, fileDatesCmd.Short)
	assert.NotNil(t, fileDatesCmd.Long)
}

func TestFileDatesCmd_Flags(t *testing.T) {
	persistentFlags := fileDatesCmd.PersistentFlags()
	assert.False(t, persistentFlags.HasFlags())

	flags := fileDatesCmd.Flags()
	assert.False(t, flags.HasFlags())
}

func TestFileDatesCmd_CommandStructure(t *testing.T) {
	assert.Equal(t, "file_dates", fileDatesCmd.Name())
	assert.Equal(t, "media-tool fix file_dates", fileDatesCmd.CommandPath())

	flags := fileDatesCmd.Flags()
	assert.Equal(t, 0, flags.NFlag())
}

func TestFileDatesCmd_ArgValidation(t *testing.T) {
	// Test no args (should fail)
	err := fileDatesCmd.Args(fileDatesCmd, []string{})
	assert.NoError(t, err)

	// Test single arg (should pass)
	err = fileDatesCmd.Args(fileDatesCmd, []string{"test.jpg"})
	assert.NoError(t, err)

	// Test multiple args (should fail)
	err = fileDatesCmd.Args(fileDatesCmd, []string{"test1.jpg", "test2.jpg"})
	assert.Error(t, err)
}

func TestRunFixFileDates(t *testing.T) {
	// Save original values to restore after test
	origRecursively := recursively
	origDryRun := dryRun
	defer func() {
		recursively = origRecursively
		dryRun = origDryRun
	}()

	tests := []tools.ExifSingleCallCaseData{
		{
			Name:      "without recursion",
			Args:      []string{"test.jpg"},
			Recursive: false,
			Expected: []string{
				"-FileModifyDate<CreateDate",
				//Win specific
				//"-FileCreateDate<CreateDate",
			},
			Unexpected: []string{
				"-r",
			},
		},
		{
			Name:      "with recursion",
			Args:      []string{"test.jpg"},
			Recursive: true,
			Expected: []string{
				"-FileModifyDate<CreateDate",
				//Win specific
				//"-FileCreateDate<CreateDate",
				"-r",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			recursively = tt.Recursive
			dryRun = false

			// Create a test exiftool wrapper
			testTool := tools.NewTestExifTool()
			defer testTool.Clear()

			// Run the command
			runFixFileDates(fileDatesCmd, tt.Args)

			// Check ExifTool calls
			testTool.AssertCalls(t, 1)

			testTool.AssertCallArgs(t, 0, tt.Args[0], tt.CallArgs())
		})
	}
}

func TestRunFixFileDates_DryRun(t *testing.T) {
	// Save original values to restore after test
	origRecursively := recursively
	origDryRun := dryRun
	defer func() {
		recursively = origRecursively
		dryRun = origDryRun
	}()

	tests := []tools.ExifSingleCallCaseData{
		{
			Name:      "without recursion",
			Args:      []string{"test.jpg"},
			Recursive: false,
			Expected: []string{
				"-WriteNothing<CreateDate",
			},
			Unexpected: []string{
				"-FileModifyDate<CreateDate",
				"-FileCreateDate<CreateDate",
				"-r",
			},
		},
		{
			Name:      "with recursion",
			Args:      []string{"test.jpg"},
			Recursive: true,
			Expected: []string{
				"-WriteNothing<CreateDate",
				"-r",
			},
			Unexpected: []string{
				"-FileModifyDate<CreateDate",
				"-FileCreateDate<CreateDate",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			recursively = tt.Recursive
			dryRun = true

			// Create a test exiftool wrapper
			testTool := tools.NewTestExifTool()
			defer testTool.Clear()

			// Run the command
			runFixFileDates(fileDatesCmd, tt.Args)

			// Check ExifTool calls
			testTool.AssertCalls(t, 1)

			testTool.AssertCallArgs(t, 0, tt.Args[0], tt.CallArgs())
		})
	}
}

func TestFixFileDates_DryRun(t *testing.T) {
	tmp := testutil.NewTempDir(t)

	tmp.MkCr3FullExif("src", "zxc.cr3")
	tmp.MkJpgFullExif("src", "zxc.jpg")
	tmp.MkJpegFullExif("src", "zxc.jpeg")
	tmp.MkVideoMkv("src", "zxc.mkv")
	tmp.MkVideoMp4("src", "zxc.mp4")

	srcDir := tmp.DirName("src")

	fixFileDates(srcDir, true, true, true)

	assert.True(t, tmp.IsDirExist("src"))

	assert.True(t, tmp.IsFileExist("src", "zxc.cr3"))
	assert.NotEqual(t, testutil.CreationDate, tmp.FileDate("src", "zxc.cr3"))
	assert.True(t, tmp.IsFileExist("src", "zxc.jpg"))
	assert.NotEqual(t, testutil.CreationDate, tmp.FileDate("src", "zxc.jpg"))
	assert.True(t, tmp.IsFileExist("src", "zxc.jpeg"))
	assert.NotEqual(t, testutil.CreationDate, tmp.FileDate("src", "zxc.jpeg"))
	assert.True(t, tmp.IsFileExist("src", "zxc.mkv"))
	assert.NotEqual(t, testutil.CreationDate, tmp.FileDate("src", "zxc.mkv"))
	assert.True(t, tmp.IsFileExist("src", "zxc.mp4"))
	assert.NotEqual(t, testutil.CreationDate, tmp.FileDate("src", "zxc.mp4"))
}

func TestFixFileDates_Fix(t *testing.T) {

	tmp := testutil.NewTempDir(t)

	tmp.MkCr3FullExif("src", "zxc.cr3")
	tmp.MkJpgFullExif("src", "zxc.jpg")
	tmp.MkJpegFullExif("src", "zxc.jpeg")
	tmp.MkVideoMkv("src", "zxc.mkv")
	tmp.MkVideoMp4("src", "zxc.mp4")

	srcDir := tmp.DirName("src")

	fixFileDates(srcDir, true, false, true)

	assert.True(t, tmp.IsDirExist("src"))

	assert.True(t, tmp.IsFileExist("src", "zxc.cr3"))
	assert.Equal(t, testutil.CreationDate, tmp.FileDate("src", "zxc.cr3"))
	assert.True(t, tmp.IsFileExist("src", "zxc.jpg"))
	assert.Equal(t, testutil.CreationDate, tmp.FileDate("src", "zxc.jpg"))
	assert.True(t, tmp.IsFileExist("src", "zxc.jpeg"))
	assert.Equal(t, testutil.CreationDate, tmp.FileDate("src", "zxc.jpeg"))
	assert.True(t, tmp.IsFileExist("src", "zxc.mkv"))
	//Warning: mkv is not fully supported by exiftool
	assert.NotEqual(t, testutil.CreationDate, tmp.FileDate("src", "zxc.mkv"))
	assert.True(t, tmp.IsFileExist("src", "zxc.mp4"))
	assert.Equal(t, testutil.CreationDate, tmp.FileDate("src", "zxc.mp4"))
}
