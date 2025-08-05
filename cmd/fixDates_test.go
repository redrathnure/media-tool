package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixDatesCmd_Initialization(t *testing.T) {
	assert.NotNil(t, fixDatesCmd)
	assert.NotNil(t, fixDatesCmd.Use)
	assert.NotNil(t, fixDatesCmd.Short)
	assert.NotNil(t, fixDatesCmd.Long)
}

func TestFixDatesCmd_Flags(t *testing.T) {
	recursiveFlag := fixDatesCmd.Flags().Lookup("recursively")
	assert.NotNil(t, recursiveFlag)
	assert.Equal(t, "r", recursiveFlag.Shorthand)
	assert.Equal(t, "false", recursiveFlag.DefValue)
	assert.Equal(t, "bool", recursiveFlag.Value.Type())
}

func TestFixDatesCmd_CommandStructure(t *testing.T) {
	assert.Equal(t, "fixDates", fixDatesCmd.Name())
	assert.Equal(t, "media-tool fixDates", fixDatesCmd.CommandPath())

	flags := fixDatesCmd.Flags()
	assert.Equal(t, 0, flags.NFlag())
}

func TestFixDatesCmd_ArgValidation(t *testing.T) {
	// Test no args (should fail)
	err := fixDatesCmd.Args(fixDatesCmd, []string{})
	assert.Error(t, err)

	// Test single arg (should pass)
	err = fixDatesCmd.Args(fixDatesCmd, []string{"test.jpg"})
	assert.NoError(t, err)

	// Test multiple args (should fail)
	err = fixDatesCmd.Args(fixDatesCmd, []string{"test1.jpg", "test2.jpg"})
	assert.Error(t, err)
}

func TestRunFixDates(t *testing.T) {
	// Save original values to restore after test
	origRecursively := recursively
	defer func() {
		recursively = origRecursively
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

			// Create a test exiftool wrapper
			testTool := newTestExifTool()
			defer testTool.clear()

			// Run the command
			runFixDates(fixDatesCmd, tt.args)

			testArgs := testTool.args

			assert.Contains(t, testArgs.args, tt.args[0], "source path not set in %s", tt.name)
			for _, tag := range tt.expectedTags {
				assert.Contains(t, testArgs.args, tag, "missing expected tag in %s", tt.name)
			}

			for _, tag := range tt.unexpectedTags {
				assert.NotContains(t, testArgs.args, tag, "found unexpected tag in %s", tt.name)
			}

			assert.True(t, testTool.execCalled, "exiftool exec should be called in %s", tt.name)
		})
	}
}
