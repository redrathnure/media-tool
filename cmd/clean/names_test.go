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

	tests := []tools.ExifSingleCallCaseData{
		{
			Name:       "rename",
			Args:       []string{"testdir"},
			Recursive:  false,
			Expected:   []string{"-filename<${filename;s/ - Copy/%-c/gi;s/ Copy/%-c/gi}"},
			Unexpected: []string{"-r"},
		},
		{
			Name:      "recursive rename",
			Args:      []string{"testdir"},
			Recursive: true,
			Expected:  []string{"-filename<${filename;s/ - Copy/%-c/gi;s/ Copy/%-c/gi}", "-r"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Setup
			recursively = tt.Recursive
			dryRun = false

			// Mock the exifTool exec calls
			testTool := tools.NewTestExifTool()
			defer func() { testTool.Clear() }()

			cmd := &cobra.Command{}

			// Run the command
			runNames(cmd, tt.Args)

			// Check ExifTool calls
			testTool.AssertCalls(t, 1)

			testTool.AssertCallArgs(t, 0, tt.Args[0], tt.CallArgs())
		})
	}
}

func TestRunCleanNames_DryRun(t *testing.T) {
	origRecursively := recursively
	origDryRun := dryRun
	defer func() {
		recursively = origRecursively
		dryRun = origDryRun
	}()

	tests := []tools.ExifSingleCallCaseData{
		{
			Name:       "dry run rename",
			Args:       []string{"testdir"},
			Recursive:  false,
			Expected:   []string{"-testname<${filename;s/ - Copy/%-c/gi;s/ Copy/%-c/gi}"},
			Unexpected: []string{"-r"},
		},
		{
			Name:      "recursive rename",
			Args:      []string{"testdir"},
			Recursive: true,
			Expected:  []string{"-testname<${filename;s/ - Copy/%-c/gi;s/ Copy/%-c/gi}", "-r"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Setup
			recursively = tt.Recursive
			dryRun = true

			// Mock the exifTool exec calls
			testTool := tools.NewTestExifTool()
			defer func() { testTool.Clear() }()

			cmd := &cobra.Command{}

			// Run the command
			runNames(cmd, tt.Args)

			// Check ExifTool calls
			testTool.AssertCalls(t, 1)

			testTool.AssertCallArgs(t, 0, tt.Args[0], tt.CallArgs())
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
