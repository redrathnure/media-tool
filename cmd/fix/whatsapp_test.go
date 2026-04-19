package fix

import (
	"testing"

	"github.com/redrathnure/media-tool/core/tools"
	"github.com/redrathnure/media-tool/core/tools/testutil"
	"github.com/stretchr/testify/assert"
)

func TestWhatsAppCmd_Initialization(t *testing.T) {
	assert.NotNil(t, whatsAppCmd)
	assert.NotNil(t, whatsAppCmd.Use)
	assert.NotNil(t, whatsAppCmd.Short)
	assert.NotNil(t, whatsAppCmd.Long)
}

func TestWhatsAppCmd_Flags(t *testing.T) {
	persistentFlags := whatsAppCmd.PersistentFlags()
	assert.False(t, persistentFlags.HasFlags())

	flags := whatsAppCmd.Flags()
	assert.False(t, flags.HasFlags())
}

func TestWhatsAppCmd_CommandStructure(t *testing.T) {
	assert.Equal(t, "whatsapp", whatsAppCmd.Name())
	assert.Equal(t, "media-tool fix whatsapp", whatsAppCmd.CommandPath())

	flags := whatsAppCmd.Flags()
	assert.Equal(t, 0, flags.NFlag())
}

func TestWhatsAppCmd_ArgValidation(t *testing.T) {
	// Test no args (should fail)
	err := whatsAppCmd.Args(whatsAppCmd, []string{})
	assert.NoError(t, err)

	// Test single arg (should pass)
	err = whatsAppCmd.Args(whatsAppCmd, []string{"test.jpg"})
	assert.NoError(t, err)

	// Test multiple args (should fail)
	err = whatsAppCmd.Args(whatsAppCmd, []string{"test1.jpg", "test2.jpg"})
	assert.Error(t, err)
}

func TestRunFixWhatsAppFiles(t *testing.T) {
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
				{Expected: []string{
					"-FileModifyDate<filename",
					"-CreateDate<filename",
					"-TrackModifyDate<filename",

					"-Exif:ProcessingSoftware=WhatsApp",
					"-Exif:Software=WhatsApp",
					"-Exif:MetadataEditingSoftware=WhatsApp",

					"-if $filename =~ /WhatsApp /i",
					"-if $filename !~ /_x265/i",
				},
					Unexpected: []string{
						"-r",
					}},
				{
					Expected: []string{
						"-FileName<CreateDate",

						"-if $filename =~ /WhatsApp /i",
						"-if $filename !~ /_x265/i",

						"-d IMG_%Y%m%d_%H%M%S_WhatsApp%%-c.%%e",
					},
					Unexpected: []string{
						"-r",
					},
				},
				{
					Expected: []string{
						"-FileName<CreateDate",

						"-if $filename =~ /WhatsApp /i",
						"-if $filename !~ /_x265/i",

						"-d VID_%Y%m%d_%H%M%S_WhatsApp%%-c.%%e",
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
				{Expected: []string{
					"-FileModifyDate<filename",
					"-CreateDate<filename",
					"-TrackModifyDate<filename",

					"-Exif:ProcessingSoftware=WhatsApp",
					"-Exif:Software=WhatsApp",
					"-Exif:MetadataEditingSoftware=WhatsApp",

					"-if $filename =~ /WhatsApp /i",
					"-if $filename !~ /_x265/i",
					"-r",
				}},
				{
					Expected: []string{
						"-FileName<CreateDate",

						"-if $filename =~ /WhatsApp /i",
						"-if $filename !~ /_x265/i",

						"-d IMG_%Y%m%d_%H%M%S_WhatsApp%%-c.%%e",
						"-r",
					},
				},
				{
					Expected: []string{
						"-FileName<CreateDate",

						"-if $filename =~ /WhatsApp /i",
						"-if $filename !~ /_x265/i",

						"-d VID_%Y%m%d_%H%M%S_WhatsApp%%-c.%%e",
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
			runFixWhatsAppFiles(whatsAppCmd, tt.Args)

			// Check ExifTool calls
			testTool.AssertCalls(t, 3)

			testTool.AssertCallArgs(t, 0, tt.Args[0], tt.CallArgs[0])
			testTool.AssertCallArgs(t, 1, tt.Args[0], tt.CallArgs[1])
			testTool.AssertCallArgs(t, 2, tt.Args[0], tt.CallArgs[2])
		})
	}
}

func TestRunFixWhatsAppFiles_DryRun(t *testing.T) {

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
						"-WriteNothing<filename",

						"-if $filename =~ /WhatsApp /i",
						"-if $filename !~ /_x265/i",
					},
					Unexpected: []string{
						"-r",
					},
				},
				{
					Expected: []string{
						"-TestName<CreateDate",

						"-if $filename =~ /WhatsApp /i",
						"-if $filename !~ /_x265/i",

						"-d IMG_%Y%m%d_%H%M%S_WhatsApp%%-c.%%e",
					},
					Unexpected: []string{
						"-FileModifyDate<CreateDate",
						"-r",
					},
				},
				{
					Expected: []string{
						"-TestName<CreateDate",

						"-if $filename =~ /WhatsApp /i",
						"-if $filename !~ /_x265/i",

						"-d VID_%Y%m%d_%H%M%S_WhatsApp%%-c.%%e",
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
						"-WriteNothing<filename",

						"-if $filename =~ /WhatsApp /i",
						"-if $filename !~ /_x265/i",
						"-r",
					}},
				{
					Expected: []string{
						"-TestName<CreateDate",

						"-if $filename =~ /WhatsApp /i",
						"-if $filename !~ /_x265/i",

						"-d IMG_%Y%m%d_%H%M%S_WhatsApp%%-c.%%e",
						"-r",
					},
					Unexpected: []string{
						"-FileModifyDate<CreateDate",
					},
				},
				{
					Expected: []string{
						"-TestName<CreateDate",

						"-if $filename =~ /WhatsApp /i",
						"-if $filename !~ /_x265/i",

						"-d VID_%Y%m%d_%H%M%S_WhatsApp%%-c.%%e",
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
			runFixWhatsAppFiles(whatsAppCmd, tt.Args)

			// Check ExifTool calls
			testTool.AssertCalls(t, 3)

			testTool.AssertCallArgs(t, 0, tt.Args[0], tt.CallArgs[0])
			testTool.AssertCallArgs(t, 1, tt.Args[0], tt.CallArgs[1])
			testTool.AssertCallArgs(t, 2, tt.Args[0], tt.CallArgs[2])
		})
	}
}

