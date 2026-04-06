package clean

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/redrathnure/media-tool/core/tools"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNamesCmd_Initialization(t *testing.T) {
	// Test that clean names command is properly initialized
	assert.NotNil(t, namesCmd)
	assert.NotNil(t, namesCmd.Use)
	assert.NotNil(t, namesCmd.Short)
	assert.NotNil(t, namesCmd.Long)
}

func TestNamesCmd_CommandStructure(t *testing.T) {
	assert.Equal(t, "names", namesCmd.Name())
	assert.Equal(t, "media-tool clean names", namesCmd.CommandPath())

	flags := namesCmd.Flags()
	assert.Equal(t, 0, flags.NFlag())
}

func TestNamesCmd_Flags(t *testing.T) {
	persistentFlags := namesCmd.PersistentFlags()
	assert.False(t, persistentFlags.HasFlags())

	flags := namesCmd.Flags()
	assert.False(t, flags.HasFlags())
}
func TestNamesCmd_ArgValidation(t *testing.T) {
	// Test no args (should fail)
	err := namesCmd.Args(namesCmd, []string{})
	assert.Error(t, err)

	// Test single arg (should pass)
	err = namesCmd.Args(namesCmd, []string{"test.jpg"})
	assert.NoError(t, err)

	// Test multiple args (should fail)
	err = namesCmd.Args(namesCmd, []string{"test1.jpg", "test2.jpg"})
	assert.Error(t, err)
}

func TestRunCleanNames(t *testing.T) {
	origRecursively := recursively
	origDryRun := dryRun
	defer func() {
		recursively = origRecursively
		dryRun = origDryRun
	}()

	tests := []struct {
		name           string
		args           []string
		recursive      bool
		dryRun         bool
		expectedTags   []string
		unexpectedTags []string
	}{
		{
			name:           "rename",
			args:           []string{"testdir"},
			recursive:      false,
			dryRun:         false,
			expectedTags:   []string{"-filename<${filename;s/ - Copy/%-c/gi;s/ Copy/%-c/gi}"},
			unexpectedTags: []string{"-r"},
		},
		{
			name:           "dry run rename",
			args:           []string{"testdir"},
			recursive:      false,
			dryRun:         true,
			expectedTags:   []string{"-testname<${filename;s/ - Copy/%-c/gi;s/ Copy/%-c/gi}"},
			unexpectedTags: []string{"-r"},
		},
		{
			name:         "recursive rename",
			args:         []string{"testdir"},
			recursive:    true,
			dryRun:       false,
			expectedTags: []string{"-filename<${filename;s/ - Copy/%-c/gi;s/ Copy/%-c/gi}", "-r"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			recursively = tt.recursive
			dryRun = tt.dryRun

			// Mock the exifTool exec calls
			testTool := tools.NewTestExifTool()
			defer func() { testTool.Clear() }()

			cmd := &cobra.Command{}

			// Run the command
			runNames(cmd, tt.args)

			testArgs := testTool.Args

			assert.Contains(t, testArgs.Args, tt.args[0], "source path not set in %s", tt.name)
			for _, tag := range tt.expectedTags {
				assert.Contains(t, testArgs.Args, tag, "missing expected tag in %s", tt.name)
			}

			for _, tag := range tt.unexpectedTags {
				assert.NotContains(t, testArgs.Args, tag, "found unexpected tag in %s", tt.name)
			}

			assert.True(t, testTool.ExecCalled, "exiftool exec should be called in %s", tt.name)
		})
	}
}

func TestNames_CopyRenameSimple(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cleannames-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	originFileCopy := createFile(t, tmpDir, "test - Copy.jpg")

	origDryRun := dryRun
	defer func() { dryRun = origDryRun }()
	dryRun = false

	cmd := &cobra.Command{}
	runNames(cmd, []string{tmpDir})

	_, err = os.Stat(originFileCopy)
	assert.True(t, os.IsNotExist(err), "copy file should not exist after renaming")

	renamedFile := filepath.Join(tmpDir, "test.jpg")
	_, err = os.Stat(renamedFile)
	assert.NoError(t, err, "renamed file should exist")
}

func TestNames_CopyRenameWithoutDash(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cleannames-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	originFileCopy := createFile(t, tmpDir, "test Copy.jpg")

	origDryRun := dryRun
	defer func() { dryRun = origDryRun }()
	dryRun = false

	cmd := &cobra.Command{}
	runNames(cmd, []string{tmpDir})

	_, err = os.Stat(originFileCopy)
	assert.True(t, os.IsNotExist(err), "copy file should not exist after renaming")

	renamedFile := filepath.Join(tmpDir, "test.jpg")
	_, err = os.Stat(renamedFile)
	assert.NoError(t, err, "renamed file should exist")
}

