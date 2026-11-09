package ticket

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

func TestTicketPriorityService_List(t *testing.T) {
	client := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
		}

		if r.URL.Path != "/api/v1/ticket_priorities" {
			t.Errorf("expected path /api/v1/ticket_priorities, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`[
			{
				"id": 1,
				"name": "1 low",
				"default_create": false,
				"active": true
			},
			{
				"id": 2,
				"name": "2 normal",
				"default_create": true,
				"active": true
			}
		]`))
	}))

	result, err := client.GetPriorities(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}

	if len(result.Priorities) != 2 {
		t.Fatalf("expected 2 priorities, got %d", len(result.Priorities))
	}

	if result.Priorities[0].Name != "1 low" {
		t.Errorf("expected priority name 1 low, got %s", result.Priorities[0].Name)
	}

	if result.Priorities[1].Name != "2 normal" {
		t.Errorf("expected priority name 2 normal, got %s", result.Priorities[1].Name)
	}

	if !result.Priorities[1].DefaultCreate {
		t.Error("expected second priority to be default create")
	}
}

func TestTicketPriorityService_Get(t *testing.T) {
	client := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
		}

		if r.URL.Path != "/api/v1/ticket_priorities/2" {
			t.Errorf("expected path /api/v1/ticket_priorities/2, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`{
			"id": 2,
			"name": "2 normal",
			"default_create": true,
			"active": true
		}`))
	}))

	result, err := client.GetPriority(context.Background(), 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != 2 {
		t.Errorf("expected ID 2, got %d", result.ID)
	}

	if result.Name != "2 normal" {
		t.Errorf("expected name 2 normal, got %s", result.Name)
	}

	if !result.DefaultCreate {
		t.Error("expected default_create to be true")
	}

	if !result.Active {
		t.Error("expected priority to be active")
	}
}

func TestTicketPriorityService_Get_InvalidID(t *testing.T) {
	client := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("HTTP request should not be executed for invalid ID")
	}))

	_, err := client.GetPriority(context.Background(), 0)
	if err == nil {
		t.Fatal("expected error for invalid priority ID")
	}
}
