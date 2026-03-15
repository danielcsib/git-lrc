package appui

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/HexmosTech/git-lrc/network"
)

func (s *connectorManagerServer) proxyJSONRequest(method, apiPath string, payload []byte) (int, []byte, error) {
	url := network.ReviewNormalizedAPIURL(s.cfg.APIURL, apiPath)

	s.mu.Lock()
	jwt := s.cfg.JWT
	orgID := s.cfg.OrgID
	s.mu.Unlock()

	if strings.TrimSpace(jwt) == "" || strings.TrimSpace(orgID) == "" {
		return http.StatusUnauthorized, []byte(`{"error":"not authenticated. Open Home and use Re-authenticate."}`), nil
	}

	status, respBody, err := s.forwardJSONRequest(method, url, payload, jwt, orgID)
	if err != nil {
		return status, nil, err
	}

	if status == http.StatusUnauthorized {
		refreshed, refreshErr := s.refreshAccessToken(jwt)
		if refreshErr != nil {
			log.Printf("failed to refresh lrc ui token: %v", refreshErr)
			return status, respBody, nil
		}
		if refreshed {
			s.mu.Lock()
			newJWT := s.cfg.JWT
			s.mu.Unlock()

			status, retryBody, retryErr := s.forwardJSONRequest(method, url, payload, newJWT, orgID)
			if retryErr != nil {
				return status, nil, retryErr
			}
			return status, retryBody, nil
		}
	}

	return status, respBody, nil
}

func (s *connectorManagerServer) forwardJSONRequest(method, url string, payload []byte, jwt string, orgID string) (int, []byte, error) {
	resp, err := network.ReviewForwardJSONWithBearer(s.client, method, url, payload, jwt, orgID)
	if err != nil {
		return http.StatusBadGateway, nil, fmt.Errorf("failed to call LiveReview API: %w", err)
	}

	return resp.StatusCode, resp.Body, nil
}

func (s *connectorManagerServer) refreshAccessToken(failedJWT string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if strings.TrimSpace(s.cfg.JWT) != strings.TrimSpace(failedJWT) {
		return true, nil
	}

	if strings.TrimSpace(s.cfg.RefreshJWT) == "" {
		return false, fmt.Errorf("refresh_token missing in %s", s.cfg.ConfigPath)
	}

	refreshURL := network.ReviewNormalizedAPIURL(s.cfg.APIURL, "/api/v1/auth/refresh")
	reqBody, err := json.Marshal(authRefreshRequest{RefreshToken: s.cfg.RefreshJWT})
	if err != nil {
		return false, fmt.Errorf("failed to marshal refresh request: %w", err)
	}

	resp, err := s.client.Do(http.MethodPost, refreshURL, reqBody, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return false, fmt.Errorf("refresh request failed: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, fmt.Errorf("refresh failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(resp.Body)))
	}

	var refreshResp authRefreshResponse
	if err := json.Unmarshal(resp.Body, &refreshResp); err != nil {
		return false, fmt.Errorf("failed to decode refresh response: %w", err)
	}

	if strings.TrimSpace(refreshResp.AccessToken) == "" {
		return false, fmt.Errorf("refresh response missing access token")
	}

	s.cfg.JWT = strings.TrimSpace(refreshResp.AccessToken)
	if strings.TrimSpace(refreshResp.RefreshToken) != "" {
		s.cfg.RefreshJWT = strings.TrimSpace(refreshResp.RefreshToken)
	}

	if err := persistAuthTokensToConfig(s.cfg.ConfigPath, s.cfg.JWT, s.cfg.RefreshJWT); err != nil {
		log.Printf("warning: refreshed token obtained but failed to update %s: %v", s.cfg.ConfigPath, err)
	}

	return true, nil
}
