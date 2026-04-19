package import_

import (
	"testing"

	"github.com/redrathnure/media-tool/core/tools/testutil"
	"github.com/stretchr/testify/assert"
)

func TestSdPhotosCmd_Initialization(t *testing.T) {
	assert.NotNil(t, sdPhotosCmd)
	assert.NotNil(t, sdPhotosCmd.Use)
	assert.NotNil(t, sdPhotosCmd.Short)
	assert.NotNil(t, sdPhotosCmd.Long)
}

func TestSdPhotosCmd_Flags(t *testing.T) {
	persistentFlags := sdPhotosCmd.PersistentFlags()
	assert.False(t, persistentFlags.HasFlags())

	flags := sdPhotosCmd.Flags()
	assert.True(t, flags.HasFlags())

	localRenameFlag := flags.Lookup("keep-names")
	assert.NotNil(t, localRenameFlag)
	assert.Equal(t, "k", localRenameFlag.Shorthand)
	assert.Equal(t, "false", localRenameFlag.DefValue)
	assert.Equal(t, "bool", localRenameFlag.Value.Type())
}

func TestSdPhotosCmd_CommandStructure(t *testing.T) {
	assert.Equal(t, "sdphotos", sdPhotosCmd.Name())
	assert.Equal(t, "media-tool import sdphotos", sdPhotosCmd.CommandPath())

	flags := sdPhotosCmd.Flags()
	assert.Equal(t, 0, flags.NFlag())
}

func TestSdPhotosCmd_ArgValidation(t *testing.T) {
	// Test no args
	err := sdPhotosCmd.Args(sdPhotosCmd, []string{})
	assert.NoError(t, err)

	// Test single arg
	err = sdPhotosCmd.Args(sdPhotosCmd, []string{"test.jpg"})
	assert.NoError(t, err)

	// Test multiple args
	err = sdPhotosCmd.Args(sdPhotosCmd, []string{"test1.jpg", "test2.jpg"})
	assert.Error(t, err)
}

func TestSdPhotosMoveToDst_DryRun(t *testing.T) {
	tmp := testutil.NewTempDir(t)
	defer tmp.Clean()

	tmp.MkTestMedia("zxc", "src")

	srcDir := tmp.DirName("src")
	resultDir := tmp.DirName("results")

	sdPhotosMoveToDst(srcDir, resultDir, false, true)

	assert.False(t, tmp.IsDirExist("results"))

	assert.True(t, tmp.IsFileExist("zxc.cr3", "src"))
	assert.True(t, tmp.IsFileExist("zxc.jpg", "src"))
	assert.True(t, tmp.IsFileExist("zxc.jpeg", "src"))
	assert.True(t, tmp.IsFileExist("zxc.mkv", "src"))
	assert.True(t, tmp.IsFileExist("zxc.mp4", "src"))
}

func TestSdPhotosMoveToDst_Rename(t *testing.T) {
	tmp := testutil.NewTempDir(t)
	defer tmp.Clean()

	tmp.MkTestMedia("zxc", "src")

	srcDir := tmp.DirName("src")
	resultDir := tmp.DirName("results")

	sdPhotosMoveToDst(srcDir, resultDir, false, false)

	assert.True(t, tmp.IsDirExist("results"))

	assert.True(t, tmp.IsFileExist("IMG_20250612_184641.cr3", "results", "2025.06.12"))
	assert.Equal(t, testutil.CreationDate, tmp.FileDate("IMG_20250612_184641.cr3", "results", "2025.06.12"))

	assert.True(t, tmp.IsFileExist("IMG_20250612_184641.jpg", "results", "2025.06.12"))
	assert.Equal(t, testutil.CreationDate, tmp.FileDate("IMG_20250612_184641.jpg", "results", "2025.06.12"))

	assert.True(t, tmp.IsFileExist("IMG_20250612_184641.jpeg", "results", "2025.06.12"))
	assert.Equal(t, testutil.CreationDate, tmp.FileDate("IMG_20250612_184641.jpeg", "results", "2025.06.12"))

	//Warning: mkv is not fully supported by exiftool
	assert.False(t, tmp.IsFileExist("VID_20250612_184641.mkv", "results", "2025.06.12"))
	assert.True(t, tmp.IsFileExist("VID_20250612_184641.mp4", "results", "2025.06.12"))
	assert.Equal(t, testutil.CreationDate, tmp.FileDate("VID_20250612_184641.mp4", "results", "2025.06.12"))
}

func TestSdPhotosMoveToDst_KeepNames(t *testing.T) {
	tmp := testutil.NewTempDir(t)
	defer tmp.Clean()

	tmp.MkTestMedia("zxc", "src")

	srcDir := tmp.DirName("src")
	resultDir := tmp.DirName("results")

	sdPhotosMoveToDst(srcDir, resultDir, true, false)

	assert.True(t, tmp.IsDirExist("results"))

	assert.True(t, tmp.IsFileExist("zxc.cr3", "results", "2025.06.12"))
	assert.Equal(t, testutil.CreationDate, tmp.FileDate("zxc.cr3", "results", "2025.06.12"))

	assert.True(t, tmp.IsFileExist("zxc.jpg", "results", "2025.06.12"))
	assert.Equal(t, testutil.CreationDate, tmp.FileDate("zxc.jpg", "results", "2025.06.12"))

	assert.True(t, tmp.IsFileExist("zxc.jpeg", "results", "2025.06.12"))
	assert.Equal(t, testutil.CreationDate, tmp.FileDate("zxc.jpeg", "results", "2025.06.12"))

	//Warning: mkv is not fully supported by exiftool
	assert.False(t, tmp.IsFileExist("zxc.mkv", "results", "2025.06.12"))
	assert.True(t, tmp.IsFileExist("zxc.mp4", "results", "2025.06.12"))
	assert.Equal(t, testutil.CreationDate, tmp.FileDate("zxc.mp4", "results", "2025.06.12"))
}
