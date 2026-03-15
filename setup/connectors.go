package setup

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/HexmosTech/git-lrc/network"
)

// ValidateGeminiKey checks the key against LiveReview's validate-key endpoint.
func ValidateGeminiKey(result *SetupResult, geminiKey string) (bool, string, error) {
	reqBody := ValidateKeyRequest{
		Provider: "gemini",
		APIKey:   geminiKey,
		Model:    DefaultGeminiModel,
	}

	client := network.NewSetupClient(30 * time.Second)
	resp, err := network.SetupValidateConnectorKey(client, CloudAPIURL, result.OrgID, reqBody, result.AccessToken)
	if err != nil {
		return false, "", fmt.Errorf("failed to validate key: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return false, "", fmt.Errorf("validate-key returned %d: %s", resp.StatusCode, string(resp.Body))
	}

	var valResp ValidateKeyResponse
	if err := json.Unmarshal(resp.Body, &valResp); err != nil {
		return false, "", fmt.Errorf("failed to parse validation response: %w", err)
	}

	return valResp.Valid, valResp.Message, nil
}

// CreateGeminiConnector creates a Gemini AI connector in LiveReview.
func CreateGeminiConnector(result *SetupResult, geminiKey string) error {
	reqBody := CreateConnectorRequest{
		ProviderName:  "gemini",
		APIKey:        geminiKey,
		ConnectorName: "Gemini Flash",
		SelectedModel: DefaultGeminiModel,
		DisplayOrder:  0,
	}

	client := network.NewSetupClient(30 * time.Second)
	resp, err := network.SetupCreateConnector(client, CloudAPIURL, result.OrgID, reqBody, result.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to create connector: %w", err)
	}
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("create connector returned %d: %s", resp.StatusCode, string(resp.Body))
	}

	return nil
}
