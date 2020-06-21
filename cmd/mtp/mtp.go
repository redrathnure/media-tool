package mtp

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/tobwithu/gowpd"
)

func LoadFromAllWpd(deviceDir string, targetDir string, removeFromOrigin bool) (string, error) {
	err := gowpd.Init()
	defer gowpd.Destroy()

	if err != nil {
		return "", err
	}
	resultDir := generateTmpDir(targetDir)
	err = os.MkdirAll(resultDir, os.ModeDir)
	if err != nil {
		return "", err
	}

	mtpDeviceCount := gowpd.GetDeviceCount()
	for i := 0; i < mtpDeviceCount; i++ {
		mtpLabel := buildMtpDeviceLabel(i)
		log.Infof("Checking MTP#%v - '%s'...", i, mtpLabel)

		dev, err := gowpd.ChooseDevice(i)
		if err != nil {
			return "", err
		}

		tmpDir := path.Join(resultDir, fmt.Sprint(i))
		err = os.MkdirAll(tmpDir, os.ModeDir)
		if err != nil {
			log.Warningf("Unable to create '%v' temp directory. MTP#%v - '%s' was skipped", tmpDir, i, mtpLabel)
			continue
		}

		log.Infof("Files from '%s' will be downloaded into '%v' temp directory", mtpLabel, tmpDir)

		wpdRootDir := gowpd.PathSeparator + deviceDir
		files := listWpdDir(dev, wpdRootDir)

		for _, file := range files {
			copyFromWpd(file, wpdRootDir, tmpDir)
		}

		if removeFromOrigin {
			removeFromWpd(files)
		}
	}

	return resultDir, nil
}

func buildMtpDeviceLabel(deviceNumber int) string {
	return fmt.Sprintf("%v (%v)", gowpd.GetDeviceName(deviceNumber), gowpd.GetDeviceDescription(deviceNumber))
}

func LoadFromWpd(deviceDescriptionFilter string, deviceDir string, targetDir string, removeFromOrigin bool) (string, error) {
	err := gowpd.Init()
	defer gowpd.Destroy()

	if err != nil {
		return "", err
	}

	mtpDeviceCount := gowpd.GetDeviceCount()
	for i := 0; i < mtpDeviceCount; i++ {
		mtpLabel := buildMtpDeviceLabel(i)
		log.Infof("Checking MTP#%v - '%s'...", i, mtpLabel)

		if strings.Contains(gowpd.GetDeviceDescription(i), deviceDescriptionFilter) || strings.Contains(gowpd.GetDeviceName(i), deviceDescriptionFilter) {

			dev, err := gowpd.ChooseDevice(i)
			if err != nil {
				return "", err
			}

			tmpDir := generateTmpDir(targetDir)
			err = os.MkdirAll(tmpDir, os.ModeDir)
			if err != nil {
				log.Warningf("Unable to create '%v' temp directory. MTP#%v - '%s' was skipped", tmpDir, i, mtpLabel)
				return "", err
			}

			log.Infof("Files from '%s' will be downloaded into '%v' temp directory", mtpLabel, tmpDir)

			wpdRootDir := gowpd.PathSeparator + deviceDir
			files := listWpdDir(dev, wpdRootDir)

			for _, file := range files {
				copyFromWpd(file, wpdRootDir, tmpDir)
			}

			if removeFromOrigin {
				removeFromWpd(files)
			}

			return tmpDir, nil
		}

		log.Infof("MTP#%v is not a GoPro device", i)
	}

	return "", fmt.Errorf("No '%v' MTP devices was found", deviceDescriptionFilter)
}

func printWpdFile(file *wpdFile) {
	log.Infof("%v (isDir: %v)", file.filePath, file.wpdObject.IsDir)

	for _, file := range file.chidren {
		printWpdFile(file)
	}
}

func generateTmpDir(targetDir string) string {
	return filepath.Join(targetDir, time.Now().Format("20060102_150405"))
}

func listWpdDir(dev *gowpd.Device, dir string) map[string]*wpdFile {
	obj := dev.FindObject(dir)
	if obj == nil {
		log.Warningf("%v was not found.", dir)
		return make(map[string]*wpdFile)
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

func copyFromWpd(wpdFile *wpdFile, wpdRootDir string, tmpDir string) {
	if wpdFile.wpdObject.IsDir {
		for _, child := range wpdFile.chidren {
			copyFromWpd(child, wpdRootDir, tmpDir)
		}
	} else {
		relWpdFilePath := wpdFile.relPath(wpdRootDir)
		targetFile := filepath.Join(tmpDir, relWpdFilePath)

		log.Debugf("Copying from '%v' to %v... ", wpdFile.filePath, targetFile)
		targetDir := filepath.Dir(targetFile)
		os.MkdirAll(targetDir, os.ModeDir)

		copyCount, error := wpdFile.copyTo(targetFile)

		if error != nil {
			log.Infof("Copy of '%v' - filed - %v", wpdFile.filePath, error)
		} else {
			log.Infof("Copy of '%v' - done ('%v')", wpdFile.filePath, sizeToLabel(copyCount))
			wpdFile.wasCopied = true
		}
	}
}

func removeFromWpd(wpdFiles map[string]*wpdFile) {
	for _, file := range wpdFiles {
		log.Debugf("Deleting '%v'...", file.filePath)

		if file.wasCopied {
			err := file.deleteFile()
			if err != nil {
				log.Infof("Deleting of '%v' - failed: %v", file.filePath, err)
			} else {
				log.Infof("Deleting of '%v' - done", file.filePath)
			}
		} else {
			log.Infof("Deleting of '%v' - skipped", file.filePath)
		}

		removeFromWpd(file.chidren)
	}
}
