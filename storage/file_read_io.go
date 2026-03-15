package storage

import (
	"fmt"
	"os"
)

// ReadPendingUpdateStateBytes reads the persisted pending-update state payload.
func ReadPendingUpdateStateBytes(statePath string) ([]byte, error) {
	data, err := os.ReadFile(statePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read pending update state file %s: %w", statePath, err)
	}
	return data, nil
}

// ReadUpdateLockMetadataBytes reads persisted update-lock metadata bytes.
func ReadUpdateLockMetadataBytes(lockPath string) ([]byte, error) {
	data, err := os.ReadFile(lockPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read update lock metadata file %s: %w", lockPath, err)
	}
	return data, nil
}

// OpenFileForRead opens a file for read-only access through the storage boundary.
func OpenFileForRead(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for read %s: %w", path, err)
	}
	return file, nil
}
