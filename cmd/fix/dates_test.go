package fix

import (
	"testing"

	"github.com/redrathnure/media-tool/core/tools"
	"github.com/stretchr/testify/assert"
)

func TestDatesCmd_Initialization(t *testing.T) {
	assert.NotNil(t, datesCmd)
	assert.NotNil(t, datesCmd.Use)
	assert.NotNil(t, datesCmd.Short)
	assert.NotNil(t, datesCmd.Long)
}

func TestDatesCmd_Flags(t *testing.T) {
	persistentFlags := datesCmd.PersistentFlags()
	assert.False(t, persistentFlags.HasFlags())

	flags := datesCmd.Flags()
	assert.False(t, flags.HasFlags())
}

func TestDatesCmd_CommandStructure(t *testing.T) {
	assert.Equal(t, "dates", datesCmd.Name())
	assert.Equal(t, "media-tool fix dates", datesCmd.CommandPath())

	flags := datesCmd.Flags()
	assert.Equal(t, 0, flags.NFlag())
}

func TestDatesCmd_ArgValidation(t *testing.T) {
	// Test no args (should fail)
	err := datesCmd.Args(datesCmd, []string{})
	assert.NoError(t, err)

	// Test single arg (should pass)
	err = datesCmd.Args(datesCmd, []string{"test.jpg"})
	assert.NoError(t, err)

	// Test multiple args (should fail)
	err = datesCmd.Args(datesCmd, []string{"test1.jpg", "test2.jpg"})
	assert.Error(t, err)
}

func TestRunFixDates(t *testing.T) {
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
				"-FileModifyDate<filename",
				"-CreateDate<filename",
				"-TrackModifyDate<filename",
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
				"-FileModifyDate<filename",
				"-CreateDate<filename",
				"-TrackModifyDate<filename",
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
			runFixDates(datesCmd, tt.Args)

			// Check ExifTool calls
			testTool.AssertCalls(t, 1)

			testTool.AssertCallArgs(t, 0, tt.Args[0], tt.CallArgs())
		})
	}
}

func TestRunFixDates_DryRun(t *testing.T) {
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
				"-WriteNothing<filename",
			},
			Unexpected: []string{
				"-FileModifyDate<filename",
				"-CreateDate<filename",
				"-TrackModifyDate<filename",
				"-r",
			},
		},
		{
			Name:      "with recursion",
			Args:      []string{"test.jpg"},
			Recursive: true,
			Expected: []string{
				"-WriteNothing<filename",
				"-r",
			},
			Unexpected: []string{
				"-FileModifyDate<filename",
				"-CreateDate<filename",
				"-TrackModifyDate<filename",
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
			runFixDates(datesCmd, tt.Args)

			// Check ExifTool calls
			testTool.AssertCalls(t, 1)

			testTool.AssertCallArgs(t, 0, tt.Args[0], tt.CallArgs())
		})
	}
}
