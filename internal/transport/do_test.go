
package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestDoGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}

			if r.URL.Path != "/api/v1/tickets/123" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}

			if r.Header.Get("Authorization") != "Token token=test-token" {
				t.Error("unexpected Authorization header")
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":    123,
				"title": "Network Problem",
			})
		},
	))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		Token:   "test-token",
	})

	var result struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	}

	err := client.Do(
		context.Background(),
		http.MethodGet,
		"/tickets/123",
		nil,
		&result,
	)
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

func TestDoPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}

			if r.Header.Get("Content-Type") != "application/json" {
				t.Error("unexpected Content-Type")
			}

			var payload map[string]string

			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Errorf("failed to decode request: %v", err)
				return
			}

			if payload["title"] != "Network Problem" {
				t.Errorf("unexpected title: %s", payload["title"])
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":    456,
				"title": payload["title"],
			})
		},
	))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		Token:   "test-token",
	})

	body := map[string]string{
		"title": "Network Problem",
	}

	var result struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	}

	err := client.Do(
		context.Background(),
		http.MethodPost,
		"/tickets",
		body,
		&result,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != 456 {
		t.Errorf("expected ID 456, got %d", result.ID)
	}
}

func TestDoAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)

			_, _ = w.Write([]byte(
				`{"error":"Not Found","message":"Ticket not found"}`,
			))
		},
	))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		Token:   "test-token",
	})

	err := client.Do(
		context.Background(),
		http.MethodGet,
		"/tickets/999",
		nil,
		nil,
	)

	if err == nil {
		t.Fatal("expected API error")
	}

	var apiErr *APIError

	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}

	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("unexpected status code: %d", apiErr.StatusCode)
	}
}

func TestDoConnectionError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {},
	))

	serverURL := server.URL
	server.Close()

	client := NewClient(Config{
		BaseURL: serverURL,
		Token:   "test-token",
	})

	err := client.Do(
		context.Background(),
		http.MethodGet,
		"/tickets",
		nil,
		nil,
	)

	if err == nil {
		t.Fatal("expected connection error")
	}

	if !strings.Contains(err.Error(), "execute HTTP request") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDoContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			<-r.Context().Done()
		},
	))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		Token:   "test-token",
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := client.Do(
		ctx,
		http.MethodGet,
		"/tickets",
		nil,
		nil,
	)

	if err == nil {
		t.Fatal("expected context cancellation error")
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestDoWithTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			<-r.Context().Done()
		},
	))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		Token:   "test-token",
		Timeout: 50 * time.Millisecond,
	})

	err := client.Do(
		context.Background(),
		http.MethodGet,
		"/tickets",
		nil,
		nil,
	)

	if err == nil {
		t.Fatal("expected timeout error")
	}
}