func TestFixWhatsAppFiles_DryRun(t *testing.T) {
	tmp := testutil.NewTempDir(t)
	defer tmp.Clean()

	tmp.MkTestMedia("WhatsApp Image 2026-04-15 at 15.57.29", "src")

	srcDir := tmp.DirName("src")

	fixWhatsAppFiles(srcDir, true, true, true)

	assert.True(t, tmp.IsDirExist("src"))

	assert.True(t, tmp.IsFileExist("WhatsApp Image 2026-04-15 at 15.57.29.jpg", "src"))
	assert.True(t, tmp.IsFileExist("WhatsApp Image 2026-04-15 at 15.57.29.jpeg", "src"))

	assert.True(t, tmp.IsFileExist("WhatsApp Image 2026-04-15 at 15.57.29.mkv", "src"))
	assert.True(t, tmp.IsFileExist("WhatsApp Image 2026-04-15 at 15.57.29.mp4", "src"))
}

func TestFixWhatsAppFiles_Fix(t *testing.T) {
	tmp := testutil.NewTempDir(t)
	defer tmp.Clean()

	tmp.MkTestMedia("WhatsApp Image 2026-01-02 at 03.04.05", "src")

	srcDir := tmp.DirName("src")

	fixWhatsAppFiles(srcDir, true, false, true)

	assert.True(t, tmp.IsDirExist("src"))

	expectedDate := "2026-01-02 03:04:05"

	assert.True(t, tmp.IsFileExist("IMG_20260102_030405_WhatsApp.jpg", "src"))
	assert.Equal(t, expectedDate, tmp.FileDate("IMG_20260102_030405_WhatsApp.jpg", "src"))
	assert.True(t, tmp.IsFileExist("IMG_20260102_030405_WhatsApp.jpeg", "src"))
	assert.Equal(t, expectedDate, tmp.FileDate("IMG_20260102_030405_WhatsApp.jpeg", "src"))

	//Warning: mkv is not fully supported by exiftool
	assert.False(t, tmp.IsFileExist("VID_20250612_184641.mkv", "src"))
	assert.True(t, tmp.IsFileExist("WhatsApp Image 2026-01-02 at 03.04.05.mkv", "src"))
	assert.True(t, tmp.IsFileExist("VID_20260102_030405_WhatsApp.mp4", "src"))
	assert.Equal(t, expectedDate, tmp.FileDate("VID_20260102_030405_WhatsApp.mp4", "src"))
}

