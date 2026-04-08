package tools

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	BackupStrategyNone    = "none"
	BackupStrategyOnPlace = "in_place"
	originalFileTemplate  = "*.*_original"
)

type Backuper struct {
	backupDir   string
	dryRun      bool
	recursively bool
	backupName  string
}

func NewBackuper(backupLocation string, dryRun bool, recursively bool, taskName string) *Backuper {
	result := Backuper{backupDir: backupLocation, dryRun: dryRun, recursively: recursively}
	result.initLocation(taskName)
	return &result
}

func (b *Backuper) initLocation(taskName string) {
	if b.ShouldPerformBackup() {
		taskName = strings.TrimPrefix(taskName, "media-tool ")

		taskName = strings.ReplaceAll(taskName, " ", "_")

		t := time.Now()
		ts := t.Format("20060102150405")

		b.backupName = taskName + "_" + ts
		b.backupDir = b.expandPath(path.Join(b.backupDir, b.backupName))
	}
}

func (b *Backuper) GetBackupDir() string {
	return b.backupDir
}

func (b *Backuper) ShouldExifDeleteOriginal() bool {
	return b.backupDir == "" || strings.ToLower(b.backupDir) == BackupStrategyNone
}

func (b *Backuper) ShouldPerformBackup() bool {
	return !(b.dryRun || b.ShouldExifDeleteOriginal() || strings.ToLower(b.backupDir) == BackupStrategyOnPlace)
}

func (b *Backuper) getWorkDir(workDirectory string) (string, error) {
	if stat, err := os.Stat(workDirectory); err == nil && stat.IsDir() {
		return workDirectory, nil
	} else {
		dir := filepath.Dir(workDirectory)
		if stat, err := os.Stat(dir); err == nil && stat.IsDir() {
			return dir, nil
		}
	}
	return "", fmt.Errorf("Unable to use '%s' working directory", workDirectory)
}

func (b *Backuper) CleanupWorkDir(workDirectory string) error {
	if b.dryRun {
		log.Debugf("Nothing to backup (dry run)")
		return nil
	} else if !b.ShouldPerformBackup() {
		log.Debugf("Nothing to backup (Backup strategy is '%s')", b.backupDir)
		return nil
	}

	workDir, err := b.getWorkDir(workDirectory)
	if err != nil {
		return err
	}

	filesToMove, err := b.findOriginals(workDir)
	if err != nil {
		log.Errorf("Unable to find origin files to backup: %s", err)
		return err
	}

	filesTotal := len(filesToMove)
	if filesTotal > 0 {
		log.Infof("Moving %d '%s' files from '%s' to '%s'", filesTotal, originalFileTemplate, workDir, b.backupDir)

		for i, fileToMove := range filesToMove {
			log.Infof("Moving [%d/%d]: '%s'...", i+1, filesTotal, fileToMove)

			oldFile := path.Join(workDir, fileToMove)
			newFile := path.Join(b.backupDir, fileToMove)

			newDir := filepath.Dir(newFile)
			if err := os.MkdirAll(newDir, os.ModePerm); err != nil {
				log.Errorf("Unable to prepare backup storage: %s", err)
				return err
			}
			err := os.Rename(oldFile, newFile)
			if err != nil {
				log.Errorf("Unable to move origin files to backup storage: %s", err)
				return err
			}
		}
	}
	return nil
}

func (b *Backuper) findOriginals(workDir string) ([]string, error) {
	if b.recursively {
		//Because **/*.blah is not supported
		return b.recursiveGlob(workDir, originalFileTemplate)
	}
	workFSys := os.DirFS(workDir)
	return fs.Glob(workFSys, originalFileTemplate)
}

func (b *Backuper) recursiveGlob(root, pattern string) ([]string, error) {
	var matches []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip problematic directories
		}

		if !d.IsDir() {
			if matched, _ := filepath.Match(pattern, d.Name()); matched {
				relPath, err := filepath.Rel(root, path)
				if err != nil {
					return nil // Skip problematic files
				}
				matches = append(matches, relPath)
			}
		}
		return nil
	})

	return matches, err
}

func (c *Backuper) expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, strings.TrimPrefix(path, "~"))
		}
	}

	return filepath.Clean(os.ExpandEnv(path))
}
