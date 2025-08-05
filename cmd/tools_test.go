package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractAbsPath(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		argPosition  int
		defaultValue string
		expected     string
	}{
		{
			name:         "extract path from args",
			args:         []string{"test", "file"},
			argPosition:  1,
			defaultValue: "/default/path",
			expected:     func() string { abs, _ := filepath.Abs("file"); return abs }(),
		},
		{
			name:         "use default when arg position out of bounds",
			args:         []string{"test"},
			argPosition:  1,
			defaultValue: "default",
			expected:     func() string { abs, _ := filepath.Abs("default"); return abs }(),
		},
		{
			name:         "use default when args empty",
			args:         []string{},
			argPosition:  0,
			defaultValue: "default",
			expected:     func() string { abs, _ := filepath.Abs("default"); return abs }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractAbsPath(tt.args, tt.argPosition, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractPath(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		argPosition  int
		defaultValue string
		expected     string
	}{
		{
			name:         "extract path from args",
			args:         []string{"test", "/path/to/file"},
			argPosition:  1,
			defaultValue: "/default/path",
			expected:     "/path/to/file",
		},
		{
			name:         "use default when arg position out of bounds",
			args:         []string{"test"},
			argPosition:  1,
			defaultValue: "/default/path",
			expected:     "/default/path",
		},
		{
			name:         "use default when args empty",
			args:         []string{},
			argPosition:  0,
			defaultValue: "/default/path",
			expected:     "/default/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractPath(tt.args, tt.argPosition, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetAbsPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "absolute path remains absolute",
			path:     "/absolute/path",
			expected: func() string { abs, _ := filepath.Abs("/absolute/path"); return abs }(),
		},
		{
			name:     "relative path becomes absolute",
			path:     "relative/path",
			expected: func() string { abs, _ := filepath.Abs("relative/path"); return abs }(),
		},
		{
			name:     "current directory",
			path:     ".",
			expected: func() string { abs, _ := filepath.Abs("."); return abs }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getAbsPath(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}


func prepareNonEmptyDir(t *testing.T, tempDir string) {
	nonEmptyDir := filepath.Join(tempDir, "nonempty")
	err := os.Mkdir(nonEmptyDir, 0755)
	require.NoError(t, err)
	testFile := filepath.Join(nonEmptyDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)
}

func TestRemoveDir_NonEmptyDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_remove_dir")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	prepareNonEmptyDir(t, tempDir)
	dirName := filepath.Join(tempDir, "nonempty")

	removeDir(dirName, true)

	_, err = os.Stat(dirName)
	assert.Error(t, err, "Directory should not exist")
}

func TestRemoveDir_EmptyDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_remove_dir")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	prepareNonEmptyDir(t, tempDir)
	dirName := filepath.Join(tempDir, "nonempty")

	removeDir(dirName, false)

	_, err = os.Stat(dirName)
	assert.NoError(t, err, "Directory should exist")
}

func TestRemoveDir_UnExistedDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_remove_dir")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dirName := filepath.Join(tempDir, "new")

	removeDir(dirName, false)

	_, err = os.Stat(dirName)
	assert.Error(t, err, "Directory should not exist")
}

func TestRemoveFiles_ExistedFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_remove_files")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testFiles := []string{
		filepath.Join(tempDir, "test.txt"),
		filepath.Join(tempDir, "other.dat"),
	}

	for _, file := range testFiles {
		err = os.WriteFile(file, []byte("test content"), 0644)
		require.NoError(t, err)
	}

	removeFiles(tempDir, "*.txt")

	_, err = os.Stat(testFiles[0])
	assert.Error(t, err, "File should be removed")

	_, err = os.Stat(testFiles[1])
	assert.NoError(t, err, "File should still exist")
}

func TestCheckDirEmpty(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_check_dir_empty")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create an empty subdirectory
	emptyDir := filepath.Join(tempDir, "empty")
	err = os.Mkdir(emptyDir, 0755)
	require.NoError(t, err)

	// Create a directory with a file
	nonEmptyDir := filepath.Join(tempDir, "nonempty")
	err = os.Mkdir(nonEmptyDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(nonEmptyDir, "test.txt"), []byte("test"), 0644)
	require.NoError(t, err)

	// Create a directory with an empty subdirectory (should be considered empty)
	nestedDir := filepath.Join(tempDir, "nested")
	err = os.Mkdir(nestedDir, 0755)
	require.NoError(t, err)
	err = os.Mkdir(filepath.Join(nestedDir, "subdir"), 0755)
	require.NoError(t, err)

	tests := []struct {
		name     string
		dirName  string
		expected bool
	}{
		{
			name:     "empty directory",
			dirName:  emptyDir,
			expected: true,
		},
		{
			name:     "non-empty directory with file",
			dirName:  nonEmptyDir,
			expected: false,
		},
		{
			name:     "nested directory with empty subdirectory",
			dirName:  nestedDir,
			expected: true,
		},
		{
			name:     "non-existent directory",
			dirName:  filepath.Join(tempDir, "nonexistent"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkDirEmpty(tt.dirName)
			assert.Equal(t, tt.expected, result)
		})
	}
}
