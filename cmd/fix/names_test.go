package fix

import (
	"testing"

	"github.com/redrathnure/media-tool/core/tools"
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
