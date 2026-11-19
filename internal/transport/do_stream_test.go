package transport

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDoStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Token token=test-token" {
			t.Error("unexpected Authorization header")
		}

		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write([]byte("png-bytes"))
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL, Token: "test-token"})

	resp, err := client.DoStream(context.Background(), http.MethodGet, "/ticket_attachment/1/2/3")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	if string(data) != "png-bytes" || resp.Header.Get("Content-Type") != "image/png" {
		t.Errorf("unexpected response: %q %q", data, resp.Header.Get("Content-Type"))
	}
}

func TestDoStreamAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error":"The file does not belong to the specified article."}`))
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL, Token: "test-token"})

	resp, err := client.DoStream(context.Background(), http.MethodGet, "/ticket_attachment/1/2/999")
	if resp != nil {
		t.Error("expected nil response on API error")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 APIError, got %v", err)
	}
}

func TestDoStreamConnectionError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	serverURL := server.URL
	server.Close()

	client := NewClient(Config{BaseURL: serverURL, Token: "test-token"})

	_, err := client.DoStream(context.Background(), http.MethodGet, "/ticket_attachment/1/2/3")
	if err == nil || !strings.Contains(err.Error(), "execute HTTP request") {
		t.Errorf("expected wrapped connection error, got %v", err)
	}
}
