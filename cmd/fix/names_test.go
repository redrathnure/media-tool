package fix

import (
	"testing"

	"github.com/redrathnure/media-tool/core/tools"
	"github.com/redrathnure/media-tool/core/tools/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNamesCmd_Initialization(t *testing.T) {
	assert.NotNil(t, namesCmd)
	assert.NotNil(t, namesCmd.Use)
	assert.NotNil(t, namesCmd.Short)
	assert.NotNil(t, namesCmd.Long)
}

func TestNamesCmd_Flags(t *testing.T) {
	persistentFlags := namesCmd.PersistentFlags()
	assert.False(t, persistentFlags.HasFlags())

	flags := namesCmd.Flags()
	assert.False(t, flags.HasFlags())
}

func TestNamesCmd_CommandStructure(t *testing.T) {
	assert.Equal(t, "names", namesCmd.Name())
	assert.Equal(t, "media-tool fix names", namesCmd.CommandPath())

	flags := namesCmd.Flags()
	assert.Equal(t, 0, flags.NFlag())
}

func TestNamesCmd_ArgValidation(t *testing.T) {
	// Test no args (should fail)
	err := namesCmd.Args(namesCmd, []string{})
	assert.NoError(t, err)

	// Test single arg (should pass)
	err = namesCmd.Args(namesCmd, []string{"test.jpg"})
	assert.NoError(t, err)

	// Test multiple args (should fail)
	err = namesCmd.Args(namesCmd, []string{"test1.jpg", "test2.jpg"})
	assert.Error(t, err)
}

func TestRunFixNames(t *testing.T) {
	// Save original values to restore after test
	origRecursively := recursively
	origDryRun := dryRun
	defer func() {
		recursively = origRecursively
		dryRun = origDryRun
	}()

	tests := []tools.ExifMultiCallCaseData{
		{
			Name:      "without recursion",
			Args:      []string{"test.jpg"},
			Recursive: false,
			CallArgs: []tools.ExifRunArgs{
				{
					Expected: []string{
						"-FileName<CreateDate",
						"-if $filename !~ /WhatsApp/i",
						"-if $filename !~ /_x265/i",
						"-d IMG_%Y%m%d_%H%M%S%%-c.%%e",
					},
					Unexpected: []string{
						"-r",
					},
				},
				{
					Expected: []string{
						"-FileName<CreateDate",
						"-if $filename !~ /WhatsApp/i",
						"-if $filename !~ /_x265/i",
						"-d VID_%Y%m%d_%H%M%S%%-c.%%e",
					},
					Unexpected: []string{
						"-r",
					},
				},
			},
		},
		{
			Name:      "with recursion",
			Args:      []string{"test.jpg"},
			Recursive: true,
			CallArgs: []tools.ExifRunArgs{
				{
					Expected: []string{
						"-FileName<CreateDate",
						"-if $filename !~ /WhatsApp/i",
						"-if $filename !~ /_x265/i",
						"-d IMG_%Y%m%d_%H%M%S%%-c.%%e",
						"-r",
					},
				},
				{
					Expected: []string{
						"-FileName<CreateDate",
						"-if $filename !~ /WhatsApp/i",
						"-if $filename !~ /_x265/i",
						"-d VID_%Y%m%d_%H%M%S%%-c.%%e",
						"-r",
					},
				},
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
			runFixNames(namesCmd, tt.Args)

			// Check ExifTool calls
			testTool.AssertCalls(t, 2)

			testTool.AssertCallArgs(t, 0, tt.Args[0], tt.CallArgs[0])
			testTool.AssertCallArgs(t, 1, tt.Args[0], tt.CallArgs[1])
		})
	}
}

func TestRunFixNames_DryRun(t *testing.T) {
	// Save original values to restore after test
	origRecursively := recursively
	origDryRun := dryRun
	defer func() {
		recursively = origRecursively
		dryRun = origDryRun
	}()

	tests := []tools.ExifMultiCallCaseData{
		{
			Name:      "without recursion",
			Args:      []string{"test.jpg"},
			Recursive: false,
			CallArgs: []tools.ExifRunArgs{
				{
					Expected: []string{
						"-TestName<CreateDate",
						"-if $filename !~ /WhatsApp/i",
						"-if $filename !~ /_x265/i",
						"-d IMG_%Y%m%d_%H%M%S%%-c.%%e",
					},
					Unexpected: []string{
						"-FileModifyDate<CreateDate",
						"-r",
					},
				},
				{
					Expected: []string{
						"-TestName<CreateDate",
						"-if $filename !~ /WhatsApp/i",
						"-if $filename !~ /_x265/i",
						"-d VID_%Y%m%d_%H%M%S%%-c.%%e",
					},
					Unexpected: []string{
						"-FileModifyDate<CreateDate",
						"-r",
					},
				},
			},
		},
		{
			Name:      "with recursion",
			Args:      []string{"test.jpg"},
			Recursive: true,
			CallArgs: []tools.ExifRunArgs{
				{
					Expected: []string{
						"-TestName<CreateDate",
						"-if $filename !~ /WhatsApp/i",
						"-if $filename !~ /_x265/i",
						"-d IMG_%Y%m%d_%H%M%S%%-c.%%e",
						"-r",
					},
					Unexpected: []string{
						"-FileModifyDate<CreateDate",
					},
				},
				{
					Expected: []string{
						"-TestName<CreateDate",
						"-if $filename !~ /WhatsApp/i",
						"-if $filename !~ /_x265/i",
						"-d VID_%Y%m%d_%H%M%S%%-c.%%e",
						"-r",
					},
					Unexpected: []string{
						"-FileModifyDate<CreateDate",
					},
				},
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
			runFixNames(namesCmd, tt.Args)

			// Check ExifTool calls
			testTool.AssertCalls(t, 2)

			testTool.AssertCallArgs(t, 0, tt.Args[0], tt.CallArgs[0])
			testTool.AssertCallArgs(t, 1, tt.Args[0], tt.CallArgs[1])
		})
	}
}

