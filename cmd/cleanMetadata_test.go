package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCleanMetadataCmd_Initialization(t *testing.T) {
	assert.NotNil(t, cleanMetadataCmd)
	assert.NotNil(t, cleanMetadataCmd.Use)
	assert.NotNil(t, cleanMetadataCmd.Short)
	assert.NotNil(t, cleanMetadataCmd.Long)
}

func TestCleanMetadataCmd_Flags(t *testing.T) {
	locationFlag := cleanMetadataCmd.Flags().Lookup("includingLocation")
	assert.NotNil(t, locationFlag)
	assert.Equal(t, "l", locationFlag.Shorthand)
	assert.Equal(t, "false", locationFlag.DefValue)
	assert.Equal(t, "bool", locationFlag.Value.Type())

	vendorFlag := cleanMetadataCmd.Flags().Lookup("includingVendor")
	assert.NotNil(t, vendorFlag)
	assert.Equal(t, "s", vendorFlag.Shorthand)
	assert.Equal(t, "true", vendorFlag.DefValue)
	assert.Equal(t, "bool", vendorFlag.Value.Type())

	cameraFlag := cleanMetadataCmd.Flags().Lookup("includingCamera")
	assert.NotNil(t, cameraFlag)
	assert.Equal(t, "p", cameraFlag.Shorthand)
	assert.Equal(t, "false", cameraFlag.DefValue)
	assert.Equal(t, "bool", cameraFlag.Value.Type())
}

func TestCleanMetadataCmd_CommandStructure(t *testing.T) {
	assert.Equal(t, "metadata", cleanMetadataCmd.Name())
	assert.Equal(t, "media-tool clean metadata", cleanMetadataCmd.CommandPath())
	//assert.Equal(t, cobra.RangeArgs(1, 1), cleanMetadataCmd.Args)

	flags := cleanMetadataCmd.Flags()
	assert.Equal(t, 0, flags.NFlag())
}

func TestCleanMetadataCmd_ArgValidation(t *testing.T) {
	// Test no args (should fail)
	err := cleanMetadataCmd.Args(cleanMetadataCmd, []string{})
	assert.Error(t, err)

	// Test single arg (should pass)
	err = cleanMetadataCmd.Args(cleanMetadataCmd, []string{"test.jpg"})
	assert.NoError(t, err)

	// Test multiple args (should fail)
	err = cleanMetadataCmd.Args(cleanMetadataCmd, []string{"test1.jpg", "test2.jpg"})
	assert.Error(t, err)
}

func TestRunCleanMetadata(t *testing.T) {
	// Save original values to restore after test
	origLocation := includingLocation
	origVendor := includingVendor
	origCamera := includingCamera
	origRecursively := recursively
	origDryRun := DryRun
	defer func() {
		includingLocation = origLocation
		includingVendor = origVendor
		includingCamera = origCamera
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
		location       bool
		vendor         bool
		camera         bool
		recursive      bool
		dryRun         bool
		expectedTags   []string
		unexpectedTags []string
	}{
		{
			name:      "all flags enabled",
			args:      []string{"test.jpg"},
			location:  true,
			vendor:    true,
			camera:    true,
			recursive: true,
			expectedTags: []string{
				"-gps:all=",
				"-Software=",
				"-Canon:all=",
				"-r",
			},
		}, {
			name:      "only vendor",
			args:      []string{"test.jpg"},
			location:  false,
			vendor:    true,
			camera:    false,
			recursive: false,
			expectedTags: []string{
				"-Software=",
			},
			unexpectedTags: []string{
				"-gps:all=",
				"-Canon:all=",
				"-r",
			},
		},
		{
			name:      "only location",
			args:      []string{"test.jpg"},
			location:  true,
			vendor:    false,
			camera:    false,
			recursive: false,
			expectedTags: []string{
				"-gps:all=",
			},
			unexpectedTags: []string{
				"-Software=",
				"-Canon:all=",
				"-r",
			},
		},
		{
			name:      "only camera",
			args:      []string{"test.jpg"},
			location:  false,
			vendor:    false,
			camera:    true,
			recursive: false,
			expectedTags: []string{
				"-Canon:all=",
			},
			unexpectedTags: []string{
				"-Software=",
				"-gps:all=",
				"-r",
			},
		},
		{
			name:      "only recursive",
			args:      []string{"test.jpg"},
			location:  false,
			vendor:    false,
			camera:    false,
			recursive: true,
			expectedTags: []string{
				"-r",
			},
			unexpectedTags: []string{
				"-Software=",
				"-gps:all=",
				"-Canon:all=",
			},
		},
		{
			name:   "dry run",
			args:   []string{"test.jpg"},
			dryRun: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			includingLocation = tt.location
			includingVendor = tt.vendor
			includingCamera = tt.camera
			recursively = tt.recursive
			DryRun = tt.dryRun

			cmd := &cobra.Command{}

			// Run the command
			runCleanMetadata(cmd, tt.args)

			testArgs := testTool.args

			assert.Contains(t, testArgs.args, tt.args[0], "source path not set in %s", tt.name)

			for _, tag := range tt.expectedTags {
				assert.Contains(t, testArgs.args, tag, "missing expected tag in %s", tt.name)
			}

			for _, tag := range tt.unexpectedTags {
				assert.NotContains(t, testArgs.args, tag, "found unexpected tag in %s", tt.name)
			}

			if tt.dryRun {
				assert.False(t, testTool.execCalled, "exiftool exec was called when DryRun is true in %s", tt.name)
			}
		})
	}
}
