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

func (t *TempDir) MkDir(relDirs ...string) string {
	result := t.DirName(relDirs...)
	err := os.MkdirAll(result, os.ModePerm)
	require.NoError(t.t, err)

	return result
}

func (t *TempDir) DirName(relDirs ...string) string {
	pathElements := []string{t.rootDir}
	pathElements = append(pathElements, relDirs...)
	return path.Join(pathElements...)
}

func (t *TempDir) MkFile(fileName string, fileContent string, relDirs ...string) string {
	dir := t.MkDir(relDirs...)

	result := path.Join(dir, fileName)
	err := os.WriteFile(result, []byte(fileContent), 0644)
	require.NoError(t.t, err)

	return result
}

func (t *TempDir) copyTestSample(sample_name string, fileName string, relDirs ...string) string {
	dir := t.MkDir(relDirs...)
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

func (t *TempDir) MkCr3FullExif(fileName string, relDirs ...string) string {
	return t.copyTestSample(ExampleFullExifCr3File, fileName, relDirs...)
}

func (t *TempDir) MkJpgFullExif(fileName string, relDirs ...string) string {
	return t.copyTestSample(ExampleFullExifJpgFile, fileName, relDirs...)
}

func (t *TempDir) MkJpegFullExif(fileName string, relDirs ...string) string {
	return t.copyTestSample(ExampleFullExifJpegFile, fileName, relDirs...)
}

func (t *TempDir) MkVideoMp4(fileName string, relDirs ...string) string {
	return t.copyTestSample(ExampleMp4File, fileName, relDirs...)
}

func (t *TempDir) MkVideoMkv(fileName string, relDirs ...string) string {
	return t.copyTestSample(ExampleMkvFile, fileName, relDirs...)
}

func (t *TempDir) MkTestMedia(fileNameBase string, relDirs ...string) {
	t.MkCr3FullExif(fileNameBase+".cr3", relDirs...)
	t.MkJpgFullExif(fileNameBase+".jpg", relDirs...)
	t.MkJpegFullExif(fileNameBase+".jpeg", relDirs...)
	t.MkVideoMkv(fileNameBase+".mkv", relDirs...)
	t.MkVideoMp4(fileNameBase+".mp4", relDirs...)
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

func (t *TempDir) IsDirEmpty(relDirs ...string) bool {
	fullDirName := t.DirName(relDirs...)
	return t.IsDirEmptyA(fullDirName)
}

func (t *TempDir) IsDirExistA(dirPath string) bool {
	stat, err := os.Stat(dirPath)
	return err == nil && stat.IsDir()
}

func (t *TempDir) IsDirExist(relDirs ...string) bool {
	fullDirName := t.DirName(relDirs...)
	return t.IsDirExistA(fullDirName)
}

func (t *TempDir) IsFileExistA(dir string, fileName string) bool {
	_, err := os.Stat(path.Join(dir, fileName))
	return err == nil
}

func (t *TempDir) IsFileExist(fileName string, relDirs ...string) bool {
	return t.IsFileExistA(t.DirName(relDirs...), fileName)
}

func (t *TempDir) FileStat(fileName string, relDirs ...string) (size int64, date string) {
	dir := t.DirName(relDirs...)
	require.True(t.t, t.IsFileExistA(dir, fileName))

	stat, err := os.Stat(path.Join(dir, fileName))
	require.NoError(t.t, err)

	size = stat.Size()
	date = stat.ModTime().Format("2006-01-02 15:04:05")
	return
}

func (t *TempDir) FileSize(fileName string, relDirs ...string) int64 {
	size, _ := t.FileStat(fileName, relDirs...)
	return size
}

func (t *TempDir) FileDate(fileName string, relDirs ...string) string {
	_, date := t.FileStat(fileName, relDirs...)
	return date
}
