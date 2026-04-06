package import_

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVideoCamCmd_Initialization(t *testing.T) {
	assert.NotNil(t, videoCamCmd)
	assert.NotNil(t, videoCamCmd.Use)
	assert.NotNil(t, videoCamCmd.Short)
	assert.NotNil(t, videoCamCmd.Long)
}

func TestVideoCamCmd_Flags(t *testing.T) {
	persistentFlags := videoCamCmd.PersistentFlags()
	assert.False(t, persistentFlags.HasFlags())

	flags := videoCamCmd.Flags()
	assert.False(t, flags.HasFlags())
}

func TestVideoCamCmd_CommandStructure(t *testing.T) {
	assert.Equal(t, "videocam", videoCamCmd.Name())
	assert.Equal(t, "media-tool import videocam", videoCamCmd.CommandPath())

	flags := videoCamCmd.Flags()
	assert.Equal(t, 0, flags.NFlag())
}

func TestVideoCamCmd_ArgValidation(t *testing.T) {
	// Test no args
	err := videoCamCmd.Args(videoCamCmd, []string{})
	assert.NoError(t, err)

	// Test single arg
	err = videoCamCmd.Args(videoCamCmd, []string{"test.jpg"})
	assert.NoError(t, err)

	// Test multiple args
	err = videoCamCmd.Args(videoCamCmd, []string{"test1.jpg", "test2.jpg"})
	assert.Error(t, err)
}
