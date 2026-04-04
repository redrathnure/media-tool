//go:build mage
// +build mage

package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
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

	if version == nil || *version == "" {
		gitVersion, err := getGitVersion()
		if err != nil {
			return fmt.Errorf("failed to get git version: %v", err)
		}
		version = &gitVersion
	}

	args := []string{"build"}
	ldflags := []string{fmt.Sprintf("-X github.com/redrathnure/media-tool/cmd.goVersion=%s", runtime.Version()),
		fmt.Sprintf("-X github.com/redrathnure/media-tool/cmd.goBuildPlatform=%s/%s", runtime.GOOS, runtime.GOARCH),
		fmt.Sprintf("-X github.com/redrathnure/media-tool/cmd.version=%s", *version),
	}

	args = append(args, "-ldflags", strings.Join(ldflags, " "))
	if err := sh.RunV("go", args...); err != nil {
		return err
	}
	return nil
}

// Test runs all tests in the project
func Test() error {
	return runTest(false)
}

func TestV() error {

	return runTest(true)
}

func runTest(verbose bool) error {
	covDir := path.Join(".", "build", "unittest")
	covFile := path.Join(covDir, "coverage.out")
	covReport := path.Join(covDir, "coverage.html")

	os.MkdirAll(covDir, 0755)

	args := []string{"test", "./...", "-cover", "-coverprofile", covFile}
	if verbose {
		args = append(args, "-v")
	}
	if err := sh.RunV("go", args...); err != nil {
		return err
	}

	if err := sh.RunV("go", "tool", "cover", "-html", covFile, "-o", covReport); err != nil {
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

// Prepare new release version. It updates version in version.go, prepares git tag and runs `release` task.
func ReleaseVersion(newVersion string, // New version to be released, e.g. 1.2.3 or MAJOR/MINOR/PATCH
) error {
	fmt.Printf("Preparing new '%s' release...\n", newVersion)

	newVersion, err := calcNewVersion(newVersion)
	if err != nil {
		return fmt.Errorf("failed to calculate a new version: %v", err)
	}

	fmt.Printf("Updating 'cmd/version.go' version...\n")
	content, err := os.ReadFile("cmd/root.go")
	if err != nil {
		return fmt.Errorf("failed to read version.go: %v", err)
	}
	reExpr := regexp.MustCompile(`var version = ".*"`)
	newContent := reExpr.ReplaceAllString(string(content), fmt.Sprintf(`var version = "%s"`, newVersion))
	err = os.WriteFile("cmd/version.go", []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write version.go: %v", err)
	}

	// Check git changes
	status, err := sh.Output("git", "status", "--porcelain", "cmd/version.go")
	if err != nil {
		return fmt.Errorf("failed to check git status: %v", err)
	}

	if status != "" {
		fmt.Printf("Committing changes...\n")

		if err := sh.RunV("git", "add", "cmd/version.go"); err != nil {
			return fmt.Errorf("failed to add file to git: %v", err)
		}
		if err := sh.RunV("git", "commit", "-m", "chore: bump version to "+newVersion); err != nil {
			return fmt.Errorf("failed to commit: %v", err)
		}

	}

	fmt.Printf("Removing existed '%s' git tag if any...\n", newVersion)
	sh.RunV("git", "tag", "-d", newVersion)

	fmt.Printf("Creating '%s' git tag...\n", newVersion)
	if err := sh.RunV("git", "tag", newVersion); err != nil {
		return fmt.Errorf("failed to create new tag: %v", err)
	}

	return Release()
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

func calcNewVersion(newVersion string) (string, error) {
	versionType := cleanupVersion(newVersion)

	if versionType != "major" && versionType != "minor" && versionType != "patch" {
		return newVersion, nil
	}

	gitVersion, err := getGitVersion()
	if err != nil {
		return "", fmt.Errorf("failed to get git version: %v", err)
	}
	latestTag := cleanupVersion(gitVersion)

	// Extract version components using regexp
	versionRegex := regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)`)
	matches := versionRegex.FindStringSubmatch(latestTag)
	if len(matches) != 4 {
		return "", fmt.Errorf("unexpected tag version format: %s", latestTag)
	}

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return "", fmt.Errorf("invalid major version: %v", err)
	}
	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return "", fmt.Errorf("invalid minor version: %v", err)
	}
	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		return "", fmt.Errorf("invalid patch version: %v", err)
	}

	switch versionType {
	case "major":
		major++
		minor = 0
		patch = 0
	case "minor":
		minor++
		patch = 0
	case "patch":
		patch++
	}

	newVersion = fmt.Sprintf("%d.%d.%d", major, minor, patch)
	fmt.Printf("New version is '%s'\n", newVersion)

	return newVersion, nil
}

func cleanupVersion(version string) string {
	version = strings.TrimSpace(version)
	version = strings.TrimPrefix(version, "v")
	return strings.ToLower(version)
}
