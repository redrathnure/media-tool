package removable

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/redrathnure/media-tool/core/removable/core"
)

func TestExecutionPlan_EmptyPlan(t *testing.T) {
	dev := core.NewMockDev("Mock")

	sut, err := BuildExecutionPlan(&dev, "test")

	require.NoError(t, err)
	assert.Equal(t, []string{}, sut.files)
	assert.True(t, sut.IsEmpty())
	assert.Equal(t, 0, sut.GetFilesCount())
	assert.Equal(t, int64(0), sut.GetTotalSize())
}

func TestExecutionPlan_AddFile_1xFile(t *testing.T) {
	dev := core.NewMockDev("Mock", "test/f1.jpg")

	sut, err := BuildExecutionPlan(&dev, "test")

	require.NoError(t, err)
	assert.False(t, sut.IsEmpty())
	assert.Equal(t, 1, sut.GetFilesCount())
	assert.Equal(t, int64(10), sut.GetTotalSize())
}

func TestExecutionPlan_AddFile_3xFile(t *testing.T) {
	dev := core.NewMockDev("Mock", "test/f1.jpg", "test/f2.jpg", "test/f3.jpg")

	sut, err := BuildExecutionPlan(&dev, "test")

	require.NoError(t, err)
	assert.False(t, sut.IsEmpty())
	assert.Equal(t, 3, sut.GetFilesCount())
	assert.Equal(t, int64(30), sut.GetTotalSize())
}
