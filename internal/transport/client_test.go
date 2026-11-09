
package transport

import (
	"net/http"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	cfg := Config{
		BaseURL: "https://zammad.example.com",
		Token:   "test-token",
		Timeout: 30 * time.Second,
	}

	client := NewClient(cfg)

	if client == nil {
		t.Fatal("expected client, got nil")
	}

	if client.baseURL != cfg.BaseURL {
		t.Errorf("expected baseURL %q, got %q", cfg.BaseURL, client.baseURL)
	}

	if client.token != cfg.Token {
		t.Errorf("expected token %q, got %q", cfg.Token, client.token)
	}

	if client.httpClient == nil {
		t.Fatal("expected HTTP client, got nil")
	}

	if client.httpClient.Timeout != cfg.Timeout {
		t.Errorf(
			"expected timeout %v, got %v",
			cfg.Timeout,
			client.httpClient.Timeout,
		)
	}
}

func TestNewClientWithCustomHTTPClient(t *testing.T) {
	customClient := &http.Client{
		Timeout: 15 * time.Second,
	}

	cfg := Config{
		BaseURL:    "https://zammad.example.com",
		Token:      "test-token",
		HTTPClient: customClient,
		Timeout:    30 * time.Second,
	}

	client := NewClient(cfg)

	if client.httpClient != customClient {
		t.Error("expected custom HTTP client to be preserved")
	}

	if customClient.Timeout != 15*time.Second {
		t.Errorf(
			"expected custom timeout to remain 15s, got %v",
			customClient.Timeout,
		)
	}
}

func TestNewClientWithoutCustomHTTPClient(t *testing.T) {
	cfg := Config{
		Timeout: 10 * time.Second,
	}

	client := NewClient(cfg)

	if client.httpClient == nil {
		t.Fatal("expected default HTTP client, got nil")
	}

	if client.httpClient.Timeout != 10*time.Second {
		t.Errorf(
			"expected timeout 10s, got %v",
			client.httpClient.Timeout,
		)
	}
}
