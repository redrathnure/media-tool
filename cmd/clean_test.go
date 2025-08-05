package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanCmd_Initialization(t *testing.T) {
	// Test that clean command is properly initialized
	assert.NotNil(t, cleanCmd)
	assert.NotNil(t, cleanCmd.Use)
	assert.NotNil(t, cleanCmd.Short)
	assert.NotNil(t, cleanCmd.Long)
}

func TestCleanCmd_Flags(t *testing.T) {
	// Test that persistent flags are properly set
	dryRunFlag := cleanCmd.PersistentFlags().Lookup("dry")
	assert.NotNil(t, dryRunFlag)
	assert.Equal(t, "d", dryRunFlag.Shorthand)
	assert.Equal(t, "false", dryRunFlag.DefValue)
	assert.Equal(t, "bool", dryRunFlag.Value.Type())

	recursivelyFlag := cleanCmd.PersistentFlags().Lookup("recursively")
	assert.NotNil(t, recursivelyFlag)
	assert.Equal(t, "r", recursivelyFlag.Shorthand)
	assert.Equal(t, "false", recursivelyFlag.DefValue)
	assert.Equal(t, "bool", recursivelyFlag.Value.Type())
}

func TestCleanCmd_CommandStructure(t *testing.T) {
	assert.Equal(t, "clean", cleanCmd.Name())
	assert.Equal(t, "media-tool clean", cleanCmd.CommandPath())

	persistentFlags := cleanCmd.PersistentFlags()
	assert.Equal(t, 0, persistentFlags.NFlag())
}
