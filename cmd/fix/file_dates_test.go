package fix

import (
	"testing"

	"github.com/redrathnure/media-tool/core/tools"
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
