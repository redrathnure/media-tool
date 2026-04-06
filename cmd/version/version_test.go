package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionCmd_Initialization(t *testing.T) {
	assert.NotNil(t, versionCmd)
	assert.NotNil(t, versionCmd.Use)
	assert.NotNil(t, versionCmd.Short)
	assert.NotNil(t, versionCmd.Long)
}

func TestVersionCmd_Flags(t *testing.T) {
	persistentFlags := versionCmd.PersistentFlags()
	assert.False(t, persistentFlags.HasFlags())

	flags := versionCmd.Flags()
	assert.False(t, flags.HasFlags())
}

func TestVersionCmd_CommandStructure(t *testing.T) {
	assert.Equal(t, "version", versionCmd.Name())
	assert.Equal(t, "media-tool version", versionCmd.CommandPath())

	flags := versionCmd.Flags()
	assert.Equal(t, 0, flags.NFlag())
}

func TestVersionCmd_ArgValidation(t *testing.T) {
	// Test no args
	err := versionCmd.Args(versionCmd, []string{})
	assert.NoError(t, err)

	// Test single arg
	err = versionCmd.Args(versionCmd, []string{"test.jpg"})
	assert.Error(t, err)

	// Test multiple args
	err = versionCmd.Args(versionCmd, []string{"test1.jpg", "test2.jpg"})
	assert.Error(t, err)
}
