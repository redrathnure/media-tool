package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	sut := NewConfig()

	assert.NotNil(t, sut)
	assert.NotNil(t, sut.ko)
	assert.Nil(t, sut.configFiles)
}

func TestLoadFile_Empty(t *testing.T) {
	sut := NewConfig()

	result := sut.loadFile("")

	require.Error(t, result)
}

func TestLoadFile_EmptyDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	sut := NewConfig()

	result := sut.loadFile(tempDir)

	require.NoError(t, result)
}

func TestLoadFile_InvalidFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	filePath := createFile(t, tempDir, "test.txt", "zxc")

	sut := NewConfig()

	result := sut.loadFile(filePath)

	require.Error(t, result)
}

func TestLoadFile_ValidFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	filePath := createFile(t, tempDir, "test.yml", "a: b")

	sut := NewConfig()

	result := sut.loadFile(filePath)

	require.NoError(t, result)
	assert.Equal(t, "b", sut.ko.Get("a"))
}

func TestSetConfigFile_ExplainPath(t *testing.T) {
	sut := NewConfig()

	sut.SetConfigFile("~")

	assert.Len(t, sut.configFiles, 1)
	assert.NotContains(t, sut.configFiles[0], "~")
}

func TestSetConfigFile_InvalidFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	filePath := createFile(t, tempDir, "test.txt", "zxc")

	sut := NewConfig()

	sut.SetConfigFile(filePath)
	result, err := sut.Build()

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestSetConfigFile_ValidFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	filePath := createFile(t, tempDir, "test.yml", "a: b")

	sut := NewConfig()

	sut.SetConfigFile(filePath)
	result, err := sut.Build()

	require.NoError(t, err)
	assert.Equal(t, "b", result.GetString("a"))
}

func TestAddConfigLocation_ExplainPath(t *testing.T) {
	sut := NewConfig()

	sut.AddConfigLocation("~")

	assert.Len(t, sut.configFiles, 1)
	assert.NotContains(t, sut.configFiles[0], "~")
}

func TestAddConfigLocation_NoFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	sut := NewConfig()

	sut.AddConfigLocation(tempDir)
	result, err := sut.Build()

	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestAddConfigLocation_InvalidFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	createFile(t, tempDir, configFileName, "zxc")

	sut := NewConfig()

	sut.AddConfigLocation(tempDir)
	result, err := sut.Build()

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestAddConfigLocation_ValidFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	createFile(t, tempDir, configFileName, "a: b")

	sut := NewConfig()

	sut.AddConfigLocation(tempDir)
	result, err := sut.Build()

	require.NoError(t, err)
	assert.Equal(t, "b", result.GetString("a"))
}

func TestAddConfigLocation_Priority(t *testing.T) {
	tempDir1, err := os.MkdirTemp("", t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(tempDir1)

	createFile(t, tempDir1, configFileName, "a: b1")

	tempDir2, err := os.MkdirTemp("", t.Name()+"2")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir2)

	createFile(t, tempDir2, configFileName, "a: b2")

	sut := NewConfig()

	sut.AddConfigLocation(tempDir1)
	sut.AddConfigLocation(tempDir2)
	result, err := sut.Build()

	require.NoError(t, err)
	assert.Equal(t, "b2", result.GetString("a"))
}

func createFile(t *testing.T, dir, fileName, content string) string {
	filePath := filepath.Join(dir, fileName)
	f, err := os.Create(filePath)
	require.NoError(t, err)

	defer f.Close()

	_, err = f.WriteString(content)
	require.NoError(t, err)

	return filePath
}
