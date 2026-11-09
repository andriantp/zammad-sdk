package transport

import (
	"context"
	"fmt"
	"io"
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

func (c *Client) Do(ctx context.Context, method, path string, body, result any) error {
	req, err := c.newRequest(ctx, method, path, body)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute HTTP request: %w", err)
	}

	return decodeResponse(resp, result)
}

func (c *Client) DoRaw(ctx context.Context, method, path string) ([]byte, error) {
	req, err := c.newRequest(ctx, method, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, decodeAPIError(resp.StatusCode, data)
	}

	return data, nil
}
