package zammad

import (
	"net/http"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient(
		WithBaseURL("https://zammad.example.com"),
		WithToken("test-token"),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if client == nil {
		t.Fatal("expected client, got nil")
	}

	if client.config == nil {
		t.Fatal("expected config, got nil")
	}

	if client.config.baseURL != "https://zammad.example.com" {
		t.Errorf("unexpected baseURL: %q", client.config.baseURL)
	}

	if client.config.token != "test-token" {
		t.Errorf("unexpected token: %q", client.config.token)
	}

	if client.config.timeout != defaultTimeout {
		t.Errorf("unexpected timeout: %v", client.config.timeout)
	}
}

func TestNewClientMissingBaseURL(t *testing.T) {
	_, err := NewClient(
		WithToken("test-token"),
	)

	if err == nil {
		t.Fatal("expected error for missing base URL")
	}
}

func TestNewClientInvalidBaseURL(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
	}{
		{
			name:    "invalid URL",
			baseURL: "://invalid-url",
		},
		{
			name:    "missing host",
			baseURL: "https://",
		},
		{
			name:    "unsupported scheme",
			baseURL: "ftp://zammad.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(
				WithBaseURL(tt.baseURL),
				WithToken("test-token"),
			)

			if err == nil {
				t.Fatal("expected error for invalid base URL")
			}
		})
	}
}

func TestNewClientMissingToken(t *testing.T) {
	_, err := NewClient(
		WithBaseURL("https://zammad.example.com"),
	)

	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewClientWithNilOption(t *testing.T) {
	client, err := NewClient(
		nil,
		WithBaseURL("https://zammad.example.com"),
		WithToken("test-token"),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if client == nil {
		t.Fatal("expected client, got nil")
	}
}

func TestNewClientWithCustomHTTPClient(t *testing.T) {
	customClient := &http.Client{
		Timeout: 15 * time.Second,
	}

	client, err := NewClient(
		WithBaseURL("https://zammad.example.com"),
		WithToken("test-token"),
		WithHTTPClient(customClient),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if client.config.httpClient != customClient {
		t.Error("expected custom HTTP client to be preserved")
	}
}

func TestNewClientWithCustomTimeout(t *testing.T) {
	client, err := NewClient(
		WithBaseURL("https://zammad.example.com"),
		WithToken("test-token"),
		WithTimeout(10*time.Second),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if client.config.timeout != 10*time.Second {
		t.Errorf(
			"expected timeout %v, got %v",
			10*time.Second,
			client.config.timeout,
		)
	}
}

func TestNewClientInitializesServices(t *testing.T) {
	client, err := NewClient(
		WithBaseURL("https://zammad.example.com"),
		WithToken("test-token"),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if client.Tickets == nil {
		t.Error("expected Tickets service to be initialized")
	}

	if client.Articles == nil {
		t.Error("expected Articles service to be initialized")
	}

	if client.Attachments == nil {
		t.Error("expected Attachments service to be initialized")
	}

	if client.Users == nil {
		t.Error("expected Users service to be initialized")
	}

	if client.Groups == nil {
		t.Error("expected Groups service to be initialized")
	}

	if client.Organizations == nil {
		t.Error("expected Organizations service to be initialized")
	}
}

func TestNewClientInitializesTransport(t *testing.T) {
	client, err := NewClient(
		WithBaseURL("https://zammad.example.com"),
		WithToken("test-token"),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if client.transport == nil {
		t.Fatal("expected transport client to be initialized")
	}
}

func TestNewClientTransportConfiguration(t *testing.T) {
	customHTTPClient := &http.Client{
		Timeout: 15 * time.Second,
	}

	client, err := NewClient(
		WithBaseURL("https://zammad.example.com"),
		WithToken("test-token"),
		WithHTTPClient(customHTTPClient),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if client.transport == nil {
		t.Fatal("expected transport client to be initialized")
	}
}
