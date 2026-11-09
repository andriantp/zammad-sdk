package user

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

func TestUserService_List(t *testing.T) {
	client := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
		}

		if r.URL.Path != "/api/v1/users" {
			t.Errorf("expected path /api/v1/users, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`[
			{
				"id": 1,
				"login": "admin",
				"firstname": "Admin",
				"lastname": "User",
				"email": "admin@example.com",
				"active": true,
				"role_ids": [1]
			},
			{
				"id": 2,
				"login": "agent",
				"firstname": "Agent",
				"lastname": "User",
				"email": "agent@example.com",
				"active": true,
				"role_ids": [2]
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

	if len(result.Users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(result.Users))
	}

	if result.Users[0].Login != "admin" {
		t.Errorf("expected login admin, got %s", result.Users[0].Login)
	}

	if result.Users[1].Email != "agent@example.com" {
		t.Errorf("expected email agent@example.com, got %s", result.Users[1].Email)
	}

	if len(result.Users[0].RoleIDs) != 1 || result.Users[0].RoleIDs[0] != 1 {
		t.Errorf("unexpected role IDs: %v", result.Users[0].RoleIDs)
	}
}

func TestUserService_Get(t *testing.T) {
	client := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
		}

		if r.URL.Path != "/api/v1/users/123" {
			t.Errorf("expected path /api/v1/users/123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`{
			"id": 123,
			"login": "andrian",
			"firstname": "Andrian",
			"lastname": "User",
			"email": "andrian@example.com",
			"active": true,
			"role_ids": [1, 2],
			"group_ids": {
				"1": ["1", "2"]
			}
		}`))
	}))

	result, err := client.Get(context.Background(), 123)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != 123 {
		t.Errorf("expected ID 123, got %d", result.ID)
	}

	if result.Login != "andrian" {
		t.Errorf("expected login andrian, got %s", result.Login)
	}

	if result.Email != "andrian@example.com" {
		t.Errorf("expected email andrian@example.com, got %s", result.Email)
	}

	if !result.Active {
		t.Error("expected user to be active")
	}

	if len(result.RoleIDs) != 2 {
		t.Errorf("expected 2 role IDs, got %d", len(result.RoleIDs))
	}

	if len(result.GroupIDs["1"]) != 2 {
		t.Errorf("expected 2 group IDs, got %d", len(result.GroupIDs["1"]))
	}
}

func TestUserService_Get_InvalidID(t *testing.T) {
	client := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("HTTP request should not be executed for invalid ID")
	}))

	_, err := client.Get(context.Background(), 0)
	if err == nil {
		t.Fatal("expected error for invalid user ID")
	}
}
