package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveCmd_Initialization(t *testing.T) {
	assert.NotNil(t, saveCmd)
	assert.NotNil(t, saveCmd.Use)
	assert.NotNil(t, saveCmd.Short)
	assert.NotNil(t, saveCmd.Long)
}

func TestSaveCmd_Flags(t *testing.T) {
	persistentFlags := saveCmd.PersistentFlags()
	assert.False(t, persistentFlags.HasFlags())

	flags := saveCmd.Flags()
	assert.False(t, flags.HasFlags())
}

func TestSaveCmd_CommandStructure(t *testing.T) {
	assert.Equal(t, "save", saveCmd.Name())
	assert.Equal(t, "media-tool config save", saveCmd.CommandPath())

	flags := saveCmd.Flags()
	assert.Equal(t, 0, flags.NFlag())
}

func TestSaveCmd_ArgValidation(t *testing.T) {
	// Test no args
	err := saveCmd.Args(saveCmd, []string{})
	assert.Error(t, err)

	// Test single arg
	err = saveCmd.Args(saveCmd, []string{"test.yml"})
	assert.NoError(t, err)

	// Test multiple args
	err = saveCmd.Args(saveCmd, []string{"test1.yml", "test2.yml"})
	assert.Error(t, err)
}
