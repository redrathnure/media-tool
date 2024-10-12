//go:build mage
// +build mage

package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	buildDir   = "./build/"
	releaseDir = "./build/release/"
	distDir    = "./dist/"
)

// Clean go modules
func GoClean() error {
	return sh.RunV("go", "clean")
}

// Clean go modules
func GoUpdateDeps() error {
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

// Clean build dir
func Clean() {
	os.RemoveAll(buildDir)
}

// Build project
func Build() error {
	if err := sh.RunV("go", "version"); err != nil {
		return err
	}
	if err := sh.RunV("go", "build"); err != nil {
		return err
	}
	return nil
}

// Clean and build project
func ReBuild() {
	mg.Deps(GoClean, Clean, Build)
}

// Prepare release package
func ReleasePkg() error {
	version, err := getGitVersion()
	if err != nil {
		return err
	}

	fmt.Printf("Building %s release\n", version)

	ReBuild()

	if err := prepareReleaseDir(); err != nil {
		return err
	}
	if err := buildReleasePackage(version); err != nil {
		return err
	}

	return nil
}

func getGitVersion() (string, error) {
	return sh.Output("git", "describe", "--tags")
}

func prepareReleaseDir() error {

	if err := os.MkdirAll(releaseDir, 0755); err != nil {
		return err
	}
	if err := copyToDir("LICENSE", releaseDir); err != nil {
		return err
	}
	if err := copyToDir("README.md", releaseDir); err != nil {
		return err
	}
	if err := copyToDir("media-tool.exe", releaseDir); err != nil {
		return err
	}
	if err := copyToDir("media-tool.example.yml", releaseDir); err != nil {
		return err
	}
	return nil
}

func copyToDir(fileName, dstDir string) error {
	return os.Link(fileName, dstDir+fileName)
}

func buildReleasePackage(version string) error {
	releaseFile := fmt.Sprintf("%s/media-tool_%s_x64.zip", distDir, version)
	fmt.Printf("Building '%s' archive\n", releaseFile)

	os.MkdirAll(distDir, 0755)
	return zipDir(releaseDir, releaseFile)
}

func zipDir(sourceDir, targetFile string) error {
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
