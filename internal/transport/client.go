package transport

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
	Timeout    time.Duration
}

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewClient(cfg Config) *Client {
	httpClient := cfg.HTTPClient

	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: cfg.Timeout,
		}
	}

	return &Client{
		baseURL:    cfg.BaseURL,
		token:      cfg.Token,
		httpClient: httpClient,
	}
}

// Do sends a request with an optional JSON body and decodes a JSON response
// into result (which may be nil).
func (c *Client) Do(ctx context.Context, method, path string, body, result any) error {
	resp, err := c.execute(ctx, method, path, body)
	if err != nil {
		return err
	}

	return decodeResponse(resp, result)
}

// DoRaw sends a request without a body and returns the raw response body,
// for binary content such as attachments.
func (c *Client) DoRaw(ctx context.Context, method, path string) ([]byte, error) {
	resp, err := c.execute(ctx, method, path, nil)
	if err != nil {
		return nil, err
	}

	return readResponse(resp)
}

// DoStream sends a request without a body and returns the response with its
// body unread, for streaming large downloads. The caller must close the
// body. A non-2xx response is read, closed and returned as an *APIError.
func (c *Client) DoStream(ctx context.Context, method, path string) (*http.Response, error) {
	resp, err := c.execute(ctx, method, path, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		_, err := readResponse(resp)
		return nil, err
	}

	return resp, nil
}

func (c *Client) execute(ctx context.Context, method, path string, body any) (*http.Response, error) {
	req, err := c.newRequest(ctx, method, path, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute HTTP request: %w", err)
	}

	return resp, nil
}
