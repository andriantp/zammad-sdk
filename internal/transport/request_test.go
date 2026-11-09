
package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestNewRequest(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "https://zammad.example.com/",
		Token:   "test-token",
	})

	req, err := client.newRequest(
		context.Background(),
		http.MethodGet,
		"/tickets",
		nil,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if req.Method != http.MethodGet {
		t.Errorf("expected method GET, got %s", req.Method)
	}

	expectedURL := "https://zammad.example.com/api/v1/tickets"

	if req.URL.String() != expectedURL {
		t.Errorf("expected URL %q, got %q", expectedURL, req.URL.String())
	}

	if req.Header.Get("Authorization") != "Token token=test-token" {
		t.Error("unexpected Authorization header")
	}

	if req.Header.Get("Accept") != "application/json" {
		t.Error("unexpected Accept header")
	}

	if req.Header.Get("Content-Type") != "" {
		t.Error("expected empty Content-Type for request without body")
	}

	if req.Body != nil {
		t.Error("expected nil request body")
	}
}

func TestNewRequestWithJSONBody(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "https://zammad.example.com",
		Token:   "test-token",
	})

	body := map[string]any{
		"title": "Network Problem",
		"group": "Support",
	}

	req, err := client.newRequest(
		context.Background(),
		http.MethodPost,
		"/tickets",
		body,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if req.Header.Get("Content-Type") != "application/json" {
		t.Error("expected Content-Type application/json")
	}

	var result map[string]any

	if err := json.NewDecoder(req.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode request body: %v", err)
	}

	if result["title"] != "Network Problem" {
		t.Errorf("unexpected title: %v", result["title"])
	}

	if result["group"] != "Support" {
		t.Errorf("unexpected group: %v", result["group"])
	}
}

func TestNewRequestWithContext(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "https://zammad.example.com",
		Token:   "test-token",
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := client.newRequest(
		ctx,
		http.MethodGet,
		"/tickets",
		nil,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if req.Context() != ctx {
		t.Error("expected request context to match provided context")
	}
}

func TestNewRequestInvalidMethod(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "https://zammad.example.com",
		Token:   "test-token",
	})

	_, err := client.newRequest(
		context.Background(),
		"INVALID METHOD",
		"/tickets",
		nil,
	)

	if err == nil {
		t.Fatal("expected error for invalid HTTP method")
	}
}

func TestNewRequestEncodingError(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "https://zammad.example.com",
		Token:   "test-token",
	})

	body := map[string]any{
		"invalid": make(chan int),
	}

	_, err := client.newRequest(
		context.Background(),
		http.MethodPost,
		"/tickets",
		body,
	)

	if err == nil {
		t.Fatal("expected error for unsupported JSON value")
	}

	if !strings.Contains(err.Error(), "encode request body") {
		t.Errorf("unexpected error: %v", err)
	}
}
