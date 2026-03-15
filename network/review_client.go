package network

import "time"

// NewReviewAPIClient creates an HTTP client for review submit/poll operations.
func NewReviewAPIClient(timeout time.Duration) *Client {
	return NewClient(timeout)
}

// NewReviewProxyClient creates an HTTP client for review-runtime proxy operations.
func NewReviewProxyClient(timeout time.Duration) *Client {
	return NewClient(timeout)
}
