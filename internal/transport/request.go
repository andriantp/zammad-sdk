package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const apiVersion = "/api/v1"

func (c *Client) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	endpoint := strings.TrimRight(c.baseURL, "/") +
		apiVersion +
		"/" + strings.TrimLeft(path, "/")

	var requestBody io.Reader
	if body != nil {
		buffer := new(bytes.Buffer)
		if err := json.NewEncoder(buffer).Encode(body); err != nil {
			return nil, fmt.Errorf("encode request body: %w", err)
		}
		requestBody = buffer
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, requestBody)
	if err != nil {
		return nil, fmt.Errorf("create HTTP request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Token token="+c.token)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}