func TestNames_CopyRenameWithPrefix(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cleannames-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	createFile(t, tmpDir, "test.jpg")
	originFileCopy := createFile(t, tmpDir, "test - Copy.jpg")

	origDryRun := dryRun
	defer func() { dryRun = origDryRun }()
	dryRun = false

	cmd := &cobra.Command{}
	runNames(cmd, []string{tmpDir})

	_, err = os.Stat(originFileCopy)
	assert.True(t, os.IsNotExist(err), "copy file should not exist after renaming")

	expectedFiles := []string{
		"test.jpg",
		"test-1.jpg", // From "test - Copy.jpg"
	}

	for _, f := range expectedFiles {
		_, err = os.Stat(filepath.Join(tmpDir, f))
		assert.NoError(t, err, "expected file %s should exist", f)
	}
}

func TestNames_CopyRenameWithoutDashWithPrefix(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cleannames-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	createFile(t, tmpDir, "test.jpg")
	originFileCopy := createFile(t, tmpDir, "test Copy.jpg")

	origDryRun := dryRun
	defer func() { dryRun = origDryRun }()
	dryRun = false

	cmd := &cobra.Command{}
	runNames(cmd, []string{tmpDir})

	_, err = os.Stat(originFileCopy)
	assert.True(t, os.IsNotExist(err), "copy file should not exist after renaming")

	expectedFiles := []string{
		"test.jpg",
		"test-1.jpg", // From "test Copy.jpg"
	}

	for _, f := range expectedFiles {
		_, err = os.Stat(filepath.Join(tmpDir, f))
		assert.NoError(t, err, "expected file %s should exist", f)
	}
}

func TestNames_CopyRenameDouble(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cleannames-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	origFileCopy := createFile(t, tmpDir, "test - Copy - Copy.jpg")

	origDryRun := dryRun
	defer func() { dryRun = origDryRun }()
	dryRun = false

	cmd := &cobra.Command{}
	runNames(cmd, []string{tmpDir})

	_, err = os.Stat(origFileCopy)
	assert.True(t, os.IsNotExist(err), "original file should not exist after rename")

	renamedFile := filepath.Join(tmpDir, "test.jpg")
	_, err = os.Stat(renamedFile)
	assert.NoError(t, err, "renamed file should exist")
}

func TestNames_CopyRenameDoubleExisted(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cleannames-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	createFile(t, tmpDir, "test.jpg")
	origFileCopy := createFile(t, tmpDir, "test - Copy.jpg")
	origFileCopyCopy := createFile(t, tmpDir, "test - Copy - Copy.jpg")

	origDryRun := dryRun
	defer func() { dryRun = origDryRun }()
	dryRun = false

	cmd := &cobra.Command{}
	runNames(cmd, []string{tmpDir})

	_, err = os.Stat(origFileCopy)
	assert.True(t, os.IsNotExist(err), "copy file should not exist after rename")
	// Verify the copy copy file no longer exists and new file exists
	_, err = os.Stat(origFileCopyCopy)
	assert.True(t, os.IsNotExist(err), "copy copy file should not exist after rename")

	expectedFiles := []string{
		"test.jpg",
		"test-1.jpg",
		"test-1-1.jpg", // From "test Copy.jpg"
	}

	for _, f := range expectedFiles {
		_, err = os.Stat(filepath.Join(tmpDir, f))
		assert.NoError(t, err, "expected file %s should exist", f)
	}
}

func TestNames_MultipleConflicts(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cleannames-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create initial files
	createFile(t, tmpDir, "test.jpg")
	createFile(t, tmpDir, "test-1.jpg")
	origFileCopy := createFile(t, tmpDir, "test Copy.jpg")

	origDryRun := dryRun
	defer func() { dryRun = origDryRun }()
	dryRun = false

	cmd := &cobra.Command{}
	runNames(cmd, []string{tmpDir})

	_, err = os.Stat(origFileCopy)
	assert.True(t, os.IsNotExist(err), "copy file should not exist after renaming")

	expectedFiles := []string{
		"test.jpg",
		"test-1.jpg",
		"test-2.jpg", // From "test Copy.jpg"
	}

	for _, f := range expectedFiles {
		_, err = os.Stat(filepath.Join(tmpDir, f))
		assert.NoError(t, err, "expected file %s should exist", f)
	}
}

func createFile(t *testing.T, tmpDir string, origFileName string) string {
	result := filepath.Join(tmpDir, origFileName)

	err := os.WriteFile(result, []byte("test data"), 0644)
	require.NoError(t, err)
	return result
}
