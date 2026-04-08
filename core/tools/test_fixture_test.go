package tools

import (
	"io"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

type TempDir struct {
	rootDir string
	t       *testing.T
}

func NewTempDir(t *testing.T) *TempDir {
	tempDir, err := os.MkdirTemp("", "test_remove_dir")
	require.NoError(t, err)

	result := TempDir{rootDir: tempDir, t: t}
	return &result
}

func (t *TempDir) Clean() {
	err := os.RemoveAll(t.rootDir)
	require.NoError(t.t, err)
}

func (t *TempDir) RootDir() string {
	return t.rootDir
}

func (t *TempDir) Dir(dirName string) string {
	result := path.Join(t.rootDir, dirName)
	err := os.MkdirAll(result, os.ModePerm)
	require.NoError(t.t, err)

	return result
}

func (t *TempDir) File(dirName string, fileName string, fileContent string) string {
	dir := path.Join(t.rootDir, dirName)
	err := os.MkdirAll(dir, os.ModePerm)
	require.NoError(t.t, err)

	result := path.Join(dir, fileName)
	err = os.WriteFile(result, []byte(fileContent), 0644)
	require.NoError(t.t, err)

	return result
}

func (t *TempDir) IsEmptyDir(dirPath string) bool {
	dir, err := os.Open(dirPath)
	if err != nil {
		return os.IsNotExist(err)
	}
	defer dir.Close()
	_, err = dir.Readdirnames(1)
	return err == io.EOF
}

func (t *TempDir) HasFile(subDir string, fileName string) bool {
	_, err := os.Stat(path.Join(t.rootDir, subDir, fileName))
	return err == nil
}
