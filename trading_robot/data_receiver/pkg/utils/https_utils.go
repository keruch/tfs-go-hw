package utils

import (
	"net"
	"net/http"
)

type RetryableClient struct {
	*http.Client
}

func NewRetryableClient(client *http.Client) *RetryableClient {
	return &RetryableClient{
		client,
	}
}

// SendRequestWithoutTimeout sends request that ignores timeout errors
func (c *RetryableClient) SendRequestWithoutTimeout(req *http.Request) error {
	_, err := c.Do(req)
	if err != nil {
		// ignore timeout errors
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return nil
		}
		return err
	}
	return nil
}
