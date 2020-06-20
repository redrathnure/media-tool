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
		fmt.Printf("Checking MTP#%v - %v (%v)...\n", i, gowpd.GetDeviceName(i), gowpd.GetDeviceDescription(i))

		dev, err := gowpd.ChooseDevice(i)
		if err != nil {
			return "", err
		}

		tmpDir := path.Join(resultDir, fmt.Sprint(i))
		err = os.MkdirAll(tmpDir, os.ModeDir)
		if err != nil {
			fmt.Printf("Unable to create '%v' temp directory. MTP#%v - %v (%v) was skipped\n", tmpDir, i, gowpd.GetDeviceName(i), gowpd.GetDeviceDescription(i))
			continue
		}

		fmt.Printf("Files will be downloaded into '%v' temp directory\n", tmpDir)

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

	//return "", fmt.Errorf("No '%v' WPD devices was found", deviceDescriptionFilter)
}

func LoadFromWpd(deviceDescriptionFilter string, deviceDir string, targetDir string, removeFromOrigin bool) (string, error) {
	err := gowpd.Init()
	defer gowpd.Destroy()

	if err != nil {
		return "", err
	}

	mtpDeviceCount := gowpd.GetDeviceCount()
	for i := 0; i < mtpDeviceCount; i++ {
		fmt.Printf("Checking MTP#%v - %v (%v)...\n", i, gowpd.GetDeviceName(i), gowpd.GetDeviceDescription(i))

		if strings.Contains(gowpd.GetDeviceDescription(i), deviceDescriptionFilter) || strings.Contains(gowpd.GetDeviceName(i), deviceDescriptionFilter) {

			dev, err := gowpd.ChooseDevice(i)
			if err != nil {
				return "", err
			}

			tmpDir := generateTmpDir(targetDir)
			err = os.MkdirAll(tmpDir, os.ModeDir)
			if err != nil {
				return "", err
			}

			fmt.Printf("Files will be downloaded into '%v' temp directory\n", tmpDir)

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

		fmt.Printf("MTP#%v is not a GoPro device\n", i)
	}

	return "", fmt.Errorf("No '%v' WPD devices was found", deviceDescriptionFilter)
}

func printWpdFile(file *wpdFile) {
	fmt.Printf("%v (isDir: %v)\n", file.filePath, file.wpdObject.IsDir)

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
		fmt.Printf("%v was not found.\n", dir)
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

		fmt.Printf("Copying from '%v' to %v... ", wpdFile.filePath, targetFile)
		targetDir := filepath.Dir(targetFile)
		os.MkdirAll(targetDir, os.ModeDir)

		copyCount, error := wpdFile.copyTo(targetFile)

		if error != nil {
			fmt.Printf("filed - %v\n", error)
		} else {
			fmt.Printf("done ('%v')\n", sizeToLabel(copyCount))
			wpdFile.wasCopied = true
		}
	}
}

func removeFromWpd(wpdFiles map[string]*wpdFile) {
	for _, file := range wpdFiles {
		fmt.Printf("Deleting '%v'...", file.filePath)
		if file.wasCopied {

			err := file.deleteFile()
			if err != nil {
				fmt.Printf(" failed: %v\n", err)
			} else {
				fmt.Printf(" done\n")
			}
		} else {
			fmt.Printf(" skipped\n")
		}

		removeFromWpd(file.chidren)
	}
}
