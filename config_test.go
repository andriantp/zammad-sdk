package zammad

import (
	"testing"
	"time"
)

func TestNewConfig(t *testing.T) {
	cfg := newConfig()

	if cfg == nil {
		t.Fatal("expected config, got nil")
	}

	if cfg.baseURL != "" {
		t.Errorf("expected empty baseURL, got %q", cfg.baseURL)
	}

	if cfg.token != "" {
		t.Errorf("expected empty token, got %q", cfg.token)
	}

	if cfg.httpClient != nil {
		t.Error("expected nil httpClient")
	}

	if cfg.timeout != defaultTimeout {
		t.Errorf(
			"expected timeout %v, got %v",
			defaultTimeout,
			cfg.timeout,
		)
	}
}

func TestDefaultTimeout(t *testing.T) {
	expected := 30 * time.Second

	if defaultTimeout != expected {
		t.Errorf(
			"expected default timeout %v, got %v",
			expected,
			defaultTimeout,
		)
	}
}
