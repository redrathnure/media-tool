package import_

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalCmd_Initialization(t *testing.T) {
	assert.NotNil(t, localCmd)
	assert.NotNil(t, localCmd.Use)
	assert.NotNil(t, localCmd.Short)
	assert.NotNil(t, localCmd.Long)
}

func TestLocalCmd_Flags(t *testing.T) {
	persistentFlags := localCmd.PersistentFlags()
	assert.False(t, persistentFlags.HasFlags())

	flags := localCmd.Flags()
	sourceSubDirFlag := flags.Lookup("sourceSubDir")
	assert.NotNil(t, sourceSubDirFlag)
	assert.Equal(t, "s", sourceSubDirFlag.Shorthand)
	assert.Equal(t, "false", sourceSubDirFlag.DefValue)
	assert.Equal(t, "bool", sourceSubDirFlag.Value.Type())

	localRenameFlag := flags.Lookup("rename")
	assert.NotNil(t, localRenameFlag)
	assert.Equal(t, "r", localRenameFlag.Shorthand)
	assert.Equal(t, "false", localRenameFlag.DefValue)
	assert.Equal(t, "bool", localRenameFlag.Value.Type())

	dateFormatFlag := flags.Lookup("dateFormat")
	assert.NotNil(t, dateFormatFlag)
	assert.Equal(t, "f", dateFormatFlag.Shorthand)
	assert.Equal(t, "%Y.%m.%d", dateFormatFlag.DefValue)
	assert.Equal(t, "string", dateFormatFlag.Value.Type())
}

func TestLocalCmd_CommandStructure(t *testing.T) {
	assert.Equal(t, "local", localCmd.Name())
	assert.Equal(t, "media-tool import local", localCmd.CommandPath())

	flags := localCmd.Flags()
	assert.Equal(t, 0, flags.NFlag())
}

func TestLocalCmd_ArgValidation(t *testing.T) {
	// Test no args
	err := localCmd.Args(localCmd, []string{})
	assert.Error(t, err)

	// Test single arg
	err = localCmd.Args(localCmd, []string{"test.jpg"})
	assert.NoError(t, err)

	// Test multiple args
	err = localCmd.Args(localCmd, []string{"test1.jpg", "test2.jpg"})
	assert.NoError(t, err)

	// Test multiple args
	err = localCmd.Args(localCmd, []string{"test1.jpg", "test2.jpg", "test3.jpg"})
	assert.Error(t, err)
}
