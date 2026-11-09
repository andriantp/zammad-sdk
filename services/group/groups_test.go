package group

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/andriantp/zammad-sdk/internal/transport"
)

func newTestService(t *testing.T, handler http.Handler) *Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	transportClient := transport.NewClient(transport.Config{
		BaseURL: server.URL,
		Token:   "test-token",
		Timeout: 30 * time.Second,
	})

	return New(transportClient)
}

func TestGroupService_List(t *testing.T) {
	client := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
		}

		if r.URL.Path != "/api/v1/groups" {
			t.Errorf("expected path /api/v1/groups, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`[
			{
				"id": 1,
				"name": "Admin",
				"active": true
			},
			{
				"id": 2,
				"name": "Users",
				"active": true
			}
		]`))
	}))

	result, err := client.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}

	if len(result.Groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(result.Groups))
	}

	if result.Groups[0].Name != "Admin" {
		t.Errorf("expected group name Admin, got %s", result.Groups[0].Name)
	}

	if result.Groups[1].Name != "Users" {
		t.Errorf("expected group name Users, got %s", result.Groups[1].Name)
	}
}

func TestGroupService_Get(t *testing.T) {
	client := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
		}

		if r.URL.Path != "/api/v1/groups/123" {
			t.Errorf("expected path /api/v1/groups/123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`{
			"id": 123,
			"name": "Admin",
			"active": true
		}`))
	}))

	result, err := client.Get(context.Background(), 123)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != 123 {
		t.Errorf("expected ID 123, got %d", result.ID)
	}

	if result.Name != "Admin" {
		t.Errorf("expected name Admin, got %s", result.Name)
	}

	if !result.Active {
		t.Error("expected group to be active")
	}
}

func TestGroupService_Get_InvalidID(t *testing.T) {
	client := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("HTTP request should not be executed for invalid ID")
	}))

	_, err := client.Get(context.Background(), 0)
	if err == nil {
		t.Fatal("expected error for invalid group ID")
	}
}
