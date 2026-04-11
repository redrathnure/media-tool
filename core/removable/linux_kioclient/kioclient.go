//go:build !windows
// +build !windows

package linux_kioclient

import (
	"bufio"
	"io"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type KioClient struct {
}

func NewKioClient() *KioClient {
	result := KioClient{}

	return &result
}

func (c *KioClient) CanWork() bool {
	_, err := c.execS("--version")
	return err == nil
}

func (c *KioClient) Ls(url string) []string {
	result, _ := c.execS("ls", url)
	result = c.removeEmpty(result)
	return result
}

func (c *KioClient) Copy(srcUrl string, dstFile string) error {
	dstFile, err := filepath.Abs(dstFile)
	if err != nil {
		return err
	}
	dstUrl := "file:" + dstFile

	_, err = c.exec("copy", srcUrl, dstUrl)
	return err
}

func (c *KioClient) Rm(url string) {
	c.exec("rm", url)
}

func (c *KioClient) GetState(url string) (isDir bool, size int64) {
	statLines, _ := c.exec("stat", url)
	typeString := c.extract("MIME_TYPE", statLines)
	if strings.Contains(typeString, "directory") {
		isDir = true
		size = 0
		return isDir, size
	}

	sizeString := c.extract("SIZE", statLines)
	result, err := strconv.ParseInt(sizeString, 10, 64)
	if err != nil {
		result = 0
	}
	return isDir, result
}

func (c *KioClient) extract(propName string, resultLines []string) string {
	for _, line := range resultLines {
		lineParts := strings.SplitN(line, " ", 2)
		if len(lineParts) != 2 {
			continue
		}
		key := strings.TrimSpace(lineParts[0])
		val := strings.TrimSpace(lineParts[1])
		if strings.EqualFold(key, propName) {
			return val
		}
	}
	return ""
}

func (c *KioClient) exec(args ...string) (result []string, err error) {
	return c.execInt(false, args...)
}

func (c *KioClient) execS(args ...string) (result []string, err error) {
	return c.execInt(true, args...)
}

func (c *KioClient) execInt(isSilent bool, args ...string) (result []string, err error) {
	cmd := exec.Command("kioclient", args...)
	log.Debugf("kioclient command: '%s'", cmd.String())

	cOut := make(chan []string, 1)
	cErr := make(chan []string, 1)

	stdOutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()
	go c.trackStdOut(&stdOutPipe, cOut, false, isSilent)
	go c.trackStdOut(&stderrPipe, cErr, true, isSilent)

	if err := cmd.Run(); err != nil {
		if !isSilent {
			log.Warningf("kioclient exec error: '%s'", err)
		}
		return result, nil
	}

	result = <-cOut
	return result, nil
}

func (tool *KioClient) trackStdOut(stdOutPipe *io.ReadCloser, ret chan<- []string, isErrorOut bool, isSilent bool) {
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

func (c *KioClient) removeEmpty(lines []string) []string {
	result := []string{}

	for _, line := range lines {
		if line != "" && line != "." {
			result = append(result, line)
		}
	}
	return result
}
