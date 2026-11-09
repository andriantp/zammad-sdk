package zammad

import (
	"net/http"
	"testing"
	"time"
)

func TestWithBaseURL(t *testing.T) {
	cfg := newConfig()

	WithBaseURL("https://zammad.example.com/")(cfg)

	expected := "https://zammad.example.com"

	if cfg.baseURL != expected {
		t.Errorf("expected %q, got %q", expected, cfg.baseURL)
	}
}

func TestWithToken(t *testing.T) {
	cfg := newConfig()

	WithToken("test-token")(cfg)

	if cfg.token != "test-token" {
		t.Errorf("expected test-token, got %q", cfg.token)
	}
}

func TestWithHTTPClient(t *testing.T) {
	cfg := newConfig()
	customClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	WithHTTPClient(customClient)(cfg)

	if cfg.httpClient != customClient {
		t.Error("expected custom HTTP client to be set")
	}
}

func TestWithHTTPClientNil(t *testing.T) {
	cfg := newConfig()

	WithHTTPClient(nil)(cfg)

	if cfg.httpClient != nil {
		t.Error("expected HTTP client to remain nil")
	}
}

func TestWithTimeout(t *testing.T) {
	cfg := newConfig()

	WithTimeout(10 * time.Second)(cfg)

	expected := 10 * time.Second

	if cfg.timeout != expected {
		t.Errorf("expected timeout %v, got %v", expected, cfg.timeout)
	}
}

func TestWithTimeoutInvalid(t *testing.T) {
	cfg := newConfig()

	WithTimeout(0)(cfg)

	if cfg.timeout != defaultTimeout {
		t.Errorf(
			"expected timeout %v, got %v",
			defaultTimeout,
			cfg.timeout,
		)
	}
}

func TestWithOptionsCombination(t *testing.T) {
	cfg := newConfig()

	customClient := &http.Client{
		Timeout: 15 * time.Second,
	}

	opts := []Option{
		WithBaseURL("https://zammad.example.com/"),
		WithToken("test-token"),
		WithHTTPClient(customClient),
		WithTimeout(20 * time.Second),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.baseURL != "https://zammad.example.com" {
		t.Errorf("unexpected baseURL: %q", cfg.baseURL)
	}

	if cfg.token != "test-token" {
		t.Errorf("unexpected token: %q", cfg.token)
	}

	if cfg.httpClient != customClient {
		t.Error("expected custom HTTP client to be preserved")
	}

	if cfg.timeout != 20*time.Second {
		t.Errorf("unexpected timeout: %v", cfg.timeout)
	}

	if customClient.Timeout != 15*time.Second {
		t.Errorf("custom HTTP client timeout was modified: %v", customClient.Timeout)
	}
}
