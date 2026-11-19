package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// readResponse reads and closes the response body. A non-2xx status is
// returned as an *APIError.
func readResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusMultipleChoices {
		return nil, decodeAPIError(resp.StatusCode, body)
	}

	return body, nil
}

// decodeResponse reads the response like readResponse and decodes a JSON
// body into result. An empty body or nil result is not an error.
func decodeResponse(resp *http.Response, result any) error {
	body, err := readResponse(resp)
	if err != nil {
		return err
	}

	if len(body) == 0 || result == nil {
		return nil
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("decode response body: %w", err)
	}

	return nil
}

func decodeAPIError(statusCode int, body []byte) error {
	apiErr := &APIError{
		StatusCode: statusCode,
		Message:    http.StatusText(statusCode),
	}

	if len(body) == 0 {
		return apiErr
	}

	var payload struct {
		Error   string `json:"error"`
		Message string `json:"message"`
		Detail  string `json:"detail"`
	}

	if err := json.Unmarshal(body, &payload); err == nil {
		if payload.Error != "" {
			apiErr.Message = payload.Error
		}

		if payload.Message != "" {
			apiErr.Message = payload.Message
		}

		apiErr.Detail = payload.Detail

		return apiErr
	}

	apiErr.Detail = string(body)

	return apiErr
}
