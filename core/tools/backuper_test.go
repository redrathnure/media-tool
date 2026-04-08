package tools

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBackuper_WithDryRun(t *testing.T) {
	sut := NewBackuper(".", true, false, "test")

	assert.False(t, sut.ShouldExifDeleteOriginal())
	assert.False(t, sut.ShouldPerformBackup())
	assert.Equal(t, ".", sut.GetBackupDir())
}

func TestNewBackuper_WithNoneLocation(t *testing.T) {
	tests := []string{"", "none", "None", "NONE"}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			sut := NewBackuper(tt, false, false, "test")

			assert.True(t, sut.ShouldExifDeleteOriginal())
			assert.False(t, sut.ShouldPerformBackup())
			assert.Equal(t, tt, sut.GetBackupDir())
		})
	}
}

func TestNewBackuper_WithInPlaceLocation(t *testing.T) {
	tests := []string{"in_place", "In_Place", "IN_PLACE"}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			sut := NewBackuper(tt, false, false, "test")

			assert.False(t, sut.ShouldExifDeleteOriginal())
			assert.False(t, sut.ShouldPerformBackup())
			assert.Equal(t, tt, sut.GetBackupDir())
		})
	}
}

func TestNewBackuper_WithDirLocation(t *testing.T) {
	tests := []string{".", "~/tmp"}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			sut := NewBackuper(tt, false, false, "test")

			assert.False(t, sut.ShouldExifDeleteOriginal())
			assert.True(t, sut.ShouldPerformBackup())
			assert.Contains(t, sut.GetBackupDir(), "test")
		})
	}
}

func TestCleanupWorkDir_Recursive(t *testing.T) {
	tmp := prepareCleanupWorkDir(t)
	defer tmp.Clean()

	backupDir := path.Join(tmp.RootDir(), "test")
	sut := NewBackuper(backupDir, false, true, "test")

	err := sut.CleanupWorkDir(tmp.Dir("working"))
	assert.NoError(t, err)

	assert.True(t, tmp.HasFile("working", "f1.jpg"))
	assert.False(t, tmp.HasFile("working", "f1.jpg_original"))

	assert.True(t, tmp.HasFile(path.Join("working", "sub"), "f2.jpg"))
	assert.False(t, tmp.HasFile(path.Join("working", "sub"), "f2.jpg_original"))

	assert.False(t, tmp.IsEmptyDir(backupDir))

	assert.True(t, tmp.HasFile(path.Join("test", sut.backupName), "f1.jpg_original"))
	assert.True(t, tmp.HasFile(path.Join("test", sut.backupName, "sub"), "f2.jpg_original"))
}

func TestCleanupWorkDir_NoneRecursive(t *testing.T) {
	tmp := prepareCleanupWorkDir(t)
	defer tmp.Clean()

	backupDir := path.Join(tmp.RootDir(), "test")
	sut := NewBackuper(backupDir, false, false, "test")

	err := sut.CleanupWorkDir(tmp.Dir("working"))
	assert.NoError(t, err)

	assert.True(t, tmp.HasFile("working", "f1.jpg"))
	assert.False(t, tmp.HasFile("working", "f1.jpg_original"))

	assert.True(t, tmp.HasFile(path.Join("working", "sub"), "f2.jpg"))
	assert.True(t, tmp.HasFile(path.Join("working", "sub"), "f2.jpg_original"))

	assert.False(t, tmp.IsEmptyDir(backupDir))

	assert.True(t, tmp.HasFile(path.Join("test", sut.backupName), "f1.jpg_original"))
}

func TestCleanupWorkDir_FileAsWorkDir(t *testing.T) {
	tmp := prepareCleanupWorkDir(t)
	defer tmp.Clean()

	backupDir := path.Join(tmp.RootDir(), "test")
	sut := NewBackuper(backupDir, false, true, "test")

	err := sut.CleanupWorkDir(path.Join(tmp.Dir("working"), "f1.jpg"))
	assert.NoError(t, err)

	assert.True(t, tmp.HasFile("working", "f1.jpg"))
	assert.False(t, tmp.HasFile("working", "f1.jpg_original"))

	assert.True(t, tmp.HasFile(path.Join("working", "sub"), "f2.jpg"))
	assert.False(t, tmp.HasFile(path.Join("working", "sub"), "f2.jpg_original"))

	assert.False(t, tmp.IsEmptyDir(backupDir))

	assert.True(t, tmp.HasFile(path.Join("test", sut.backupName), "f1.jpg_original"))
	assert.True(t, tmp.HasFile(path.Join("test", sut.backupName, "sub"), "f2.jpg_original"))
}

func TestCleanupWorkDir_WithDryRun(t *testing.T) {
	tmp := prepareCleanupWorkDir(t)
	defer tmp.Clean()

	backupDir := path.Join(tmp.RootDir(), "test")
	sut := NewBackuper(backupDir, true, true, "test")

	err := sut.CleanupWorkDir(tmp.Dir("working"))
	assert.NoError(t, err)

	assert.True(t, tmp.HasFile("working", "f1.jpg"))
	assert.True(t, tmp.HasFile("working", "f1.jpg_original"))

	assert.True(t, tmp.HasFile(path.Join("working", "sub"), "f2.jpg"))
	assert.True(t, tmp.HasFile(path.Join("working", "sub"), "f2.jpg_original"))

	assert.True(t, tmp.IsEmptyDir(backupDir))
}

func TestCleanupWorkDir_NoneLocation(t *testing.T) {
	tmp := prepareCleanupWorkDir(t)
	defer tmp.Clean()

	backupDir := path.Join(tmp.RootDir(), "test")
	sut := NewBackuper("none", false, true, "test")

	err := sut.CleanupWorkDir(tmp.Dir("working"))
	assert.NoError(t, err)

	assert.True(t, tmp.HasFile("working", "f1.jpg"))
	assert.True(t, tmp.HasFile("working", "f1.jpg_original"))

	assert.True(t, tmp.HasFile(path.Join("working", "sub"), "f2.jpg"))
	assert.True(t, tmp.HasFile(path.Join("working", "sub"), "f2.jpg_original"))

	assert.True(t, tmp.IsEmptyDir(backupDir))
}

func TestCleanupWorkDir_InPlaceLocation(t *testing.T) {
	tmp := prepareCleanupWorkDir(t)
	defer tmp.Clean()

	backupDir := path.Join(tmp.RootDir(), "test")
	sut := NewBackuper("in_place", false, true, "test")

	err := sut.CleanupWorkDir(tmp.Dir("working"))
	assert.NoError(t, err)

	assert.True(t, tmp.HasFile("working", "f1.jpg"))
	assert.True(t, tmp.HasFile("working", "f1.jpg_original"))

	assert.True(t, tmp.HasFile(path.Join("working", "sub"), "f2.jpg"))
	assert.True(t, tmp.HasFile(path.Join("working", "sub"), "f2.jpg_original"))

	assert.True(t, tmp.IsEmptyDir(backupDir))
}

func prepareCleanupWorkDir(t *testing.T) *TempDir {
	tmp := NewTempDir(t)

	tmp.File("working", "f1.jpg", "")
	tmp.File("working", "f1.jpg_original", "")

	tmp.File(path.Join("working", "sub"), "f2.jpg", "")
	tmp.File(path.Join("working", "sub"), "f2.jpg_original", "")

	return tmp
}
