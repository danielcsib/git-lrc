package network

import "time"

// NewSetupClient creates an HTTP client for setup-provisioning network operations.
func NewSetupClient(timeout time.Duration) *Client {
	return NewClient(timeout)
}
