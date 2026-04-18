package import_

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSdPhotosCmd_Initialization(t *testing.T) {
	assert.NotNil(t, sdPhotosCmd)
	assert.NotNil(t, sdPhotosCmd.Use)
	assert.NotNil(t, sdPhotosCmd.Short)
	assert.NotNil(t, sdPhotosCmd.Long)
}

func TestSdPhotosCmd_Flags(t *testing.T) {
	persistentFlags := sdPhotosCmd.PersistentFlags()
	assert.False(t, persistentFlags.HasFlags())

	flags := sdPhotosCmd.Flags()
	assert.True(t, flags.HasFlags())

	localRenameFlag := flags.Lookup("rename")
	assert.NotNil(t, localRenameFlag)
	assert.Equal(t, "r", localRenameFlag.Shorthand)
	assert.Equal(t, "false", localRenameFlag.DefValue)
	assert.Equal(t, "bool", localRenameFlag.Value.Type())
}

func TestSdPhotosCmd_CommandStructure(t *testing.T) {
	assert.Equal(t, "sdphotos", sdPhotosCmd.Name())
	assert.Equal(t, "media-tool import sdphotos", sdPhotosCmd.CommandPath())

	flags := sdPhotosCmd.Flags()
	assert.Equal(t, 0, flags.NFlag())
}

func TestSdPhotosCmd_ArgValidation(t *testing.T) {
	// Test no args
	err := sdPhotosCmd.Args(sdPhotosCmd, []string{})
	assert.NoError(t, err)

	// Test single arg
	err = sdPhotosCmd.Args(sdPhotosCmd, []string{"test.jpg"})
	assert.NoError(t, err)

	// Test multiple args
	err = sdPhotosCmd.Args(sdPhotosCmd, []string{"test1.jpg", "test2.jpg"})
	assert.Error(t, err)
}
