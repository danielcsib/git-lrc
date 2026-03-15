package appui

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"github.com/HexmosTech/git-lrc/storage"
)

func persistConnectorsToConfig(configPath string, connectors []aiConnectorRemote) error {
	originalBytes, err := storage.ReadConfigFile(configPath)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("failed to read config for connector snapshot: %w", err)
		}
		originalBytes = []byte{}
	}

	originalContent := string(originalBytes)
	cleanedContent := stripManagedAIConnectorsSection(originalContent)
	managedSection := renderManagedAIConnectorsSection(connectors)

	trimmed := strings.TrimRight(cleanedContent, "\n\r\t ")
	var updatedContent string
	if trimmed == "" {
		updatedContent = managedSection + "\n"
	} else {
		updatedContent = trimmed + "\n\n" + managedSection + "\n"
	}

	if err := storage.WriteFileAtomically(configPath, []byte(updatedContent), 0600); err != nil {
		return fmt.Errorf("failed to replace config file: %w", err)
	}

	return nil
}

func persistAuthTokensToConfig(configPath string, jwt string, refreshToken string) error {
	originalBytes, err := storage.ReadConfigFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config for token update: %w", err)
	}

	content := string(originalBytes)
	updated := upsertQuotedConfigValue(content, "jwt", jwt)
	if strings.TrimSpace(refreshToken) != "" {
		updated = upsertQuotedConfigValue(updated, "refresh_token", refreshToken)
	}

	if err := storage.WriteFileAtomically(configPath, []byte(updated), 0600); err != nil {
		return fmt.Errorf("failed to replace config file: %w", err)
	}

	return nil
}
