package testutil

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	ExampleFullExifCr3File  = "canon_full_exif.jpg"
	ExampleFullExifJpgFile  = "canon_full_exif.jpg"
	ExampleFullExifJpegFile = "canon_full_exif.jpg"
	ExampleMp4File          = "VID_20250612_184641.mp4"
	ExampleMkvFile          = "VID_20250612_184641.mkv"
	CreationDate            = "2025-06-12 18:46:41"
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

func (t *TempDir) MkDir(relDir string) string {
	result := t.DirName(relDir)
	err := os.MkdirAll(result, os.ModePerm)
	require.NoError(t.t, err)

	return result
}

func (t *TempDir) DirName(relDirName string) string {
	return path.Join(t.rootDir, relDirName)
}

func (t *TempDir) MkFile(relDir string, fileName string, fileContent string) string {
	dir := t.MkDir(relDir)

	result := path.Join(dir, fileName)
	err := os.WriteFile(result, []byte(fileContent), 0644)
	require.NoError(t.t, err)

	return result
}

func (t *TempDir) copyTestSample(sample_name string, relDir string, fileName string) string {
	dir := t.MkDir(relDir)
	dstFile := path.Join(dir, fileName)

	srcFile, _ := filepath.Abs(path.Join("..", "..", "test_samples", sample_name))

	srcF, err := os.Open(srcFile)
	require.NoError(t.t, err)
	defer srcF.Close()

	destF, err := os.Create(dstFile) // creates if file doesn't exist
	require.NoError(t.t, err)
	defer destF.Close()

	_, err = io.Copy(destF, srcF) // check first var for number of bytes copied
	require.NoError(t.t, err)

	err = destF.Sync()
	require.NoError(t.t, err)
	return srcFile
}

func (t *TempDir) MkCr3FullExif(relDir string, fileName string) string {
	return t.copyTestSample(ExampleFullExifCr3File, relDir, fileName)
}

func (t *TempDir) MkJpgFullExif(relDir string, fileName string) string {
	return t.copyTestSample(ExampleFullExifJpgFile, relDir, fileName)
}

func (t *TempDir) MkJpegFullExif(relDir string, fileName string) string {
	return t.copyTestSample(ExampleFullExifJpegFile, relDir, fileName)
}

func (t *TempDir) MkVideoMp4(relDir string, fileName string) string {
	return t.copyTestSample(ExampleMp4File, relDir, fileName)
}

func (t *TempDir) MkVideoMkv(relDir string, fileName string) string {
	return t.copyTestSample(ExampleMkvFile, relDir, fileName)
}

func (t *TempDir) IsDirEmptyA(dirPath string) bool {
	dir, err := os.Open(dirPath)
	if err != nil {
		return os.IsNotExist(err)
	}
	defer dir.Close()
	_, err = dir.Readdirnames(1)
	return err == io.EOF
}

func (t *TempDir) IsDirEmpty(relDir string) bool {
	fullDirName := t.DirName(relDir)
	return t.IsDirEmptyA(fullDirName)
}

func (t *TempDir) IsDirExistA(dirPath string) bool {
	stat, err := os.Stat(dirPath)
	return err == nil && stat.IsDir()
}

func (t *TempDir) IsDirExist(relDir string) bool {
	fullDirName := t.DirName(relDir)
	return t.IsDirExistA(fullDirName)
}

func (t *TempDir) IsFileExistA(dir string, fileName string) bool {
	_, err := os.Stat(path.Join(dir, fileName))
	return err == nil
}

func (t *TempDir) IsFileExist(relDir string, fileName string) bool {
	return t.IsFileExistA(t.DirName(relDir), fileName)
}

func (t *TempDir) FileStat(relDir string, fileName string) (size int64, date string) {
	dir := t.DirName(relDir)
	require.True(t.t, t.IsFileExistA(dir, fileName))

	stat, err := os.Stat(path.Join(dir, fileName))
	require.NoError(t.t, err)

	size = stat.Size()
	date = stat.ModTime().Format("2006-01-02 15:04:05")
	return
}

func (t *TempDir) FileSize(relDir string, fileName string) int64 {
	size, _ := t.FileStat(relDir, fileName)
	return size
}

func (t *TempDir) FileDate(relDir string, fileName string) string {
	_, date := t.FileStat(relDir, fileName)
	return date
}
