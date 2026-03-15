package network

import "time"

// NewUIConnectorClient creates an HTTP client for UI connector manager operations.
func NewUIConnectorClient(timeout time.Duration) *Client {
	return NewClient(timeout)
}
