package mtp

import (
	"io"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
	"github.com/tobwithu/gowpd"
)

var ignoreFiles = []string{"System Volume Information", "$RECYCLE.BIN"}

type wpdFile struct {
	filePath  string
	fileName  string
	parentDir string
	wasCopied bool
	wpdObject *gowpd.Object
	wpdDevice *gowpd.Device
	chidren   []*wpdFile
}

func newWpdFile(parentDir string, dev *gowpd.Device, obj *gowpd.Object) wpdFile {
	result := wpdFile{
		wpdObject: obj,
		wpdDevice: dev,
		fileName:  obj.Name,
		parentDir: parentDir,
		filePath:  filepath.Join(parentDir, obj.Name),
		wasCopied: false,
		chidren:   make([]*wpdFile, 0),
	}
	result.initChildren()
	return result
}

func (wf *wpdFile) initChildren() {
	objs, err := wf.wpdDevice.GetChildObjects(wf.wpdObject.Id)
	if err != nil {
		log.Warningf("Unable to read children for %v: %v", wf.filePath, err)
	}

	curPath := wf.filePath
	for _, o := range objs {

		if isIgnored(o.Name) {
			log.Debugf("Skipping '%v' file", o.Name)
			continue
		}

		rel := filepath.Join(curPath, o.Name)

		log.Debugf("Found: %v", rel)

		child := newWpdFile(wf.filePath, wf.wpdDevice, o)
		wf.chidren = append(wf.chidren, &child)
	}
}

func (wf *wpdFile) relPath(basepath string) string {
	result, err := filepath.Rel(basepath, wf.filePath)
	if err != nil {
		log.Warningf("Unable to calculate relative path for '%v' against to '%v'", wf.filePath, basepath)
		result = wf.filePath
	}
	return result
}

func (wf *wpdFile) copyTo(targetFile string, progressBar *pb.ProgressBar) (int64, error) {
	obj := wf.wpdObject
	id := obj.Id

	reader, err := wf.wpdDevice.GetReader(id)
	if err != nil {
		return 0, err
	}
	defer reader.Close()

	f, err := os.Create(targetFile)
	if err != nil {
		return 0, err
	}
	writer := gowpd.NewBufWriteCloser(f, 0)
	defer writer.Close()

	proxyWriter := progressBar.NewProxyWriter(writer)

	written, err := io.Copy(proxyWriter, reader)

	if err != nil {
		return 0, err
	}
	return written, gowpd.SetFileTime(targetFile, obj.ModTime)

	//return wf.wpdDevice.CopyObjectFromDevice(targetFile, wf.wpdObject)
}

func (wf *wpdFile) deleteFile() error {
	if !wf.wpdObject.IsDir {
		return wf.wpdDevice.Delete(wf.wpdObject.Id)
	}
	return nil
}

func isIgnored(fileName string) bool {
	for _, ignore := range ignoreFiles {
		if fileName == ignore {
			return true
		}
	}
	return false
}
