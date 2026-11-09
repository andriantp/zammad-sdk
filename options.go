package zammad

import (
	"net/http"
	"strings"
	"time"
)

type Option func(*config)

func WithBaseURL(url string) Option {
	return func(c *config) {
		c.baseURL = strings.TrimRight(url, "/")
	}
}

func WithToken(token string) Option {
	return func(c *config) {
		c.token = token
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(c *config) {
		if client != nil {
			c.httpClient = client
		}
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *config) {
		if timeout > 0 {
			c.timeout = timeout
		}
	}
}
