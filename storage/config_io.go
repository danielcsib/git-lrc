package storage

import (
	"fmt"
	"os"
)

// ReadConfigFile reads a persisted lrc config file.
func ReadConfigFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}
	return data, nil
}
