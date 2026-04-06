package import_

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImportCmd_Initialization(t *testing.T) {
	assert.NotNil(t, importCmd)
	assert.NotNil(t, importCmd.Use)
	assert.NotNil(t, importCmd.Short)
	assert.NotNil(t, importCmd.Long)
}

func TestImportCmd_Flags(t *testing.T) {
	persistentFlags := importCmd.PersistentFlags()
	assert.True(t, persistentFlags.HasFlags())

	dryRunFlag := persistentFlags.Lookup("dry")
	assert.NotNil(t, dryRunFlag)
	assert.Equal(t, "d", dryRunFlag.Shorthand)
	assert.Equal(t, "false", dryRunFlag.DefValue)
	assert.Equal(t, "bool", dryRunFlag.Value.Type())

			flags := importCmd.Flags()
	assert.False(t, flags.HasFlags())
}

func TestImportCmd_CommandStructure(t *testing.T) {
	assert.Equal(t, "import", importCmd.Name())
	assert.Equal(t, "media-tool import", importCmd.CommandPath())

	flags := importCmd.Flags()
	assert.Equal(t, 0, flags.NFlag())
}
