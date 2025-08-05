package cmd

import (
	"os/exec"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewExifTool(t *testing.T) {
	sut := newExifTool()

	assert.NotNil(t, sut)
	assert.Equal(t, "exiftool", sut.cmd)
	assert.Equal(t, []string{"-v0", "-progress"}, sut.defaultArgs)
}

func TestGetExifTool_Singleton(t *testing.T) {
	sut1 := getExifTool()
	sut := getExifTool()

	assert.NotNil(t, sut1)
	assert.NotNil(t, sut)
	assert.Equal(t, sut1, sut)
}

func TestExifToolWrapper_InitCmd_Default(t *testing.T) {
	sut := &exifToolWrapper{
		cmd:         "exiftool",
		defaultArgs: []string{"-v0", "-progress"},
	}

	// Reset viper to ensure clean state
	viper.Reset()

	sut.initCmd()

	assert.Equal(t, "exiftool", sut.cmd)
}

func TestExifToolWrapper_InitCmd_CustomPath(t *testing.T) {
	// Test initCmd with custom path
	sut := &exifToolWrapper{
		cmd:         "exiftool",
		defaultArgs: []string{"-v0", "-progress"},
	}

	// Reset viper to ensure clean state
	viper.Reset()

	viper.Set(cfgExifToolPath, "/custom/path/exiftool")

	sut.initCmd()

	assert.Equal(t, "/custom/path/exiftool", sut.cmd)

	// Reset viper and singleton for next tests
	viper.Reset()
	exifToolObj = nil
}

func TestExifToolWrapper_InitCmd_AppDirSubstitution(t *testing.T) {
	// Test initCmd with $APP_DIR substitution
	sut := &exifToolWrapper{
		cmd:         "exiftool",
		defaultArgs: []string{"-v0", "-progress"},
	}

	// Reset viper to ensure clean state
	viper.Reset()

	// Set app directory path (this matches the name used in the exiftool.go implementation)
	viper.Set(cfgExifToolPath, "$APP_DIR/..")

	sut.initCmd()

	// Verify $APP_DIR was replaced with the configured path
	assert.NotEqual(t, "$APP_DIR/..", sut.cmd)

	// Reset viper and singleton for next tests
	viper.Reset()
	exifToolObj = nil
}

func TestExifToolArgs_NewArgs(t *testing.T) {
	tool := newExifTool()
	sut := tool.newArgs()

	assert.Equal(t, &tool.args, sut)

	assert.NotNil(t, sut)
	assert.NotNil(t, sut.args)
	assert.Len(t, sut.args, 2)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
}

func TestExifToolArgs_Add(t *testing.T) {
	sut := newExifTool().newArgs()
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

	tool := &exifToolWrapper{
		cmd:         "test-exiftool",
		defaultArgs: []string{"-v0", "-progress"},
		execCommand: mockExecCommand,
	}
	tool.newArgs()
	tool.args.add("-test", "value")

	// Execute
	tool.exec()

	// Verify command and arguments
	assert.Equal(t, "test-exiftool", capturedCmd, "Incorrect command executed")
	assert.Equal(t, []string{"-v0", "-progress", "-test", "value"}, capturedArgs, "Incorrect arguments passed")
}

func TestExifToolArgs_Recursively(t *testing.T) {
	sut := newExifTool().newArgs()
	assert.Len(t, sut.args, 2)

	sut.recursively()

	assert.Len(t, sut.args, 3)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-r", sut.args[2])
}

func TestExifToolArgs_Src(t *testing.T) {
	sut := newExifTool().newArgs()
	assert.Len(t, sut.args, 2)

	sut.src("some_path")

	assert.Len(t, sut.args, 3)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "some_path", sut.args[2])
}

func TestExifToolArgs_ForImages(t *testing.T) {
	sut := newExifTool().newArgs()

	sut.forImages()

	assert.Len(t, sut.args, 10)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-ext", sut.args[2])
	assert.Equal(t, "jpg", sut.args[3])
	assert.Equal(t, "-ext", sut.args[4])
	assert.Equal(t, "nef", sut.args[5])
	assert.Equal(t, "-ext", sut.args[6])
	assert.Equal(t, "cr2", sut.args[7])
	assert.Equal(t, "-ext", sut.args[8])
	assert.Equal(t, "cr3", sut.args[9])
}

func TestExifToolArgs_ForVideoMp4(t *testing.T) {
	sut := newExifTool().newArgs()

	sut.forVideoMp4()

	assert.Len(t, sut.args, 4)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-ext", sut.args[2])
	assert.Equal(t, "mp4", sut.args[3])
}

