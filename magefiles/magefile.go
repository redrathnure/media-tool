//go:build mage
// +build mage

package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	buildDir   = "./build/"
	releaseDir = "./build/release/"
	distDir    = "./dist/"
)

var releasePlatforms = [...]string{"windows/amd64", "windows/arm64", "linux/amd64", "linux/arm64"}

/* Cleanup tasks */

// Remove build dir and all temp files
func Clean() {
	fmt.Printf("Removing '%s' dir...\n", buildDir)
	os.RemoveAll(buildDir)
}

// Clean go modules
func CleanGo() error {
	fmt.Printf("Running 'go clean'...\n")
	return sh.RunV("go", "clean")
}

// Remove all build files and perform cleanGo
func CleanAll() {
	mg.Deps(CleanGo, Clean)
}

/* Project maintenance tasks */

// Update Go dependencies
func UpdateDeps() error {
	fmt.Printf("Updating Go dependencies...\n")
	if err := sh.RunV("go", "get", "-u"); err != nil {
		return err
	}
	if err := sh.RunV("go", "mod", "tidy"); err != nil {
		return err
	}
	if err := sh.RunV("go", "mod", "vendor"); err != nil {
		return err
	}
	return nil
}

/* Build related */

// Build project
func Build(platform *string, // target architecture, e.g. linux/arm64, windows/amd64 or linux
	version *string, // Release version, e.g. 1.2.3
) error {

	os_name, arch := parsePlatform(platform)
	fmt.Printf("Building for %s/%s platform...\n", os_name, arch)

	os.Setenv("GOOS", os_name)
	os.Setenv("GOARCH", arch)

	if err := sh.RunV("go", "version"); err != nil {
		return err
	}

	args := []string{"build"}
	if version != nil && *version != "" {
		args = append(args, "-ldflags", fmt.Sprintf("-X github.com/redrathnure/media-tool/cmd.version=%s", *version))
	}
	if err := sh.RunV("go", args...); err != nil {
		return err
	}
	return nil
}

// Test runs all tests in the project
func Test() error {
	if err := sh.RunV("go", "test", "./..."); err != nil {
		return err
	}
	return nil
}

func TestV() error {
	if err := sh.RunV("go", "test", "./...", "-v"); err != nil {
		return err
	}
	return nil
}

// Clean and build project
func BuildClean(platform *string, // target architecture, e.g. linux/arm64, windows/amd64 or linux
	version *string, // Release version, e.g. 1.2.3
) {
	//mg.Deps(GoClean, Clean, Build)
	Clean()
	Build(platform, version)
}

/* Release related */

// Prepare release package
func Release() error {
	version, err := getGitVersion()
	if err != nil {
		return err
	}

	fmt.Printf("Preparing %s release...\n", version)
	for _, platform := range releasePlatforms {

		fmt.Printf("\nPreparing %s package...\n", platform)
		BuildClean(&platform, &version)

		if err := prepareReleaseDir(platform); err != nil {
			return err
		}
		if err := buildReleasePackage(version, platform); err != nil {
			return err
		}
	}

	return nil
}

/* Helpers*/

func getGitVersion() (string, error) {
	return sh.Output("git", "describe", "--tags")
}

func prepareReleaseDir(platform string) error {
	fmt.Printf("Preparing '%s' dir...\n", releaseDir)

	os_name, _ := parsePlatform(&platform)

	if err := os.MkdirAll(releaseDir, 0755); err != nil {
		return err
	}
	if err := copyToDir("LICENSE", releaseDir); err != nil {
		return err
	}
	if err := copyToDir("README.md", releaseDir); err != nil {
		return err
	}
	var binFile = "media-tool"
	if os_name == "windows" {
		binFile = "media-tool.exe"
	}
	if err := copyToDir(binFile, releaseDir); err != nil {
		return err
	}

	var exampleFile = "media-tool.example.linux.yml"
	if os_name == "windows" {
		exampleFile = "media-tool.example.windows.yml"
	}
	if err := copyToDir2(exampleFile, releaseDir, "media-tool.example.yml"); err != nil {
		return err
	}
	return nil
}

func copyToDir(srcFileName, dstDir string) error {
	return copyToDir2(srcFileName, dstDir, srcFileName)
}

func copyToDir2(srcFileName, dstDir string, dstFileName string) error {
	var dstFile = filepath.Join(dstDir, dstFileName)
	return os.Link(srcFileName, dstFile)
}

func buildReleasePackage(version string, platform string) error {
	os_name, arch := parsePlatform(&platform)

	releaseFile := filepath.Join(distDir, fmt.Sprintf("media-tool_%s_%s-%s.zip", version, os_name, arch))
	fmt.Printf("Building '%s' archive\n", releaseFile)

	os.MkdirAll(distDir, 0755)
	return zipDir(releaseDir, releaseFile)
}

func zipDir(sourceDir, targetFile string) error {
	fmt.Printf("Compressing '%s' dir...\n", sourceDir)

	// 1. Create a ZIP file and zip.Writer
	f, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// 2. Go through all the files of the source
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if sourceDir == path {
			return nil
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(sourceDir), path)
		if err != nil {
			return err
		}
		fmt.Printf(" - adding '%s' file...\n", header.Name)

		if info.IsDir() {
			header.Name += "/"
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}

func parsePlatform(platform *string) (os_name string, arch string) {
	if platform != nil {

		parts := strings.Split(*platform, "/")
		if len(parts) == 1 {
			return parts[0], runtime.GOARCH
		}
		if len(parts) == 2 {
			return parts[0], parts[1]
		}
	}

	return runtime.GOOS, runtime.GOARCH
}
