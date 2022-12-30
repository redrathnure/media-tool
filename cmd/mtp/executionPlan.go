package mtp

import "fmt"

type ExecutionPlan struct {
	files      []*wpdFile
	totalSize  int64
	wpdRootDir string
}

type ExecutionFileIterator struct {
	plan         *ExecutionPlan
	currentIndex int
	sizeLeft     int64
}

func (executionPlan *ExecutionPlan) AddFile(file *wpdFile) {
	executionPlan.files = append(executionPlan.files, file)
	executionPlan.totalSize += file.wpdObject.Size
}

func (executionPlan *ExecutionPlan) GetFilesCount() int {
	return len(executionPlan.files)
}

func (executionPlan *ExecutionPlan) GetTotalSize() int64 {
	return executionPlan.totalSize
}

func (executionPlan *ExecutionPlan) IsEmpty() bool {
	return executionPlan.GetFilesCount() == 0
}

func (executionPlan *ExecutionPlan) GetTotalSizeString() string {
	size := executionPlan.totalSize
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(size)/float64(div), "KMGTPE"[exp])
}

func (executionPlan *ExecutionPlan) GetFileInterator() *ExecutionFileIterator {
	return &ExecutionFileIterator{plan: executionPlan, currentIndex: 0, sizeLeft: executionPlan.GetTotalSize()}
}

func (fileIt *ExecutionFileIterator) Current() *wpdFile {
	if fileIt.currentIndex < len(fileIt.plan.files) {
		return fileIt.plan.files[fileIt.currentIndex]
	}
	return nil
}

func (fileIt *ExecutionFileIterator) HasNext() bool {
	return fileIt.GetFilesTotal() > fileIt.GetFilesCount()
}

func (fileIt *ExecutionFileIterator) Next() *wpdFile {
	if !fileIt.HasNext() {
		fileIt.currentIndex = len(fileIt.plan.files)
		fileIt.sizeLeft = 0
		return nil
	}

	fileIt.sizeLeft -= fileIt.Current().wpdObject.Size

	fileIt.currentIndex += 1

	return fileIt.Current()
}

func (fileIt *ExecutionFileIterator) GetSizeLeft() int64 {
	return fileIt.sizeLeft
}

func (fileIt *ExecutionFileIterator) GetFilesCount() int {
	if fileIt.currentIndex >= fileIt.GetFilesTotal() {
		return fileIt.GetFilesTotal()
	}
	return fileIt.currentIndex + 1
}

func (fileIt *ExecutionFileIterator) GetFilesTotal() int {
	return fileIt.plan.GetFilesCount()
}

func BuildExecutionPlan(rootDirs []*wpdFile, wpdRootDir string) *ExecutionPlan {
	var result = new(ExecutionPlan)
	result.wpdRootDir = wpdRootDir
	for _, file := range rootDirs {
		addToPlan(file, result)
	}
	return result
}

func addToPlan(wpdFile *wpdFile, executionPlan *ExecutionPlan) {
	if wpdFile.wpdObject.IsDir {
		children := wpdFile.chidren

		for _, child := range children {
			addToPlan(child, executionPlan)
		}
	} else {
		executionPlan.AddFile(wpdFile)
	}
}
