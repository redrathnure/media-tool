package clean

import (
	"testing"

	"github.com/redrathnure/media-tool/core/tools"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestMetadata_Initialization(t *testing.T) {
	assert.NotNil(t, metadataCmd)
	assert.NotNil(t, metadataCmd.Use)
	assert.NotNil(t, metadataCmd.Short)
	assert.NotNil(t, metadataCmd.Long)
}

func TestMetadata_Flags(t *testing.T) {
	persistentFlags := metadataCmd.PersistentFlags()
	assert.False(t, persistentFlags.HasFlags())

	flags := metadataCmd.Flags()
	locationFlag := flags.Lookup("includingLocation")
	assert.NotNil(t, locationFlag)
	assert.Equal(t, "l", locationFlag.Shorthand)
	assert.Equal(t, "false", locationFlag.DefValue)
	assert.Equal(t, "bool", locationFlag.Value.Type())

	vendorFlag := metadataCmd.Flags().Lookup("includingVendor")
	assert.NotNil(t, vendorFlag)
	assert.Equal(t, "s", vendorFlag.Shorthand)
	assert.Equal(t, "true", vendorFlag.DefValue)
	assert.Equal(t, "bool", vendorFlag.Value.Type())

	cameraFlag := metadataCmd.Flags().Lookup("includingCamera")
	assert.NotNil(t, cameraFlag)
	assert.Equal(t, "p", cameraFlag.Shorthand)
	assert.Equal(t, "false", cameraFlag.DefValue)
	assert.Equal(t, "bool", cameraFlag.Value.Type())
}

func TestMetadata_CommandStructure(t *testing.T) {
	assert.Equal(t, "metadata", metadataCmd.Name())
	assert.Equal(t, "media-tool clean metadata", metadataCmd.CommandPath())
	//assert.Equal(t, cobra.RangeArgs(1, 1), cleanMetadataCmd.Args)

	flags := metadataCmd.Flags()
	assert.Equal(t, 0, flags.NFlag())
}

func TestMetadata_ArgValidation(t *testing.T) {
	// Test no args (should fail)
	err := metadataCmd.Args(metadataCmd, []string{})
	assert.Error(t, err)

	// Test single arg (should pass)
	err = metadataCmd.Args(metadataCmd, []string{"test.jpg"})
	assert.NoError(t, err)

	// Test multiple args (should fail)
	err = metadataCmd.Args(metadataCmd, []string{"test1.jpg", "test2.jpg"})
	assert.Error(t, err)
}

func TestRunCleanMetadata(t *testing.T) {
	// Save original values to restore after test
	origLocation := includingLocation
	origVendor := includingVendor
	origCamera := includingCamera
	origRecursively := recursively
	origDryRun := dryRun
	defer func() {
		includingLocation = origLocation
		includingVendor = origVendor
		includingCamera = origCamera
		recursively = origRecursively
		dryRun = origDryRun
	}()

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
			dryRun = tt.dryRun

			// Mock the exifTool exec calls
			testTool := tools.NewTestExifTool()
			defer func() { testTool.Clear() }()

			cmd := &cobra.Command{}

			// Run the command
			runMetadata(cmd, tt.args)

			if tt.dryRun {
				assert.Len(t, testTool.Calls, 0, "exiftool exec was called when DryRun is true in %s", tt.name)
			} else {
				assert.Len(t, testTool.Calls, 1, "exiftool exec should be called in %s", tt.name)

				testArgs := testTool.Calls[0]

				assert.Contains(t, testArgs, tt.args[0], "source path not set in %s", tt.name)

				for _, tag := range tt.expectedTags {
					assert.Contains(t, testArgs, tag, "missing expected tag in %s", tt.name)
				}

				for _, tag := range tt.unexpectedTags {
					assert.NotContains(t, testArgs, tag, "found unexpected tag in %s", tt.name)
				}
			}
		})
	}
}
