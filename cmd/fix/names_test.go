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
				"-FileName<CreateDate",
				"-if",
				"$filename !~ /WhatsApp/i",
				"-if",
				"$filename !~ /_x265/i",
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
				"-FileName<CreateDate",
				"-if",
				"$filename !~ /WhatsApp/i",
				"-if",
				"$filename !~ /_x265/i",
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
			runFixNames(namesCmd, tt.args)

			assert.Len(t, testTool.Calls, 2, "exiftool exec should be called in %s", tt.name)
			testArgs := testTool.Calls[1]

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

func TestRunFixNames_DryRun(t *testing.T) {
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
				"-TestName<CreateDate",
				"-if",
				"$filename !~ /WhatsApp/i",
				"-if",
				"$filename !~ /_x265/i",
			},
			unexpectedTags: []string{
				"-FileModifyDate<CreateDate",
				"-r",
			},
		},
		{
			name:      "with recursion",
			args:      []string{"test.jpg"},
			recursive: true,
			expectedTags: []string{
				"-TestName<CreateDate",
				"-if",
				"$filename !~ /WhatsApp/i",
				"-if",
				"$filename !~ /_x265/i",
				"-r",
			},
			unexpectedTags: []string{
				"-FileModifyDate<CreateDate",
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
			runFixNames(namesCmd, tt.args)

			assert.Len(t, testTool.Calls, 2, "exiftool exec should be called in %s", tt.name)
			testArgs := testTool.Calls[1]

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
