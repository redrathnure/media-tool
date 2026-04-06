package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigCmd_Initialization(t *testing.T) {
	assert.NotNil(t, configCmd)
	assert.NotNil(t, configCmd.Use)
	assert.NotNil(t, configCmd.Short)
	assert.NotNil(t, configCmd.Long)
}

func TestConfigCmd_Flags(t *testing.T) {
	persistentFlags := configCmd.PersistentFlags()
	assert.False(t, persistentFlags.HasFlags())

	flags := configCmd.Flags()
	assert.False(t, flags.HasFlags())
}

func TestConfigCmd_CommandStructure(t *testing.T) {
	assert.Equal(t, "config", configCmd.Name())
	assert.Equal(t, "media-tool config", configCmd.CommandPath())

	flags := configCmd.Flags()
	assert.Equal(t, 0, flags.NFlag())
}
