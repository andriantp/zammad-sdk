package transport

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestDoRaw(t *testing.T) {
	// Bytes that are not valid UTF-8 or JSON must come back untouched.
	want := []byte{0x89, 'P', 'N', 'G', 0x00, 0xff, 0xfe}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/ticket_attachment/1/2/3" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if r.Header.Get("Authorization") != "Token token=test-token" {
			t.Error("unexpected Authorization header")
		}

		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(want)
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL, Token: "test-token"})

	got, err := client.DoRaw(context.Background(), http.MethodGet, "/ticket_attachment/1/2/3")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !bytes.Equal(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestDoRawAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error":"The file does not belong to the specified article."}`))
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL, Token: "test-token"})

	data, err := client.DoRaw(context.Background(), http.MethodGet, "/ticket_attachment/1/2/999")
	if data != nil {
		t.Errorf("expected nil data on error, got %q", data)
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T (%v)", err, err)
	}

	if apiErr.StatusCode != http.StatusForbidden ||
		apiErr.Message != "The file does not belong to the specified article." {
		t.Errorf("unexpected API error: %+v", apiErr)
	}
}

// DoRaw must report connection failures the same way as Do.
func TestDoRawConnectionError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	serverURL := server.URL
	server.Close()

	client := NewClient(Config{BaseURL: serverURL, Token: "test-token"})

	_, err := client.DoRaw(context.Background(), http.MethodGet, "/ticket_attachment/1/2/3")
	if err == nil {
		t.Fatal("expected connection error")
	}

	if !strings.Contains(err.Error(), "execute HTTP request") {
		t.Errorf("expected wrapped error, got %v", err)
	}

	var urlErr *url.Error
	if !errors.As(err, &urlErr) {
		t.Errorf("expected underlying *url.Error to stay reachable, got %T", err)
	}
}

func TestDoRawContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL, Token: "test-token"})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.DoRaw(ctx, http.MethodGet, "/ticket_attachment/1/2/3")
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}
