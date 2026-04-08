package tools

import (
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewExifTool(t *testing.T) {
	sut := newExifTool()

	assert.NotNil(t, sut)
	assert.Equal(t, "exiftool", sut.cmd)
	assert.Equal(t, []string{"-v0", "-progress"}, sut.defaultArgs)
}

func TestGetExifTool_Singleton(t *testing.T) {
	sut1 := GetExifTool()
	sut := GetExifTool()

	assert.NotNil(t, sut1)
	assert.NotNil(t, sut)
	assert.Equal(t, sut1, sut)
}

func TestExifToolWrapper_InitCmd_Default(t *testing.T) {
	sut := &ExifToolWrapper{
		cmd:         "exiftool",
		defaultArgs: []string{"-v0", "-progress"},
	}

	// Reset config to ensure clean state
	ExifToolPath = ""

	sut.initCmd()

	assert.Equal(t, "exiftool", sut.cmd)
}

func TestExifToolWrapper_InitCmd_CustomPath(t *testing.T) {
	// Test initCmd with custom path
	sut := &ExifToolWrapper{
		cmd:         "exiftool",
		defaultArgs: []string{"-v0", "-progress"},
	}

	// Reset config to ensure clean state
	ExifToolPath = "/custom/path/exiftool"

	sut.initCmd()

	assert.Equal(t, "/custom/path/exiftool", sut.cmd)

	// Reset config and singleton for next tests
	ExifToolPath = ""
	exifToolObj = nil
}

func TestExifToolWrapper_InitCmd_AppDirSubstitution(t *testing.T) {
	// Test initCmd with $APP_DIR substitution
	sut := &ExifToolWrapper{
		cmd:         "exiftool",
		defaultArgs: []string{"-v0", "-progress"},
	}

	// Reset config to ensure clean state
	ExifToolPath = ""

	// Set app directory path (this matches the name used in the exiftool.go implementation)
	ExifToolPath = path.Join("$APP_DIR", "..")

	sut.initCmd()

	// Verify $APP_DIR was replaced with the configured path
	assert.NotEqual(t, "$APP_DIR/..", sut.cmd)

	// Reset config and singleton for next tests
	ExifToolPath = ""
	exifToolObj = nil
}

func TestExifToolWrapper_DeleteOriginals_True(t *testing.T) {
	sut := newExifTool()

	sut.DeleteOriginals(true)

	args := sut.NewArgs().args

	// 2 defaults and 1 extra
	assert.Len(t, args, 3)

	assert.Equal(t, "-v0", args[0])
	assert.Equal(t, "-progress", args[1])
	assert.Equal(t, "-overwrite_original", args[2])
}

func TestExifToolWrapper_DeleteOriginals_False(t *testing.T) {
	sut := newExifTool()

	sut.DeleteOriginals(false)

	args := sut.NewArgs().args

	assert.Len(t, args, 2)

	assert.Equal(t, "-v0", args[0])
	assert.Equal(t, "-progress", args[1])
}

func TestExifToolArgs_NewArgs(t *testing.T) {
	tool := newExifTool()
	sut := tool.NewArgs()

	assert.Equal(t, &tool.args, sut)

	assert.NotNil(t, sut)
	assert.NotNil(t, sut.args)
	assert.Len(t, sut.args, 2)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
}

func TestExifToolArgs_Add(t *testing.T) {
	sut := newExifTool().NewArgs()
	assert.Len(t, sut.args, 2)

	sut.add("arg1", "arg2", "arg3")

	assert.Len(t, sut.args, 5)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "arg1", sut.args[2])
	assert.Equal(t, "arg2", sut.args[3])
	assert.Equal(t, "arg3", sut.args[4])
}

