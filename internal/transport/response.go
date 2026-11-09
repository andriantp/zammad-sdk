package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func decodeResponse(resp *http.Response, result any) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusMultipleChoices {
		return decodeAPIError(resp.StatusCode, body)
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
