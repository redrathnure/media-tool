package import_

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoProCmd_Initialization(t *testing.T) {
	assert.NotNil(t, goProCmd)
	assert.NotNil(t, goProCmd.Use)
	assert.NotNil(t, goProCmd.Short)
	assert.NotNil(t, goProCmd.Long)
}

func TestGoProCmd_Flags(t *testing.T) {
	persistentFlags := goProCmd.PersistentFlags()
	assert.False(t, persistentFlags.HasFlags())

	flags := goProCmd.Flags()
	assert.False(t, flags.HasFlags())
}

func TestGoProCmd_CommandStructure(t *testing.T) {
	assert.Equal(t, "gopro", goProCmd.Name())
	assert.Equal(t, "media-tool import gopro", goProCmd.CommandPath())

	flags := goProCmd.Flags()
	assert.Equal(t, 0, flags.NFlag())
}

func TestGoProCmd_ArgValidation(t *testing.T) {
	// Test no args
	err := goProCmd.Args(goProCmd, []string{})
	assert.NoError(t, err)

	// Test single arg
	err = goProCmd.Args(goProCmd, []string{"test.jpg"})
	assert.NoError(t, err)

	// Test multiple args
	err = goProCmd.Args(goProCmd, []string{"test1.jpg", "test2.jpg"})
	assert.Error(t, err)
}