func TestExifToolWrapper_Exec(t *testing.T) {
	var capturedCmd string
	var capturedArgs []string

	// Create test instance with mock command
	mockExecCommand := func(name string, args ...string) *exec.Cmd {
		capturedCmd = name
		capturedArgs = args
		return exec.Command("echo", "test") // Use a harmless command
	}

	tool := &ExifToolWrapper{
		cmd:         "test-exiftool",
		defaultArgs: []string{"-v0", "-progress"},
		execCommand: mockExecCommand,
	}
	tool.NewArgs()
	tool.args.add("-test", "value")

	// Execute
	tool.Exec(false)

	// Verify command and arguments
	assert.Equal(t, "test-exiftool", capturedCmd, "Incorrect command executed")
	assert.Equal(t, []string{"-v0", "-progress", "-test", "value"}, capturedArgs, "Incorrect arguments passed")
}

func TestExifToolArgs_Recursively(t *testing.T) {
	sut := newExifTool().NewArgs()
	assert.Len(t, sut.args, 2)

	sut.Recursively(true)

	assert.Len(t, sut.args, 3)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-r", sut.args[2])
}

func TestExifToolArgs_NonRecursively(t *testing.T) {
	sut := newExifTool().NewArgs()
	assert.Len(t, sut.args, 2)

	sut.Recursively(false)

	assert.Len(t, sut.args, 2)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
}

func TestExifToolArgs_Src(t *testing.T) {
	sut := newExifTool().NewArgs()
	assert.Len(t, sut.args, 2)

	sut.Src("some_path")

	assert.Len(t, sut.args, 3)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "some_path", sut.args[2])
}

func TestExifToolArgs_ForImages(t *testing.T) {
	sut := newExifTool().NewArgs()

	sut.ForImages()

	assert.Len(t, sut.args, 12)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-ext", sut.args[2])
	assert.Equal(t, "jpg", sut.args[3])
	assert.Equal(t, "-ext", sut.args[4])
	assert.Equal(t, "jpeg", sut.args[5])
	assert.Equal(t, "-ext", sut.args[6])
	assert.Equal(t, "nef", sut.args[7])
	assert.Equal(t, "-ext", sut.args[8])
	assert.Equal(t, "cr2", sut.args[9])
	assert.Equal(t, "-ext", sut.args[10])
	assert.Equal(t, "cr3", sut.args[11])
}

func TestExifToolArgs_ForVideoMp4(t *testing.T) {
	sut := newExifTool().NewArgs()

	sut.ForVideoMp4()

	assert.Len(t, sut.args, 4)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-ext", sut.args[2])
	assert.Equal(t, "mp4", sut.args[3])
}

func TestExifToolArgs_ForVideoLrv(t *testing.T) {
	sut := newExifTool().NewArgs()

	sut.ForVideoLrv()

	assert.Len(t, sut.args, 4)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-ext", sut.args[2])
	assert.Equal(t, "LRV", sut.args[3])
}

func TestExifToolArgs_ForVideoAvchd(t *testing.T) {
	sut := newExifTool().NewArgs()

	sut.ForVideoAvchd()

	assert.Len(t, sut.args, 4)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-ext", sut.args[2])
	assert.Equal(t, "mts", sut.args[3])
}

func TestExifToolArgs_ForDateFormat(t *testing.T) {
	sut := newExifTool().NewArgs()

	assert.Len(t, sut.args, 2)

	sut.ForDateFormat("YYYY:MM:DD HH:MM:SS")

	assert.Len(t, sut.args, 4) // 2 default + 2 date format
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-d", sut.args[2])
	assert.Equal(t, "YYYY:MM:DD HH:MM:SS", sut.args[3])
}

func TestExifToolArgs_CopyTag(t *testing.T) {
	sut := newExifTool().NewArgs()
	assert.Len(t, sut.args, 2)

	sut.CopyTag("DateTime", "someOtherTag")

	assert.Len(t, sut.args, 3)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-DateTime<someOtherTag", sut.args[2])
}

