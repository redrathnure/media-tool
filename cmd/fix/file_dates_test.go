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
				"-FileModifyDate<CreateDate",
				"-FileCreateDate<CreateDate",
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
				"-FileModifyDate<CreateDate",
				"-FileCreateDate<CreateDate",
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
			runFixFileDates(fileDatesCmd, tt.args)

			assert.Len(t, testTool.Calls, 1, "exiftool exec should be called in %s", tt.name)
			testArgs := testTool.Calls[0]

			assert.Contains(t, testArgs, tt.args[0], "source path not set in %s", tt.name)
			for _, tag := range tt.expectedTags {
				assert.Contains(t, testArgs, tag, "missing expected tag in %s", tt.name)
			}

			for _, tag := range tt.unexpectedTags {
				assert.NotContains(t, testArgs, tag, "found unexpected tag in %s", tt.name)
			}

			assert.Len(t, testTool.Calls, 1, "exiftool exec should be called in %s", tt.name)
		})
	}
}

func TestRunFixFileDates_DryRun(t *testing.T) {
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
				"-WriteNothing<CreateDate",
			},
			unexpectedTags: []string{
				"-FileModifyDate<CreateDate",
				"-FileCreateDate<CreateDate",
				"-r",
			},
		},
		{
			name:      "with recursion",
			args:      []string{"test.jpg"},
			recursive: true,
			expectedTags: []string{
				"-WriteNothing<CreateDate",
				"-r",
			},
			unexpectedTags: []string{
				"-FileModifyDate<CreateDate",
				"-FileCreateDate<CreateDate",
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
			runFixFileDates(fileDatesCmd, tt.args)

			assert.Len(t, testTool.Calls, 1, "exiftool exec should be called in %s", tt.name)
			testArgs := testTool.Calls[0]

			assert.Contains(t, testArgs, tt.args[0], "source path not set in %s", tt.name)
			for _, tag := range tt.expectedTags {
				assert.Contains(t, testArgs, tag, "missing expected tag in %s", tt.name)
			}

			for _, tag := range tt.unexpectedTags {
				assert.NotContains(t, testArgs, tag, "found unexpected tag in %s", tt.name)
			}
		})
	}
}
