package removable

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/cheggaaa/pb/v3"

	"github.com/redrathnure/media-tool/core/removable/core"
	"github.com/redrathnure/media-tool/core/tools"
)

var CopyProgressTemplate pb.ProgressBarTemplate = `{{with string . "prefix"}}{{.}} {{end}}{{counters . "%s/%s" "%s/?"}} ({{speed . "%s/s" "..."}}) {{bar . }} {{percent . "%.0f%%" "?"}} {{rtime . "ETA %s"}}{{with string . "suffix"}} {{.}}{{end}}`
var DeletingProgressTemplate pb.ProgressBarTemplate = `{{with string . "prefix"}}{{.}} {{end}}{{counters . "%s/%s" "%s/?"}} {{bar . }} {{percent . "%.0f%%" "?"}} {{rtime . "ETA %s"}}{{with string . "suffix"}} {{.}}{{end}}`

type MtpDownloader struct {
	resultDir string
	tmpDir    string
	dryRun    bool
}

func LoadSdPhotos(targetDir string, dryRun bool) (contentDir string, err error) {
	return loadFromAllWpd(core.SdPhotosFilter, core.DCIM_DIR, targetDir, dryRun)
}

func LoadGoProVideos(targetDir string, dryRun bool) (contentDir string, err error) {
	return loadFromAllWpd(core.GoProFilter, core.GOPRO_DIR, targetDir, dryRun)
}

func LoadCamVideos(targetDir string, dryRun bool) (contentDir string, err error) {
	return loadFromAllWpd(core.CamFilter, core.CAM_FILES_DIR, targetDir, dryRun)
}

func loadFromAllWpd(deviceFilter core.DeviceFilter, deviceDir string, targetDir string, dryRun bool) (contentDir string, err error) {
	result := MtpDownloader{dryRun: dryRun}
	defer result.close()
	if err := result.init(targetDir); err != nil {
		return "", err
	}

	return result.loadFromMatchedDevices(deviceFilter, deviceDir)
}

func (d *MtpDownloader) init(targetDir string) error {
	d.resultDir = filepath.Join(targetDir, time.Now().Format("20060102_150405"))
	d.resultDir = tools.ExpandPath(d.resultDir)

	if !d.dryRun {
		return os.MkdirAll(d.resultDir, os.ModePerm)
	}
	return nil
}

func (d *MtpDownloader) loadFromMatchedDevices(deviceFilter core.DeviceFilter, deviceDir string) (contentDir string, err error) {
	devices := d.findDevices(deviceFilter)
	for i, dev := range devices {
		if err := d.copyContentToTempDir(i, dev, deviceDir); err != nil {
			return d.resultDir, err
		}
	}
	return d.resultDir, nil
}

func (d *MtpDownloader) copyContentToTempDir(devIndex int, dev *core.RemovableDevice, deviceDir string) error {
	if err := d.prepareTempDir(devIndex, dev); err != nil {
		return err
	}

	executionPlan, err := BuildExecutionPlan(dev, deviceDir)
	if err != nil {
		return err
	}
	if executionPlan.IsEmpty() {
		log.Infof("'%v' device dir is empty", deviceDir)
		return nil
	}
	log.Infof("%v file(s) (%v) will be downloaded to '%v' temp directory", executionPlan.GetFilesCount(), executionPlan.GetTotalSizeString(), d.tmpDir)

	if err := d.copyToTmpDir(dev, executionPlan); err != nil {
		return err
	}

	if err := d.removeSrcFiles(dev, executionPlan); err != nil {
		return err
	}

	return nil
}

func (d *MtpDownloader) removeSrcFiles(dev *core.RemovableDevice, executionPlan *ExecutionPlan) error {
	if d.dryRun {
		log.Infof("Source files will not be removed ('DryRun' flag is true)")
		return nil
	}

	log.Infof("Deleting origin files from %v", (*dev).Name())

	filesTotal := len(executionPlan.files)
	progressBar := DeletingProgressTemplate.Start(filesTotal)
	defer progressBar.Finish()

	for i, srcFile := range executionPlan.files {
		progressBar.Set("prefix", fmt.Sprintf("(%v/%v) '%v'", i+1, filesTotal, srcFile))

		if err := (*dev).DeleteFile(srcFile); err != nil {
			log.Warningf("Unable to remove'%v' file: %v", srcFile, err)
		}
		progressBar.Increment()
	}
	return nil
}

func (d *MtpDownloader) prepareTempDir(devIndex int, dev *core.RemovableDevice) error {
	d.tmpDir = path.Join(d.resultDir, fmt.Sprintf("%v_%v", devIndex, (*dev).Name()))
	return os.MkdirAll(d.tmpDir, os.ModePerm)
}

func (d *MtpDownloader) copyToTmpDir(dev *core.RemovableDevice, executionPlan *ExecutionPlan) error {
	progressBar := CopyProgressTemplate.Start64(executionPlan.GetTotalSize())
	defer progressBar.Finish()

	filesTotal := len(executionPlan.files)
	for i, srcFile := range executionPlan.files {
		dstFile := filepath.Join(d.tmpDir, srcFile)

		log.Debugf("Copying from '%v' to %v... ", srcFile, dstFile)
		progressBar.Set("prefix", fmt.Sprintf("(%v/%v) '%v'", i+1, filesTotal, srcFile))

		if !d.dryRun {
			dstDir := filepath.Dir(dstFile)
			if err := os.MkdirAll(dstDir, os.ModePerm); err != nil {
				return err
			}

			bytesCount, err := (*dev).CopyFile(srcFile, dstFile, progressBar)

			if err != nil {
				log.Errorf("Unable to copy '%v' file: %v", srcFile, err)
				return err
			} else {
				log.Debugf("Copy of '%v' - done ('%v')", srcFile, tools.HumanizeFileSize(bytesCount))
			}
		}
	}
	return nil
}
