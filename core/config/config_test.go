package config

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteConfigAs_NewFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dstFile := path.Join(tempDir, "a", "b", "result.yml")

	sut := Empty()
	sut.Set("a", "b")

	result := sut.WriteConfigAs(dstFile)

	require.NoError(t, result)
	assert.Equal(t, "a: b\n", readFile(t, dstFile))
}

func TestWriteConfigAs_ExistedFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dstFile := path.Join(tempDir, "a", "b", "result.yml")

	sut := Empty()
	sut.Set("a", "b")

	result := sut.WriteConfigAs(dstFile)
	require.NoError(t, result)

	result = sut.WriteConfigAs(dstFile)
	require.NoError(t, result)

	assert.Equal(t, "a: b\n", readFile(t, dstFile))
}

func TestWriteConfigAs_WithDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	sut := Empty()
	sut.Set("a", "b")

	result := sut.WriteConfigAs(tempDir)
	require.Error(t, result)

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
