package network

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func validateSelfUpdateURL(fullURL string) error {
	parsed, err := url.Parse(fullURL)
	if err != nil {
		return fmt.Errorf("invalid self-update URL %q: %w", fullURL, err)
	}
	if parsed.Scheme != "https" {
		return fmt.Errorf("self-update URL must use https: %s", fullURL)
	}

	manifestURL, err := url.Parse(SelfUpdateManifestURL())
	if err != nil {
		return fmt.Errorf("invalid manifest URL configuration: %w", err)
	}
	if !strings.EqualFold(parsed.Host, manifestURL.Host) {
		return fmt.Errorf("self-update URL host %q does not match trusted host %q", parsed.Host, manifestURL.Host)
	}

	return nil
}

// SelfUpdateFetchManifest fetches the top-level self-update manifest.
func SelfUpdateFetchManifest(client *Client) (*Response, error) {
	return client.DoJSON(http.MethodGet, SelfUpdateManifestURL(), nil, "", "", nil)
}

// SelfUpdateFetchReleaseManifest fetches the release manifest at a provided URL.
func SelfUpdateFetchReleaseManifest(client *Client, fullURL string) (*Response, error) {
	if err := validateSelfUpdateURL(fullURL); err != nil {
		return nil, err
	}
	return client.DoJSON(http.MethodGet, fullURL, nil, "", "", nil)
}

// SelfUpdateDownloadBinaryTo streams a self-update binary into dst.
func SelfUpdateDownloadBinaryTo(client *Client, fullURL string, dst io.Writer) (int, error) {
	if err := validateSelfUpdateURL(fullURL); err != nil {
		return 0, err
	}

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create self-update request: %w", err)
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if _, err := io.Copy(dst, resp.Body); err != nil {
		return resp.StatusCode, fmt.Errorf("failed to stream self-update body: %w", err)
	}

	return resp.StatusCode, nil
}
