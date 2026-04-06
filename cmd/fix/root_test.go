package fix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixCmd_Initialization(t *testing.T) {
	assert.NotNil(t, fixCmd)
	assert.NotNil(t, fixCmd.Use)
	assert.NotNil(t, fixCmd.Short)
	assert.NotNil(t, fixCmd.Long)
}

func TestFixCmd_Flags(t *testing.T) {
	persistentFlags := fixCmd.PersistentFlags()
	assert.True(t, persistentFlags.HasFlags())

	dryRunFlag := persistentFlags.Lookup("dry")
	assert.NotNil(t, dryRunFlag)
	assert.Equal(t, "d", dryRunFlag.Shorthand)
	assert.Equal(t, "false", dryRunFlag.DefValue)
	assert.Equal(t, "bool", dryRunFlag.Value.Type())

	recursivelyFlag := persistentFlags.Lookup("recursively")
	assert.NotNil(t, recursivelyFlag)
	assert.Equal(t, "r", recursivelyFlag.Shorthand)
	assert.Equal(t, "false", recursivelyFlag.DefValue)
	assert.Equal(t, "bool", recursivelyFlag.Value.Type())

	flags := fixCmd.Flags()
	assert.False(t, flags.HasFlags())
}

func TestFixCmd_CommandStructure(t *testing.T) {
	assert.Equal(t, "fix", fixCmd.Name())
	assert.Equal(t, "media-tool fix", fixCmd.CommandPath())

	flags := fixCmd.Flags()
	assert.Equal(t, 0, flags.NFlag())
}
