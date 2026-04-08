package config

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveConfig_NewFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dstFile := path.Join(tempDir, "a", "b", "result.yml")

	sut := Empty()
	sut.Set("a", "b")

	result := sut.SaveConfig(dstFile)

	require.NoError(t, result)
	assert.Equal(t, "a: b\n", readFile(t, dstFile))
}

func TestSaveConfig_ExistedFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dstFile := path.Join(tempDir, "a", "b", "result.yml")

	sut := Empty()
	sut.Set("a", "b")

	result := sut.SaveConfig(dstFile)
	require.NoError(t, result)

	result = sut.SaveConfig(dstFile)
	require.NoError(t, result)

	assert.Equal(t, "a: b\n", readFile(t, dstFile))
}

func TestSaveConfig_WithDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	sut := Empty()
	sut.Set("a", "b")

	result := sut.SaveConfig(tempDir)
	require.Error(t, result)
}

func TestGetAsYaml_WithDir(t *testing.T) {
	sut := Empty()
	sut.Set("a", "b")

	result, err := sut.GetAsYaml()
	require.NoError(t, err)
	assert.Contains(t, string(result[:]), "a: b")
}

func readFile(t *testing.T, fileName string) string {
	data, err := os.ReadFile(fileName)
	require.NoError(t, err)

	return string(data)
}

func TestReset(t *testing.T) {
	sut := Empty()
	sut.Set("a", "b")

	assert.Equal(t, "b", sut.ko.Get("a"))

	sut.Reset()

	assert.Equal(t, nil, sut.ko.Get("a"))
}

func TestSetDefault_Empty(t *testing.T) {
	sut := Empty()
	sut.Set("a", "")

	sut.SetDefault("a", "default")

	assert.Equal(t, "", sut.ko.Get("a"))
}

func TestSetDefault_Existed(t *testing.T) {
	sut := Empty()
	sut.Set("a", "b")

	sut.SetDefault("a", "default")

	assert.Equal(t, "b", sut.ko.Get("a"))
}

func TestSetDefault_NoValue(t *testing.T) {
	sut := Empty()

	sut.SetDefault("a", "default")

	assert.Equal(t, "default", sut.ko.Get("a"))
}

func TestDefaultValues(t *testing.T) {
	sut := Empty()

	for key, val := range defaults {
		assert.Equal(t, 5, len(val), "No defaults for "+key)
	}

	assert.NotNil(t, sut.GetExifToolPath())
	assert.NotNil(t, sut.GetBackupLocation())
	assert.NotNil(t, sut.GetImportCamVideoDefaultDst())
	assert.NotNil(t, sut.GetImportGoProDefaultDst())
	assert.NotNil(t, sut.GetImportSdPhotosDefaultDst())
}
