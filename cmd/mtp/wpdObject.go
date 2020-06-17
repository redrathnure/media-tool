package mtp

import (
	"fmt"
	"path/filepath"

	"github.com/tobwithu/gowpd"
)

const (
	//TODO take care about recycle bin
	systemDir = "System Volume Information"
)

type wpdFile struct {
	filePath  string
	fileName  string
	parentDir string
	wasCopied bool
	wpdObject *gowpd.Object
	wpdDevice *gowpd.Device
	chidren   map[string]*wpdFile
}

func newWpdFile(parentDir string, dev *gowpd.Device, obj *gowpd.Object) wpdFile {
	result := wpdFile{
		wpdObject: obj,
		wpdDevice: dev,
		fileName:  obj.Name,
		parentDir: parentDir,
		filePath:  filepath.Join(parentDir, obj.Name),
		wasCopied: false,
		chidren:   make(map[string]*wpdFile),
	}
	result.initChildren()
	return result
}

func (wf wpdFile) initChildren() {
	objs, _ := wf.wpdDevice.GetChildObjects(wf.wpdObject.Id)
	//TODO error handling

	curPath := wf.filePath
	for _, o := range objs {

		if o.Name == systemDir {
			continue
		}
		rel := filepath.Join(curPath, o.Name)

		fmt.Printf("Reading info: %v \n", rel)

		child := newWpdFile(wf.filePath, wf.wpdDevice, o)
		wf.chidren[child.fileName] = &child
	}
}

func (wf wpdFile) relPath(basepath string) string {
	result, err := filepath.Rel(basepath, wf.filePath)
	if err != nil {
		fmt.Printf("Unable to calculate relative path for %v regarding %v \n", wf.filePath, basepath)
		result = wf.filePath
	}
	return result
}

func (wf wpdFile) copyTo(targetFile string) (int64, error) {
	return wf.wpdDevice.CopyObjectFromDevice(targetFile, wf.wpdObject)
}

func (wf wpdFile) deleteFile() error {
	if !wf.wpdObject.IsDir {
		return wf.wpdDevice.Delete(wf.wpdObject.Id)
	}
	return nil
}