func TestExifToolArgs_ForVideoLrv(t *testing.T) {
	sut := newExifTool().newArgs()

	sut.forVideoLrv()

	assert.Len(t, sut.args, 4)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-ext", sut.args[2])
	assert.Equal(t, "LRV", sut.args[3])
}

func TestExifToolArgs_ForVideoAvchd(t *testing.T) {
	sut := newExifTool().newArgs()

	sut.forVideoAvchd()

	assert.Len(t, sut.args, 4)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-ext", sut.args[2])
	assert.Equal(t, "mts", sut.args[3])
}

func TestExifToolArgs_ForDateFormat(t *testing.T) {
	sut := newExifTool().newArgs()

	assert.Len(t, sut.args, 2)

	sut.forDateFormat("YYYY:MM:DD HH:MM:SS")

	assert.Len(t, sut.args, 4) // 2 default + 2 date format
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-d", sut.args[2])
	assert.Equal(t, "YYYY:MM:DD HH:MM:SS", sut.args[3])
}

func TestExifToolArgs_ChangeTag(t *testing.T) {
	sut := newExifTool().newArgs()
	assert.Len(t, sut.args, 2)

	sut.changeTag("DateTime", "2023:01:01 12:00:00")

	assert.Len(t, sut.args, 3)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-DateTime<2023:01:01 12:00:00", sut.args[2])
}

func TestExifToolArgs_ChangeFileDate(t *testing.T) {
	sut := newExifTool().newArgs()
	assert.Len(t, sut.args, 2)

	sut.changeFileDate("2023:01:01 12:00:00")

	assert.Len(t, sut.args, 4)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-FileModifyDate<2023:01:01 12:00:00", sut.args[2])
	assert.Equal(t, "-FileCreateDate<2023:01:01 12:00:00", sut.args[3])
}

func TestExifToolArgs_ChangeExifDate(t *testing.T) {
	sut := newExifTool().newArgs()
	assert.Len(t, sut.args, 2)

	sut.changeExifDate("2023:01:01 12:00:00")

	assert.Len(t, sut.args, 4)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-CreateDate<2023:01:01 12:00:00", sut.args[2])
	assert.Equal(t, "-DateTimeOriginal<2023:01:01 12:00:00", sut.args[3])
}

func TestExifToolArgs_ChangeMp4Date(t *testing.T) {
	sut := newExifTool().newArgs()
	assert.Len(t, sut.args, 2)

	sut.changeMp4Date("2023:01:01 12:00:00")

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
	sut := newExifTool().newArgs()
	assert.Len(t, sut.args, 2)

	sut.cleanTag("DateTime")

	assert.Len(t, sut.args, 3)
	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-DateTime=", sut.args[2])
}

func TestExifToolArgs_CleanVendorTags(t *testing.T) {
	sut := newExifTool().newArgs()
	assert.Len(t, sut.args, 2)

	sut.cleanVendorTags()

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
	sut := newExifTool().newArgs()

	assert.Len(t, sut.args, 2)

	sut.cleanCameraTags()

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
	sut := newExifTool().newArgs()

	assert.Len(t, sut.args, 2)

	sut.cleanLocationTags()

	// 2 defaults and 1 extra
	assert.Len(t, sut.args, 3)

	assert.Equal(t, "-v0", sut.args[0])
	assert.Equal(t, "-progress", sut.args[1])
	assert.Equal(t, "-gps:all=", sut.args[2])
}

func TestExifToolWrapper_ComplexUseCase1(t *testing.T) {
	// Test a complex scenario with multiple operations
	sut := newExifTool().newArgs()

	// Add multiple operations
	sut.recursively()
	sut.src("/path/to/files")
	sut.forImages()
	sut.changeTag("DateTime", "2023:01:01 12:00:00")
	sut.cleanVendorTags()

	// Verify all operations were added (2 default + 1 recursively + 1 src + 8 forImages + 1 changeTag + 10 cleanVendorTags)
	assert.Len(t, sut.args, 23)
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
	sut := tool.newArgs()

	sut.add("-v0", "-progress")
	sut.recursively()
	sut.src("/test/path")
	sut.forImages()
	sut.changeTag("DateTime", "2023:01:01 12:00:00")

	// Verify the command structure
	assert.Equal(t, "exiftool", tool.cmd)
	assert.Equal(t, []string{"-v0", "-progress"}, tool.defaultArgs)
	assert.Len(t, sut.args, 15)
}