func TestFixNames_DryRun(t *testing.T) {
	tmp := testutil.NewTempDir(t)
	defer tmp.Clean()

	tmp.MkTestMedia("zxc", "src")

	srcDir := tmp.DirName("src")

	fixFileNames(srcDir, true, true, true)

	assert.True(t, tmp.IsDirExist("src"))

	assert.True(t, tmp.IsFileExist("zxc.cr3", "src"))
	assert.NotEqual(t, testutil.CreationDate, tmp.FileDate("zxc.cr3", "src"))
	assert.True(t, tmp.IsFileExist("zxc.jpg", "src"))
	assert.NotEqual(t, testutil.CreationDate, tmp.FileDate("zxc.jpg", "src"))
	assert.True(t, tmp.IsFileExist("zxc.jpeg", "src"))
	assert.NotEqual(t, testutil.CreationDate, tmp.FileDate("zxc.jpeg", "src"))
	assert.True(t, tmp.IsFileExist("zxc.mkv", "src"))
	assert.NotEqual(t, testutil.CreationDate, tmp.FileDate("zxc.mkv", "src"))
	assert.True(t, tmp.IsFileExist("zxc.mp4", "src"))
	assert.NotEqual(t, testutil.CreationDate, tmp.FileDate("zxc.mp4", "src"))
}

func TestFixNames_Fix(t *testing.T) {
	tmp := testutil.NewTempDir(t)
	defer tmp.Clean()

	tmp.MkTestMedia("zxc", "src")

	srcDir := tmp.DirName("src")

	fixFileNames(srcDir, true, false, true)

	assert.True(t, tmp.IsDirExist("src"))

	//No any date correction, just rename files
	assert.True(t, tmp.IsFileExist("IMG_20250612_184641.cr3", "src"))
	assert.NotEqual(t, testutil.CreationDate, tmp.FileDate("IMG_20250612_184641.cr3", "src"))
	assert.True(t, tmp.IsFileExist("IMG_20250612_184641.jpg", "src"))
	assert.NotEqual(t, testutil.CreationDate, tmp.FileDate("IMG_20250612_184641.jpg", "src"))
	assert.True(t, tmp.IsFileExist("IMG_20250612_184641.jpeg", "src"))
	assert.NotEqual(t, testutil.CreationDate, tmp.FileDate("IMG_20250612_184641.jpeg", "src"))
	//Warning: mkv is not fully supported by exiftool
	assert.False(t, tmp.IsFileExist("VID_20250612_184641.mkv", "src"))
	assert.True(t, tmp.IsFileExist("zxc.mkv", "src"))
	assert.True(t, tmp.IsFileExist("VID_20250612_184641.mp4", "src"))
	assert.NotEqual(t, testutil.CreationDate, tmp.FileDate("VID_20250612_184641.mp4", "src"))
}

func TestFixNames_ExcludeWhatsAppFiles(t *testing.T) {
	tmp := testutil.NewTempDir(t)
	defer tmp.Clean()

	tmp.MkTestMedia("zxc-whatsapp", "src")

	srcDir := tmp.DirName("src")

	fixFileNames(srcDir, true, false, true)

	assert.True(t, tmp.IsDirExist("src"))

	assert.True(t, tmp.IsFileExist("zxc-whatsapp.cr3", "src"))
	assert.NotEqual(t, testutil.CreationDate, tmp.FileDate("zxc-whatsapp.cr3", "src"))
	assert.True(t, tmp.IsFileExist("zxc-whatsapp.jpg", "src"))
	assert.NotEqual(t, testutil.CreationDate, tmp.FileDate("zxc-whatsapp.jpg", "src"))
	assert.True(t, tmp.IsFileExist("zxc-whatsapp.jpeg", "src"))
	assert.NotEqual(t, testutil.CreationDate, tmp.FileDate("zxc-whatsapp.jpeg", "src"))
	assert.True(t, tmp.IsFileExist("zxc-whatsapp.mkv", "src"))
	assert.NotEqual(t, testutil.CreationDate, tmp.FileDate("zxc-whatsapp.mkv", "src"))
	assert.True(t, tmp.IsFileExist("zxc-whatsapp.mp4", "src"))
	assert.NotEqual(t, testutil.CreationDate, tmp.FileDate("zxc-whatsapp.mp4", "src"))
}