func TestExifToolArgs_SetTag(t *testing.T) {
	sut := newExifTool().NewArgs()
	assert.Len(t, sut.args, 2)

	sut.SetTag("DateTime", "123")

	assert.Len(t, sut.args, 3)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-DateTime=123", sut.args[2])
}

func TestExifToolArgs_ChangeFileDate(t *testing.T) {
	sut := newExifTool().NewArgs()
	assert.Len(t, sut.args, 2)

	sut.ChangeFileDate("someOtherTag")

	assert.Len(t, sut.args, 4)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-FileModifyDate<someOtherTag", sut.args[2])
	assert.Equal(t, "-FileCreateDate<someOtherTag", sut.args[3])
}

func TestExifToolArgs_ChangeExifDate(t *testing.T) {
	sut := newExifTool().NewArgs()
	assert.Len(t, sut.args, 2)

	sut.ChangeExifDate("2023:01:01 12:00:00")

	assert.Len(t, sut.args, 4)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-CreateDate<2023:01:01 12:00:00", sut.args[2])
	assert.Equal(t, "-DateTimeOriginal<2023:01:01 12:00:00", sut.args[3])
}

func TestExifToolArgs_ChangeMp4Date(t *testing.T) {
	sut := newExifTool().NewArgs()
	assert.Len(t, sut.args, 2)

	sut.ChangeMp4Date("2023:01:01 12:00:00")

	assert.Len(t, sut.args, 8) // 2 default + 6 changeMp4Date (multiple date fields)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-CreateDate<2023:01:01 12:00:00", sut.args[2])
	assert.Equal(t, "-ModifyDate<2023:01:01 12:00:00", sut.args[3])
	assert.Equal(t, "-TrackCreateDate<2023:01:01 12:00:00", sut.args[4])
	assert.Equal(t, "-TrackModifyDate<2023:01:01 12:00:00", sut.args[5])
	assert.Equal(t, "-MediaCreateDate<2023:01:01 12:00:00", sut.args[6])
	assert.Equal(t, "-MediaModifyDate<2023:01:01 12:00:00", sut.args[7])
}

func TestExifToolArgs_CleanTag(t *testing.T) {
	sut := newExifTool().NewArgs()
	assert.Len(t, sut.args, 2)

	sut.CleanTag("DateTime")

	assert.Len(t, sut.args, 3)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-DateTime=", sut.args[2])
}

func TestExifToolArgs_CleanVendorTags(t *testing.T) {
	sut := newExifTool().NewArgs()
	assert.Len(t, sut.args, 2)

	sut.CleanVendorTags()

	assert.Len(t, sut.args, 12)
	commonTags := []string{"-Software=", "-WriterName=", "-ReaderName="}
	for _, tag := range commonTags {
		found := false
		for _, arg := range sut.args {
			if arg == tag {
				found = true
				break
			}
		}
		assert.True(t, found, "Should contain vendor tag: %s", tag)
	}
}

func TestExifToolArgs_CleanCameraTags(t *testing.T) {
	sut := newExifTool().NewArgs()

	assert.Len(t, sut.args, 2)

	sut.CleanCameraTags()

	assert.Greater(t, len(sut.args), 30)
	tagsToCheck := []string{"-Canon:all=", "-Sony:all=", "-Nikon:all="}
	for _, tag := range tagsToCheck {
		found := false
		for _, arg := range sut.args {
			if arg == tag {
				found = true
				break
			}
		}
		assert.True(t, found, "Should contain camera tag: %s", tag)
	}
}

func TestExifToolArgs_CleanLocationTags(t *testing.T) {
	sut := newExifTool().NewArgs()

	assert.Len(t, sut.args, 2)

	sut.CleanLocationTags()

	// 2 defaults and 1 extra
	assert.Len(t, sut.args, 3)

	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-gps:all=", sut.args[2])
}