func TestFixWhatsAppFiles_ExcludeNonWhatsAppFiles(t *testing.T) {
	tmp := testutil.NewTempDir(t)
	defer tmp.Clean()

	tmp.MkTestMedia("IMG_20260102_030405", "src")

	srcDir := tmp.DirName("src")

	fixWhatsAppFiles(srcDir, true, false, true)

	assert.True(t, tmp.IsDirExist("src"))

	expectedDate := "2026-01-02 03:04:05"

	assert.True(t, tmp.IsFileExist("IMG_20260102_030405.jpg", "src"))
	assert.NotEqual(t, expectedDate, tmp.FileDate("IMG_20260102_030405.jpg", "src"))
	assert.True(t, tmp.IsFileExist("IMG_20260102_030405.jpeg", "src"))
	assert.NotEqual(t, expectedDate, tmp.FileDate("IMG_20260102_030405.jpeg", "src"))
	assert.True(t, tmp.IsFileExist("IMG_20260102_030405.mkv", "src"))
	assert.NotEqual(t, expectedDate, tmp.FileDate("IMG_20260102_030405.mkv", "src"))
	assert.True(t, tmp.IsFileExist("IMG_20260102_030405.mp4", "src"))
	assert.NotEqual(t, expectedDate, tmp.FileDate("IMG_20260102_030405.mp4", "src"))
}

func TestFixWhatsAppFiles_ExcludeAlreadyProcessedFiles(t *testing.T) {
	tmp := testutil.NewTempDir(t)
	defer tmp.Clean()

	tmp.MkTestMedia("IMG_20260102_030405_WhatsApp", "src")

	srcDir := tmp.DirName("src")

	fixWhatsAppFiles(srcDir, true, false, true)

	assert.True(t, tmp.IsDirExist("src"))

	expectedDate := "2026-01-02 03:04:05"

	assert.True(t, tmp.IsFileExist("IMG_20260102_030405_WhatsApp.jpg", "src"))
	assert.NotEqual(t, expectedDate, tmp.FileDate("IMG_20260102_030405_WhatsApp.jpg", "src"))
	assert.True(t, tmp.IsFileExist("IMG_20260102_030405_WhatsApp.jpeg", "src"))
	assert.NotEqual(t, expectedDate, tmp.FileDate("IMG_20260102_030405_WhatsApp.jpeg", "src"))
	assert.True(t, tmp.IsFileExist("IMG_20260102_030405_WhatsApp.mkv", "src"))
	assert.NotEqual(t, expectedDate, tmp.FileDate("IMG_20260102_030405_WhatsApp.mkv", "src"))
	assert.True(t, tmp.IsFileExist("IMG_20260102_030405_WhatsApp.mp4", "src"))
	assert.NotEqual(t, expectedDate, tmp.FileDate("IMG_20260102_030405_WhatsApp.mp4", "src"))
}

func TestFixWhatsAppFiles_BadFileName(t *testing.T) {
	tmp := testutil.NewTempDir(t)
	defer tmp.Clean()

	//No information about dateime in file names
	tmp.MkTestMedia("WhatApp zxc", "src")

	srcDir := tmp.DirName("src")

	fixWhatsAppFiles(srcDir, true, false, true)

	assert.True(t, tmp.IsDirExist("src"))

	assert.True(t, tmp.IsFileExist("WhatApp zxc.jpg", "src"))
	assert.True(t, tmp.IsFileExist("WhatApp zxc.jpeg", "src"))
	assert.True(t, tmp.IsFileExist("WhatApp zxc.mkv", "src"))
	assert.True(t, tmp.IsFileExist("WhatApp zxc.mp4", "src"))
}
