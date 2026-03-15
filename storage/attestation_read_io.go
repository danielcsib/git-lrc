package storage

import (
	"fmt"
	"os"
)

// ReadAttestationFile reads an attestation JSON payload from disk.
func ReadAttestationFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read attestation file %s: %w", path, err)
	}
	return data, nil
}

// EnsureAttestationOutputDir creates the .git/lrc/attestations directory.
func EnsureAttestationOutputDir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create attestation output directory %s: %w", path, err)
	}
	return nil
}
