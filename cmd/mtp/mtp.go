package mtp

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/tobwithu/gowpd"

	"github.com/cheggaaa/pb/v3"
)

var ProgressTemplate pb.ProgressBarTemplate = `{{with string . "prefix"}}{{.}} {{end}}{{counters . "%s/%s" "%s/?"}} ({{speed . "%s/s" "..."}}) {{bar . }} {{percent . "%.0f%%" "?"}} {{rtime . "ETA %s"}}{{with string . "suffix"}} {{.}}{{end}}`

type MtpDownloader struct {
	resultDir string
	tmpDir    string
	error     error
	dryRun    bool

	currentDeviceId    int
	currentDeviceLabel string
	currentDevice      *gowpd.Device
}

func LoadFromAllWpd(deviceDir string, targetDir string, dryRun bool) (string, error) {
	result := MtpDownloader{dryRun: dryRun}
	result.init(targetDir)
	defer result.close()

	result.loadFromAllDevices(deviceDir)

	return result.GetResultDir(), result.GetError()
}

func LoadFromMatchedWpd(deviceDescriptionFilter string, deviceDir string, targetDir string, dryRun bool) (string, error) {
	result := MtpDownloader{dryRun: dryRun}
	result.init(targetDir)
	defer result.close()

	result.loadFromMatchedDevices(deviceDescriptionFilter, deviceDir)

	return result.GetResultDir(), result.GetError()
}

func (downloader *MtpDownloader) HasError() bool {
	return downloader.GetError() != nil
}

func (downloader *MtpDownloader) GetError() error {
	return downloader.error
}

func (downloader *MtpDownloader) GetResultDir() string {
	if downloader.HasError() {
		return ""
	}
	return downloader.resultDir
}

func (downloader *MtpDownloader) init(targetDir string) {
	downloader.error = gowpd.Init()

	if !downloader.HasError() {
		downloader.resultDir = downloader.generateTmpDir(targetDir)
		downloader.error = os.MkdirAll(downloader.resultDir, os.ModeDir)
	}
}

func (downloader *MtpDownloader) close() {
	gowpd.Destroy()
}

func (downloader *MtpDownloader) loadFromAllDevices(deviceDir string) {
	if downloader.HasError() {
		return
	}

	mtpDeviceCount := gowpd.GetDeviceCount()
	for i := 0; i < mtpDeviceCount; i++ {
		downloader.initCurrentDevice(i)
		if downloader.HasError() {
			//return
			log.Warningf("Unable to read %s!", downloader.currentDeviceLabel)
			continue
		}

		downloader.copyContentToTempDir(deviceDir)
	}
}

func (downloader *MtpDownloader) loadFromMatchedDevices(deviceDescriptionFilter string, deviceDir string) {
	if downloader.HasError() {
		return
	}

	mtpDeviceCount := gowpd.GetDeviceCount()
	for i := 0; i < mtpDeviceCount; i++ {
		downloader.initCurrentDevice(i)
		if downloader.HasError() {
			//return
			log.Warningf("Unable to read %s!", downloader.currentDeviceLabel)
			continue
		}

		//TODO add filter by file name and filter by file/directory exists
		if !strings.Contains(downloader.currentDeviceLabel, deviceDescriptionFilter) {
			log.Infof("Skipping %s because does not match name filter", downloader.currentDeviceLabel)
			continue
		}

		downloader.copyContentToTempDir(deviceDir)
	}
}

func (downloader *MtpDownloader) copyContentToTempDir(deviceDir string) {
	downloader.prepareTempDir()
	if downloader.HasError() {
		log.Warningf("Unable to create '%v' temp directory. %s was skipped", downloader.tmpDir, downloader.currentDeviceLabel)
		return
	}

	wpdRootDirName := gowpd.PathSeparator + deviceDir
	wpdRootDirs := listWpdDir(downloader.currentDevice, wpdRootDirName)

	if len(wpdRootDirs) > 0 {
		log.Infof("Scanning %s...", downloader.currentDeviceLabel)

		executionPlan := BuildExecutionPlan(wpdRootDirs)
		log.Infof("%v file(s) (%v) will be downloaded into '%v' temp directory", executionPlan.GetFilesCount(), executionPlan.GetTotalSizeString(), downloader.tmpDir)

		progressBar := ProgressTemplate.Start64(executionPlan.GetTotalSize())
		defer progressBar.Finish()

		downloader.copyToTmpDir(executionPlan, wpdRootDirName, progressBar)

		if !downloader.dryRun {
			downloader.removeFromWpd(wpdRootDirs)
		}
	}
}

