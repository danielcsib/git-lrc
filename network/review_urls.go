package network

import "strings"

// ReviewNormalizedAPIURL builds a normalized LiveReview API URL for a relative API path.
func ReviewNormalizedAPIURL(baseURL, apiPath string) string {
	base := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	base = strings.TrimSuffix(base, "/api/v1")
	base = strings.TrimSuffix(base, "/api")
	return base + apiPath
}

// ReviewProxyRequestURL builds a proxy backend URL from API base, request path, and raw query.
func ReviewProxyRequestURL(apiBase, path, rawQuery string) string {
	backendURL := strings.TrimSuffix(apiBase, "/") + path
	if strings.TrimSpace(rawQuery) != "" {
		backendURL += "?" + rawQuery
	}
	return backendURL
}
