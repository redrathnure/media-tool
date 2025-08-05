package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCleanNamesCmd_Initialization(t *testing.T) {
	// Test that clean names command is properly initialized
	assert.NotNil(t, cleanNamesCmd)
	assert.NotNil(t, cleanNamesCmd.Use)
	assert.NotNil(t, cleanNamesCmd.Short)
	assert.NotNil(t, cleanNamesCmd.Long)
}

func TestCleanNamesCmd_CommandStructure(t *testing.T) {
	assert.Equal(t, "names", cleanNamesCmd.Name())
	assert.Equal(t, "media-tool clean names", cleanNamesCmd.CommandPath())

	flags := cleanNamesCmd.Flags()
	assert.Equal(t, 0, flags.NFlag())
}

func TestCleanNamesCmd_ArgValidation(t *testing.T) {
	// Test no args (should fail)
	err := cleanNamesCmd.Args(cleanNamesCmd, []string{})
	assert.Error(t, err)

	// Test single arg (should pass)
	err = cleanNamesCmd.Args(cleanNamesCmd, []string{"test.jpg"})
	assert.NoError(t, err)

	// Test multiple args (should fail)
	err = cleanNamesCmd.Args(cleanNamesCmd, []string{"test1.jpg", "test2.jpg"})
	assert.Error(t, err)
}

func TestRunCleanNames(t *testing.T) {
	origRecursively := recursively
	origDryRun := DryRun
	defer func() {
		recursively = origRecursively
		DryRun = origDryRun
	}()

	// Create test exiftool wrapper to capture arguments and exec calls
	testTool := newTestExifTool()
	exifToolObj = &testTool.exifToolWrapper
	defer func() { exifToolObj = nil }()

	tests := []struct {
		name           string
		args           []string
		recursive      bool
		dryRun         bool
		expectedTags   []string
		unexpectedTags []string
	}{
		{
			name:           "rename",
			args:           []string{"testdir"},
			recursive:      false,
			dryRun:         false,
			expectedTags:   []string{"-filename<${filename;s/ - Copy/%-c/i}"},
			unexpectedTags: []string{"-r"},
		},
		{
			name:           "dry run rename",
			args:           []string{"testdir"},
			recursive:      false,
			dryRun:         true,
			expectedTags:   []string{"-testname<${filename;s/ - Copy/%-c/i}"},
			unexpectedTags: []string{"-r"},
		},
		{
			name:         "recursive rename",
			args:         []string{"testdir"},
			recursive:    true,
			dryRun:       false,
			expectedTags: []string{"-filename<${filename;s/ - Copy/%-c/i}", "-r"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			recursively = tt.recursive
			DryRun = tt.dryRun

			cmd := &cobra.Command{}

			// Run the command
			runCleanNames(cmd, tt.args)

			testArgs := testTool.args

			assert.Contains(t, testArgs.args, tt.args[0], "source path not set in %s", tt.name)
			for _, tag := range tt.expectedTags {
				assert.Contains(t, testArgs.args, tag, "missing expected tag in %s", tt.name)
			}

			for _, tag := range tt.unexpectedTags {
				assert.NotContains(t, testArgs.args, tag, "found unexpected tag in %s", tt.name)
			}

			// TODO fix mocks
			//assert.True(t, testTool.execCalled, "exiftool exec should be called in %s", tt.name)
		})
	}
}