func TestExcludeIfNameContains_NoArgs(t *testing.T) {
	sut := newExifTool().NewArgs()

	assert.Len(t, sut.args, 2)

	sut.ExcludeIfNameContains()

	// 2 defaults and no extra
	assert.Len(t, sut.args, 2)

	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
}

func TestExcludeIfNameContains_SingleArg(t *testing.T) {
	sut := newExifTool().NewArgs()

	assert.Len(t, sut.args, 2)

	sut.ExcludeIfNameContains("a")

	// 2 defaults and 2 extra
	assert.Len(t, sut.args, 4)

	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-if", sut.args[2])
	assert.Equal(t, "$filename !~ /a/i", sut.args[3])
}

func TestExcludeIfNameContains_MultipleArgs(t *testing.T) {
	sut := newExifTool().NewArgs()

	assert.Len(t, sut.args, 2)

	sut.ExcludeIfNameContains("a", "b")

	// 2 defaults and 2 extra
	assert.Len(t, sut.args, 4)

	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-if", sut.args[2])
	assert.Equal(t, "$filename !~ /a|b/i", sut.args[3])
}

func TestExcludeWhatsAppFiles_NormalCase(t *testing.T) {
	sut := newExifTool().NewArgs()

	assert.Len(t, sut.args, 2)

	sut.ExcludeWhatsAppFiles()

	// 2 defaults and 2 extra
	assert.Len(t, sut.args, 4)

	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-if", sut.args[2])
	assert.Equal(t, "$filename !~ /WhatsApp/i", sut.args[3])
}

func TestIncludeWhatsAppFiles_NormalCase(t *testing.T) {
	sut := newExifTool().NewArgs()

	assert.Len(t, sut.args, 2)

	sut.IncludeWhatsAppFiles()

	// 2 defaults and 2 extra
	assert.Len(t, sut.args, 4)

	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-if", sut.args[2])
	assert.Equal(t, "$filename =~ /WhatsApp /i", sut.args[3])
}

func TestExcludeTranscodedVideoFiles_NormalCase(t *testing.T) {
	sut := newExifTool().NewArgs()

	assert.Len(t, sut.args, 2)

	sut.ExcludeTranscodedVideoFiles()

	// 2 defaults and 2 extra
	assert.Len(t, sut.args, 4)

	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-if", sut.args[2])
	assert.Equal(t, "$filename !~ /_x265/i", sut.args[3])
}

func TestExifToolWrapper_ComplexUseCase1(t *testing.T) {
	// Test a complex scenario with multiple operations
	sut := newExifTool().NewArgs()

	// Add multiple operations
	sut.Recursively(true)
	sut.Src("/path/to/files")
	sut.ForImages()
	sut.CopyTag("DateTime", "2023:01:01 12:00:00")
	sut.CleanVendorTags()

	// Verify all operations were added (2 default + 1 recursively + 1 src + 8 forImages + 1 changeTag + 10 cleanVendorTags)
	assert.Len(t, sut.args, 25)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-r", sut.args[2])
	assert.Equal(t, "/path/to/files", sut.args[3])
	assert.Equal(t, "-ext", sut.args[4])
	assert.Equal(t, "jpg", sut.args[5])

	found := false
	for _, arg := range sut.args {
		if arg == "-DateTime<2023:01:01 12:00:00" {
			found = true
			break
		}
	}
	assert.True(t, found, "Should contain DateTime tag")
}

func TestExifToolWrapper_ComplexUseCase2(t *testing.T) {
	tool := newExifTool()
	sut := tool.NewArgs()

	sut.add("-v0", "-progress")
	sut.Recursively(true)
	sut.Src("/test/path")
	sut.ForImages()
	sut.CopyTag("DateTime", "2023:01:01 12:00:00")

	// Verify the command structure
	assert.Equal(t, "exiftool", tool.cmd)
	assert.Equal(t, []string{"-v0", "-progress"}, tool.defaultArgs)
	assert.Len(t, sut.args, 17)
}
