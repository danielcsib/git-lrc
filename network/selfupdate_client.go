package network

import "time"

// NewSelfUpdateClient creates an HTTP client for self-update network operations.
func NewSelfUpdateClient(timeout time.Duration) *Client {
	return NewClient(timeout)
}
