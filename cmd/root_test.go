package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCmd_Initialization(t *testing.T) {
	assert.NotNil(t, RootCmd)
	assert.NotNil(t, RootCmd.Use)
	assert.NotNil(t, RootCmd.Short)
	assert.NotNil(t, RootCmd.Long)
}

func TestRootCmd_Flags(t *testing.T) {
	flags := RootCmd.Flags()
	assert.False(t, flags.HasFlags())
}

func TestRootCmd_CommandStructure(t *testing.T) {
	assert.Equal(t, "media-tool", RootCmd.Name())
	assert.Equal(t, "media-tool", RootCmd.CommandPath())

	flags := RootCmd.Flags()
	assert.Equal(t, 0, flags.NFlag())
}
