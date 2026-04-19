//go:build windows
// +build windows

package windows_removable

import (
	"bufio"
	"io"
	"os/exec"
	"strings"

	"github.com/redrathnure/media-tool/core/removable/core"
)

func FindDevices(deviceFilter core.DeviceFilter) []*core.RemovableDevice {
	result := []*core.RemovableDevice{}

	log.Infof("Checking Win removable devices...")

	drives, err := getWmicResult()
	if err != nil {
		log.Warningf("Unable to execute 'wmic logicaldisk where drivetype=2 get deviceid' command: %s", err)
		return result
	}

	for _, drive := range drives {

		udev := NewWinDevice(drive)
		log.Infof("Found removable device: '%s'", udev.Name())
		if deviceFilter.Accept(&udev) {
			result = append(result, &udev)
		} else {
			log.Infof("Device was rejected by filters. Skipping...")
		}
	}
	return result
}

func getWmicResult() (drives []string, err error) {
	drives = []string{}

	cmd := exec.Command("wmic", "logicaldisk", "where", "drivetype=2", "get", "deviceid")
	log.Debugf("wmic command: '%s'", cmd.String())

	isSilent := false

	cOut := make(chan []string, 1)
	cErr := make(chan []string, 1)

	stdOutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()
	go trackStdOut(&stdOutPipe, cOut, false, isSilent)
	go trackStdOut(&stderrPipe, cErr, true, isSilent)

	if err := cmd.Run(); err != nil {
		if !isSilent {
			log.Warningf("wmic exec error: '%s'", err)
		}
		return drives, nil
	}

	result := <-cOut
	for _, resultLine := range result {
		resultLine := strings.TrimSpace(resultLine)
		if strings.HasSuffix(resultLine, ":") {
			drives = append(drives, resultLine)
		}
	}
	return drives, nil
}

func trackStdOut(stdOutPipe *io.ReadCloser, ret chan<- []string, isErrorOut bool, isSilent bool) {
	result := []string{}
	scanner := bufio.NewScanner(*stdOutPipe)
	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, line)

		if isErrorOut && !isSilent {
			log.Warning(line)
		} else {
			log.Debugf(line)
		}

	}
	ret <- result
}
