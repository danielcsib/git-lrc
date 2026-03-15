package storage

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
)

// EnsureManagedHooksDir creates the directory that stores lrc-managed hook scripts.
func EnsureManagedHooksDir(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create managed hooks directory %s: %w", dir, err)
	}
	return nil
}

// EnsureHooksPathDir creates the resolved hooksPath root when it does not exist.
func EnsureHooksPathDir(hooksPath string) error {
	if err := os.MkdirAll(hooksPath, 0755); err != nil {
		return fmt.Errorf("failed to create hooks path %s: %w", hooksPath, err)
	}
	return nil
}

// EnsureHooksBackupDir creates the hooks backup directory.
func EnsureHooksBackupDir(backupDir string) error {
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create hooks backup directory %s: %w", backupDir, err)
	}
	return nil
}

// EnsureRepoLRCStateDir creates the repo-local .git/lrc directory used by hook state.
func EnsureRepoLRCStateDir(lrcDir string) error {
	if err := os.MkdirAll(lrcDir, 0755); err != nil {
		return fmt.Errorf("failed to create repository lrc state directory %s: %w", lrcDir, err)
	}
	return nil
}

// ReadHookFile reads hook script content.
func ReadHookFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read hook file %s: %w", path, err)
	}
	return data, nil
}

// ReadHookMetaFile reads hook metadata JSON content.
func ReadHookMetaFile(hooksPath, metaFilename string) ([]byte, error) {
	metaPath := filepath.Join(hooksPath, metaFilename)
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read hook metadata %s: %w", metaPath, err)
	}
	return data, nil
}

// RemoveHookMetaFile removes hook metadata JSON file.
func RemoveHookMetaFile(hooksPath, metaFilename string) error {
	metaPath := filepath.Join(hooksPath, metaFilename)
	if err := Remove(metaPath); err != nil {
		return fmt.Errorf("failed to remove hook metadata %s: %w", metaPath, err)
	}
	return nil
}

// RemoveDirIfEmpty removes a directory only when it has no entries.
func RemoveDirIfEmpty(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("failed to read directory %s: %w", dir, err)
	}
	if len(entries) == 0 {
		if err := Remove(dir); err != nil {
			if errors.Is(err, fs.ErrNotExist) || errors.Is(err, syscall.ENOTEMPTY) || errors.Is(err, syscall.EEXIST) {
				return nil
			}
			return fmt.Errorf("failed to remove empty directory %s: %w", dir, err)
		}
	}
	return nil
}
