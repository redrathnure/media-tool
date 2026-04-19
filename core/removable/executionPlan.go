package removable

import (
	"fmt"

	"github.com/redrathnure/media-tool/core/removable/core"
)

type ExecutionPlan struct {
	files     []string
	totalSize int64
}

func (p *ExecutionPlan) GetFilesCount() int {
	return len(p.files)
}

func (p *ExecutionPlan) GetTotalSize() int64 {
	return p.totalSize
}

func (p *ExecutionPlan) IsEmpty() bool {
	return len(p.files) == 0
}

func (p *ExecutionPlan) GetTotalSizeString() string {
	size := p.totalSize
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

func BuildExecutionPlan(dev *core.RemovableDevice, deviceDir string) (*ExecutionPlan, error) {
	var result = ExecutionPlan{files: []string{}}

	if err := result.addChildren(dev, deviceDir); err != nil {
		return nil, err
	}

	return &result, nil
}

func (p *ExecutionPlan) addChildren(dev *core.RemovableDevice, deviceDir string) error {
	log.Infof("Scanning %s -> '%s'...", (*dev).Name(), deviceDir)

	children, err := (*dev).GetChildren(deviceDir)
	if err != nil {
		return err
	}

	log.Debugf("'%s' dir has %d children...", deviceDir, len(children))

	for _, child := range children {
		if child.IsDir {
			p.addChildren(dev, child.Name)
		} else {
			p.files = append(p.files, child.Name)
			p.totalSize += child.Size
		}
	}

	return nil
}
