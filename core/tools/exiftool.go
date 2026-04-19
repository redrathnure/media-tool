package tools

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const (
	whatsAppMarker        = "WhatsApp"
	transcodedVideoMarker = "_x265"
)

type ExifToolWrapper struct {
	cmd         string
	defaultArgs []string
	args        ExifToolArgs
	execCommand func(name string, args ...string) *exec.Cmd
}

func (tool *ExifToolWrapper) GetFileNameTag(dryRun bool) string {
	if dryRun {
		return "TestName"
	}
	return "FileName"
}

type ExifToolArgs struct {
	args []string
}

var exifToolObj *ExifToolWrapper

var ExifToolPath = ""

func newExifTool() *ExifToolWrapper {
	result := ExifToolWrapper{
		cmd:         "exiftool",
		defaultArgs: []string{"-v0", "-progress"},
		execCommand: exec.Command,
	}
	result.initCmd()
	result.NewArgs()
	return &result
}

func GetExifTool() *ExifToolWrapper {
	if exifToolObj == nil {
		exifToolObj = newExifTool()
	}
	return exifToolObj
}

func (tool *ExifToolWrapper) initCmd() {
	customPath := ExifToolPath
	if customPath != "" {

		customPath = ExpandPath(customPath)

		if _, err := os.Stat(customPath); err != nil {
			log.Infof("Unable to find custom exiftool: '%s'. Trying to use '%s' from $PATH", err, tool.cmd)
			return
		}

		tool.cmd = customPath
	}
}

func (tool *ExifToolWrapper) Exec(verbose bool) {
	tool.args.filterPlatformSpecific()
	cmd := tool.execCommand(tool.cmd, tool.args.args...)

	log.Debugf("ExifTool command: '%s'\n", cmd.String())

	stdOutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()
	go tool.trackStdOut(&stdOutPipe, false)
	go tool.trackStdOut(&stderrPipe, true)

	err := cmd.Run()
	if err != nil {
		if cmd.ProcessState.ExitCode() != 2 {
			log.Warningf("ExifTool exec error: '%s'", err)
		}
	}
}

func (tool *ExifToolWrapper) trackStdOut(stdOutPipe *io.ReadCloser, isErrorOut bool) {
	scanner := bufio.NewScanner(*stdOutPipe)
	for scanner.Scan() {
		line := scanner.Text()

		if isErrorOut {
			log.Warning(line)
		} else {
			log.Info(line)
		}

	}
}

func (tool *ExifToolWrapper) NewArgs() *ExifToolArgs {
	tool.args = ExifToolArgs{args: tool.defaultArgs}
	return &tool.args
}

func (tool *ExifToolWrapper) DeleteOriginals(deleteOriginal bool) {
	if deleteOriginal {
		tool.defaultArgs = append(tool.defaultArgs, "-overwrite_original")
	}
}

func (toolArgs *ExifToolArgs) add(args ...string) {
	toolArgs.args = append(toolArgs.args, args...)
}

func (toolArgs *ExifToolArgs) Recursively(recursively bool) {
	if recursively {
		toolArgs.add("-r")
	}
}

func (toolArgs *ExifToolArgs) Src(dirOrFilepath string) {
	toolArgs.add(dirOrFilepath)
}

func (toolArgs *ExifToolArgs) ForImages() {
	toolArgs.add("-ext", "jpg")
	toolArgs.add("-ext", "jpeg")
	toolArgs.add("-ext", "nef")
	toolArgs.add("-ext", "cr2")
	toolArgs.add("-ext", "cr3")
}

func (toolArgs *ExifToolArgs) ForVideoMp4() {
	toolArgs.add("-ext", "mp4")
}

func (toolArgs *ExifToolArgs) ForVideoLrv() {
	toolArgs.add("-ext", "LRV")
}

func (toolArgs *ExifToolArgs) ForVideoAvchd() {
	toolArgs.add("-ext", "mts")
}

func (toolArgs *ExifToolArgs) ForDateFormat(dateFormat string) {
	toolArgs.add("-d", dateFormat)
}

func (toolArgs *ExifToolArgs) CopyTag(dstTagName string, srcTagName string) {
	toolArgs.add(fmt.Sprintf("-%s<%s", dstTagName, srcTagName))
}

func (toolArgs *ExifToolArgs) SetTag(tagName string, tagValue string) {
	toolArgs.add(fmt.Sprintf("-%s=%s", tagName, tagValue))
}

func (toolArgs *ExifToolArgs) ChangeFileDate(tagName string) {
	//File:
	toolArgs.CopyTag("FileModifyDate", tagName)
	toolArgs.CopyTag("FileCreateDate", tagName)
}

