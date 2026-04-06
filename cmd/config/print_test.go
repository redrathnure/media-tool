package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintCmd_Initialization(t *testing.T) {
	assert.NotNil(t, printCmd)
	assert.NotNil(t, printCmd.Use)
	assert.NotNil(t, printCmd.Short)
	assert.NotNil(t, printCmd.Long)
}

func TestPrintCmd_Flags(t *testing.T) {
	persistentFlags := printCmd.PersistentFlags()
	assert.False(t, persistentFlags.HasFlags())

	flags := printCmd.Flags()
	assert.False(t, flags.HasFlags())
}

func TestPrintCmd_CommandStructure(t *testing.T) {
	assert.Equal(t, "print", printCmd.Name())
	assert.Equal(t, "media-tool config print", printCmd.CommandPath())

	flags := printCmd.Flags()
	assert.Equal(t, 0, flags.NFlag())
}

func TestPrintCmd_ArgValidation(t *testing.T) {
	// Test no args
	err := printCmd.Args(printCmd, []string{})
	assert.NoError(t, err)

	// Test single arg
	err = printCmd.Args(printCmd, []string{"test.jpg"})
	assert.Error(t, err)

	// Test multiple args
	err = printCmd.Args(printCmd, []string{"test1.jpg", "test2.jpg"})
	assert.Error(t, err)
}
