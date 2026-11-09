package transport

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func newTestResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func TestDecodeResponse(t *testing.T) {
	resp := newTestResponse(
		http.StatusOK,
		`{"id": 123, "title": "Network Problem"}`,
	)

	var result struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	}

	err := decodeResponse(resp, &result)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != 123 {
		t.Errorf("expected ID 123, got %d", result.ID)
	}

	if result.Title != "Network Problem" {
		t.Errorf("unexpected title: %s", result.Title)
	}
}

func TestDecodeResponseNoContent(t *testing.T) {
	resp := newTestResponse(http.StatusNoContent, "")

	err := decodeResponse(resp, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDecodeResponseEmptyBody(t *testing.T) {
	resp := newTestResponse(http.StatusOK, "")

	var result map[string]any

	err := decodeResponse(resp, &result)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != nil {
		t.Error("expected nil result")
	}
}

func TestDecodeResponseNilResult(t *testing.T) {
	resp := newTestResponse(
		http.StatusOK,
		`{"id": 123}`,
	)

	err := decodeResponse(resp, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDecodeResponseAPIError(t *testing.T) {
	resp := newTestResponse(
		http.StatusNotFound,
		`{"error": "Not Found", "message": "Ticket not found"}`,
	)

	err := decodeResponse(resp, nil)
	if err == nil {
		t.Fatal("expected API error")
	}

	var apiErr *APIError

	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}

	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusNotFound,
			apiErr.StatusCode,
		)
	}

	if apiErr.Message != "Ticket not found" {
		t.Errorf("unexpected message: %s", apiErr.Message)
	}
}

func TestDecodeResponseAPIErrorNonJSON(t *testing.T) {
	resp := newTestResponse(
		http.StatusInternalServerError,
		"Internal Server Error",
	)

	err := decodeResponse(resp, nil)
	if err == nil {
		t.Fatal("expected API error")
	}

	var apiErr *APIError

	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}

	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("unexpected status code: %d", apiErr.StatusCode)
	}

	if apiErr.Detail != "Internal Server Error" {
		t.Errorf("unexpected detail: %s", apiErr.Detail)
	}
}

func TestDecodeResponseInvalidJSON(t *testing.T) {
	resp := newTestResponse(
		http.StatusOK,
		`{"id": 123`,
	)

	var result map[string]any

	err := decodeResponse(resp, &result)
	if err == nil {
		t.Fatal("expected JSON decoding error")
	}

	if !strings.Contains(err.Error(), "decode response body") {
		t.Errorf("unexpected error: %v", err)
	}
}
