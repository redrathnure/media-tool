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
	defer func() {
		recursively = origRecursively
	}()

	origDryRun := dryRun
	defer func() {
		dryRun = origDryRun
	}()

	tests := []struct {
		name           string
		args           []string
		recursive      bool
		expectedTags   []string
		unexpectedTags []string
	}{
		{
			name:      "without recursion",
			args:      []string{"test.jpg"},
			recursive: false,
			expectedTags: []string{
				"-FileModifyDate<filename",
				"-CreateDate<filename",
				"-TrackModifyDate<filename",
			},
			unexpectedTags: []string{
				"-r",
			},
		},
		{
			name:      "with recursion",
			args:      []string{"test.jpg"},
			recursive: true,
			expectedTags: []string{
				"-FileModifyDate<filename",
				"-CreateDate<filename",
				"-TrackModifyDate<filename",
				"-r",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recursively = tt.recursive
			dryRun = false

			// Create a test exiftool wrapper
			testTool := tools.NewTestExifTool()
			defer testTool.Clear()

			// Run the command
			runFixDates(datesCmd, tt.args)

			testArgs := testTool.Args

			assert.Contains(t, testArgs.Args, tt.args[0], "source path not set in %s", tt.name)
			for _, tag := range tt.expectedTags {
				assert.Contains(t, testArgs.Args, tag, "missing expected tag in %s", tt.name)
			}

			for _, tag := range tt.unexpectedTags {
				assert.NotContains(t, testArgs.Args, tag, "found unexpected tag in %s", tt.name)
			}

			assert.True(t, testTool.ExecCalled, "exiftool exec should be called in %s", tt.name)
		})
	}
}

func TestRunFixDates_DryRun(t *testing.T) {
	// Save original values to restore after test
	origRecursively := recursively
	defer func() {
		recursively = origRecursively
	}()

	origDryRun := dryRun
	defer func() {
		dryRun = origDryRun
	}()

	tests := []struct {
		name           string
		args           []string
		recursive      bool
		expectedTags   []string
		unexpectedTags []string
	}{
		{
			name:      "without recursion",
			args:      []string{"test.jpg"},
			recursive: false,
			expectedTags: []string{
				"-WriteNothing<filename",
			},
			unexpectedTags: []string{
				"-FileModifyDate<filename",
				"-CreateDate<filename",
				"-TrackModifyDate<filename",
				"-r",
			},
		},
		{
			name:      "with recursion",
			args:      []string{"test.jpg"},
			recursive: true,
			expectedTags: []string{
				"-WriteNothing<filename",
				"-r",
			},
			unexpectedTags: []string{
				"-FileModifyDate<filename",
				"-CreateDate<filename",
				"-TrackModifyDate<filename",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recursively = tt.recursive
			dryRun = true

			// Create a test exiftool wrapper
			testTool := tools.NewTestExifTool()
			defer testTool.Clear()

			// Run the command
			runFixDates(datesCmd, tt.args)

			testArgs := testTool.Args

			assert.Contains(t, testArgs.Args, tt.args[0], "source path not set in %s", tt.name)
			for _, tag := range tt.expectedTags {
				assert.Contains(t, testArgs.Args, tag, "missing expected tag in %s", tt.name)
			}

			for _, tag := range tt.unexpectedTags {
				assert.NotContains(t, testArgs.Args, tag, "found unexpected tag in %s", tt.name)
			}

			assert.True(t, testTool.ExecCalled, "exiftool exec should be called in %s", tt.name)
		})
	}
}
