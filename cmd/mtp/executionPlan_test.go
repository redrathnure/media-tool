package mtp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tobwithu/gowpd"
)

func TestExecutionPlan_EmptyPlan(t *testing.T) {
	r := require.New(t)

	sut := new(ExecutionPlan)

	r.True(sut.IsEmpty())
	r.Equal(0, sut.GetFilesCount())
	r.Equal(int64(0), sut.GetTotalSize())
}

func TestExecutionPlan_AddFile_1xFile(t *testing.T) {
	r := require.New(t)

	sut := addFile(new(ExecutionPlan), 1)

	r.False(sut.IsEmpty())
	r.Equal(1, sut.GetFilesCount())
	r.Equal(int64(100), sut.GetTotalSize())
}

func TestExecutionPlan_AddFile_3xFile(t *testing.T) {
	r := require.New(t)

	sut := addFile(new(ExecutionPlan), 3)

	r.False(sut.IsEmpty())
	r.Equal(3, sut.GetFilesCount())
	r.Equal(int64(600), sut.GetTotalSize())
}

func TestExecutionPlan_GetTotalSizeString_EmptyPlan(t *testing.T) {
	r := require.New(t)

	sut := new(ExecutionPlan)

	r.True(sut.IsEmpty())
	r.Equal("0 B", sut.GetTotalSizeString())
}

func TestExecutionPlan_GetTotalSizeString_Bytes(t *testing.T) {
	r := require.New(t)

	r.Equal("3 B", setTotalSize(new(ExecutionPlan), 3).GetTotalSizeString())
	r.Equal("3.0 KiB", setTotalSize(new(ExecutionPlan), 3*1024).GetTotalSizeString())
	r.Equal("3.0 MiB", setTotalSize(new(ExecutionPlan), 3*1024*1024+3).GetTotalSizeString())
	r.Equal("3.0 GiB", setTotalSize(new(ExecutionPlan), 3*1024*1024*1024+3).GetTotalSizeString())
	r.Equal("3.0 TiB", setTotalSize(new(ExecutionPlan), 3*1024*1024*1024*1024+3).GetTotalSizeString())
}

func TestExecutionFileIterator_EmptyPlan(t *testing.T) {
	r := require.New(t)

	executionPlan := new(ExecutionPlan)

	sut := executionPlan.GetFileInterator()

	r.False(sut.HasNext())
	r.Nil(sut.Current())
	r.Equal(0, sut.GetFilesCount())
	r.Equal(0, sut.GetFilesTotal())
	r.Equal(int64(0), sut.GetSizeLeft())
}

func TestExecutionFileIterator_1xFile(t *testing.T) {
	r := require.New(t)

	executionPlan := addFile(new(ExecutionPlan), 1)

	sut := executionPlan.GetFileInterator()

	r.False(sut.HasNext())
	r.NotNil(sut.Current())
	r.Equal("file_0", sut.Current().fileName)
	r.Equal(1, sut.GetFilesCount())
	r.Equal(1, sut.GetFilesTotal())
	r.Equal(int64(100), sut.GetSizeLeft())

	r.Nil(sut.Next())
	r.False(sut.HasNext())
	r.Nil(sut.Current())
	r.Equal(1, sut.GetFilesCount())
	r.Equal(1, sut.GetFilesTotal())
	r.Equal(int64(0), sut.GetSizeLeft())
}

func TestExecutionFileIterator_3xFile(t *testing.T) {
	r := require.New(t)

	executionPlan := addFile(new(ExecutionPlan), 3)

	sut := executionPlan.GetFileInterator()

	r.True(sut.HasNext())
	r.NotNil(sut.Current())
	r.Equal("file_0", sut.Current().fileName)
	r.Equal(1, sut.GetFilesCount())
	r.Equal(3, sut.GetFilesTotal())
	r.Equal(int64(600), sut.GetSizeLeft())

	r.NotNil(sut.Next())
	r.True(sut.HasNext())
	r.NotNil(sut.Current())
	r.Equal("file_1", sut.Current().fileName)
	r.Equal(2, sut.GetFilesCount())
	r.Equal(3, sut.GetFilesTotal())
	r.Equal(int64(500), sut.GetSizeLeft())

	r.NotNil(sut.Next())
	r.False(sut.HasNext())
	r.NotNil(sut.Current())
	r.Equal("file_2", sut.Current().fileName)
	r.Equal(3, sut.GetFilesCount())
	r.Equal(3, sut.GetFilesTotal())
	r.Equal(int64(300), sut.GetSizeLeft())

	r.Nil(sut.Next())
	r.False(sut.HasNext())
	r.Nil(sut.Current())
	r.Equal(3, sut.GetFilesCount())
	r.Equal(3, sut.GetFilesTotal())
	r.Equal(int64(0), sut.GetSizeLeft())
}

func addFile(executionPlan *ExecutionPlan, files int) *ExecutionPlan {
	for indx := 0; indx < files; indx++ {
		wpdObject := &gowpd.Object{}
		wpdObject.Size = int64((indx + 1) * 100)
		fileName := fmt.Sprintf("file_%v", indx)
		file := wpdFile{filePath: fileName, fileName: fileName, wpdObject: wpdObject}
		executionPlan.AddFile(&file)
	}
	return executionPlan
}

func setTotalSize(executionPlan *ExecutionPlan, totalSize int) *ExecutionPlan {
	wpdObject := &gowpd.Object{}
	wpdObject.Size = int64(totalSize)
	file := wpdFile{wpdObject: wpdObject}
	executionPlan.AddFile(&file)
	return executionPlan
}