func (toolArgs *ExifToolArgs) ChangeExifDate(tagName string) {
	//'EXIF:
	toolArgs.CopyTag("CreateDate", tagName)
	toolArgs.CopyTag("DateTimeOriginal", tagName)
}

func (toolArgs *ExifToolArgs) ChangeMp4Date(tagName string) {
	//quicktime:
	toolArgs.CopyTag("CreateDate", tagName)
	toolArgs.CopyTag("ModifyDate", tagName)
	toolArgs.CopyTag("TrackCreateDate", tagName)
	toolArgs.CopyTag("TrackModifyDate", tagName)
	toolArgs.CopyTag("MediaCreateDate", tagName)
	toolArgs.CopyTag("MediaModifyDate", tagName)
}

func (toolArgs *ExifToolArgs) CleanTag(tagName string) {
	toolArgs.add(fmt.Sprintf("-%s=", tagName))
}

func (toolArgs *ExifToolArgs) CleanVendorTags() {
	toolArgs.CleanTag("Software")
	toolArgs.CleanTag("WriterName")
	toolArgs.CleanTag("ReaderName")
	toolArgs.CleanTag("HistorySoftwareAgent")
	toolArgs.CleanTag("LookCopyright")
	toolArgs.CleanTag("XMPToolkit")
	toolArgs.CleanTag("photoshop:all")
	toolArgs.CleanTag("NikonCapture:all")
	toolArgs.CleanTag("GIMP:all")
	toolArgs.CleanTag("history*")
}

func (toolArgs *ExifToolArgs) CleanCameraTags() {
	// Camera vendor specific
	toolArgs.CleanTag("Canon:all")
	toolArgs.CleanTag("Sony:all")
	toolArgs.CleanTag("GoPro:all")
	toolArgs.CleanTag("Nikon:all")
	toolArgs.CleanTag("FujiFilm:all")
	toolArgs.CleanTag("HP:all")
	toolArgs.CleanTag("Kodak:all")
	toolArgs.CleanTag("Minolta:all")
	toolArgs.CleanTag("Nintendo:all")
	toolArgs.CleanTag("Olympus:all")
	toolArgs.CleanTag("Panasonic:all")
	toolArgs.CleanTag("Pentax:all")
	toolArgs.CleanTag("Samsung:all")
	toolArgs.CleanTag("Sanyo:all")
	toolArgs.CleanTag("Sigma:all")
	toolArgs.CleanTag("Sony:all")
	toolArgs.CleanTag("CanonRaw:all")
	toolArgs.CleanTag("MinoltaRaw:all")
	toolArgs.CleanTag("PanasonicRaw:all")
	toolArgs.CleanTag("SigmaRaw:all")

	// Common shot parameters
	toolArgs.CleanTag("all:canonexposuremode")
	toolArgs.CleanTag("EXIF:Make")
	toolArgs.CleanTag("EXIF:Model")
	toolArgs.CleanTag("EXIF:FNumber")
	toolArgs.CleanTag("Exposure*")
	toolArgs.CleanTag("ISO")
	toolArgs.CleanTag("Lens*")
	toolArgs.CleanTag("Focal*")
	toolArgs.CleanTag("Flash*")
	toolArgs.CleanTag("Camera*")
	toolArgs.CleanTag("Metering*")
	toolArgs.CleanTag("Shutter*")
	toolArgs.CleanTag("Megapixels*")
	toolArgs.CleanTag("HasCrop")
	toolArgs.CleanTag("Format")

}

func (toolArgs *ExifToolArgs) CleanLocationTags() {
	toolArgs.CleanTag("gps:all")
}

func (toolArgs *ExifToolArgs) ExcludeIfNameContains(fileNameFragments ...string) {
	if len(fileNameFragments) > 0 {
		toolArgs.add("-if", "$filename !~ /"+strings.Join(fileNameFragments, "|")+"/i")
	}
}

func (toolArgs *ExifToolArgs) ExcludeWhatsAppFiles() {
	toolArgs.ExcludeIfNameContains(whatsAppMarker)
}

func (toolArgs *ExifToolArgs) IncludeWhatsAppFiles() {
	toolArgs.add("-if", "$filename =~ /"+whatsAppMarker+" /i")
}

func (toolArgs *ExifToolArgs) ExcludeTranscodedVideoFiles() {
	toolArgs.ExcludeIfNameContains(transcodedVideoMarker)
}

func (toolArgs *ExifToolArgs) filterPlatformSpecific() {
	if runtime.GOOS != "windows" && runtime.GOOS != "darwin" {
		newArgs := []string{}
		for _, element := range toolArgs.args {
			if !strings.Contains(element, "FileCreateDate") {
				newArgs = append(newArgs, element)
			}
		}
		toolArgs.args = newArgs
	}
}
