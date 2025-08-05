package cmd

import (
	"testing"

	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	assert.NotPanics(t, func() {
		initLogger()
	})

	assert.NotNil(t, log)
	assert.Equal(t, "cmd", log.Module)
}

func TestInitLoggerLevel_VerboseTrue(t *testing.T) {
	// Save original verbose value
	originalVerbose := verbose
	defer func() { verbose = originalVerbose }()

	verbose = true
	assert.NotPanics(t, func() {
		initLoggerLevel()
	})

	assert.Equal(t, logging.DEBUG, logging.GetLevel(""))
}

func TestInitLoggerLevel_VerboseFalse(t *testing.T) {
	// Save original verbose value
	originalVerbose := verbose
	defer func() { verbose = originalVerbose }()

	verbose = false
	assert.NotPanics(t, func() {
		initLoggerLevel()
	})

	assert.Equal(t, logging.INFO, logging.GetLevel(""))
}