func (downloader *MtpDownloader) prepareTempDir() {
	downloader.tmpDir = path.Join(downloader.resultDir, fmt.Sprint(downloader.currentDeviceId))
	downloader.error = os.MkdirAll(downloader.tmpDir, os.ModeDir)
}

func (downloader *MtpDownloader) initCurrentDevice(i int) {
	downloader.currentDeviceId = i
	downloader.currentDeviceLabel = downloader.buildMtpDeviceLabel()
	log.Infof("Found %s device", downloader.currentDeviceLabel)

	downloader.currentDevice, downloader.error = gowpd.ChooseDevice(downloader.currentDeviceId)
}

func (downloader *MtpDownloader) buildMtpDeviceLabel() string {
	return fmt.Sprintf("MTP#%v - '%v (%v)'", downloader.currentDeviceId, gowpd.GetDeviceName(downloader.currentDeviceId), gowpd.GetDeviceDescription(downloader.currentDeviceId))
}

func (downloader *MtpDownloader) generateTmpDir(targetDir string) string {
	return filepath.Join(targetDir, time.Now().Format("20060102_150405"))
}

func listWpdDir(dev *gowpd.Device, dir string) []*wpdFile {
	obj := dev.FindObject(dir)
	if obj == nil {
		log.Debugf("%v was not found.", dir)
		return make([]*wpdFile, 0)
	}

	wpdFile := newWpdFile(filepath.Dir(dir), dev, obj)
	return wpdFile.chidren
}

func sizeToLabel(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(size)/float64(div), "KMGTPE"[exp])
}

func (downloader *MtpDownloader) copyToTmpDir(executionPlan *ExecutionPlan, wpdRootDir string, progressBar *pb.ProgressBar) {
	fileIterator := executionPlan.GetFileInterator()
	for fileIterator.Current() != nil {
		wpdFile := fileIterator.Current()

		relWpdFilePath := wpdFile.relPath(wpdRootDir)
		targetFile := filepath.Join(downloader.tmpDir, relWpdFilePath)

		log.Debugf("Copying from '%v' to %v... ", wpdFile.filePath, targetFile)
		progressBar.Set("prefix", fmt.Sprintf("(%v/%v) '%v'", fileIterator.GetFilesCount(), fileIterator.GetFilesTotal(), relWpdFilePath))

		targetDir := filepath.Dir(targetFile)
		os.MkdirAll(targetDir, os.ModeDir)

		copyCount, error := wpdFile.copyTo(targetFile, progressBar)

		if error != nil {
			log.Infof("Copy of '%v' - filed - %v", wpdFile.filePath, error)
		} else {
			log.Debugf("Copy of '%v' - done ('%v')", wpdFile.filePath, sizeToLabel(copyCount))
			wpdFile.wasCopied = true
		}

		fileIterator.Next()
	}
}

func (downloader *MtpDownloader) removeFromWpd(filesToRemove []*wpdFile) {
	for _, file := range filesToRemove {
		log.Debugf("Deleting '%v'...", file.filePath)

		if file.wasCopied {
			err := file.deleteFile()
			if err != nil {
				log.Infof("Deleting of '%v' - failed: %v", file.filePath, err)
			} else {
				log.Debugf("Deleting of '%v' - done", file.filePath)
			}
		} else {
			log.Infof("Deleting of '%v' - skipped", file.filePath)
		}

		downloader.removeFromWpd(file.chidren)
	}
}
