package storage

import (
	"fmt"
	"os"
)

// ReadInitialMessageFile reads the pre-seeded review initial message file.
func ReadInitialMessageFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read initial message file %s: %w", path, err)
	}
	return data, nil
}

// ReadDiffFile reads a diff payload from disk.
func ReadDiffFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read diff file %s: %w", path, err)
	}
	return data, nil
}

// CreateTempReviewHTMLFile creates a temporary HTML output file for served reviews.
// The caller owns cleanup and must close and remove the file when done.
func CreateTempReviewHTMLFile() (*os.File, error) {
	file, err := os.CreateTemp("", "lrc-review-*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary review HTML file: %w", err)
	}
	return file, nil
}

// ReadEditorBackupFile reads the saved editor backup value for hook management.
func ReadEditorBackupFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read editor backup file %s: %w", path, err)
	}
	return data, nil
}
