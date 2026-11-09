package ticket

import (
	"context"
	"net/http"
	"testing"
)

func TestService_GetStates(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
		}

		if r.URL.Path != "/api/v1/ticket_states" {
			t.Errorf("expected path /api/v1/ticket_states, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`[
			{
				"id": 1,
				"state_type_id": 1,
				"name": "new",
				"default_create": true,
				"active": true
			},
			{
				"id": 2,
				"state_type_id": 2,
				"name": "open",
				"default_follow_up": true,
				"active": true
			}
		]`))
	}))

	result, err := service.GetStates(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}

	if len(result.States) != 2 {
		t.Fatalf("expected 2 states, got %d", len(result.States))
	}

	if result.States[0].Name != "new" {
		t.Errorf("expected state name new, got %s", result.States[0].Name)
	}

	if !result.States[0].DefaultCreate {
		t.Error("expected first state default_create to be true")
	}

	if result.States[1].Name != "open" {
		t.Errorf("expected state name open, got %s", result.States[1].Name)
	}

	if !result.States[1].DefaultFollowUp {
		t.Error("expected second state default_follow_up to be true")
	}
}

func TestService_GetState(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
		}

		if r.URL.Path != "/api/v1/ticket_states/2" {
			t.Errorf("expected path /api/v1/ticket_states/2, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`{
			"id": 2,
			"state_type_id": 2,
			"name": "open",
			"default_follow_up": true,
			"active": true
		}`))
	}))

	result, err := service.GetState(context.Background(), 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != 2 {
		t.Errorf("expected ID 2, got %d", result.ID)
	}

	if result.Name != "open" {
		t.Errorf("expected name open, got %s", result.Name)
	}

	if !result.DefaultFollowUp {
		t.Error("expected default_follow_up to be true")
	}

	if !result.Active {
		t.Error("expected state to be active")
	}
}

func TestService_GetState_InvalidID(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("HTTP request should not be executed for invalid ID")
	}))

	_, err := service.GetState(context.Background(), 0)
	if err == nil {
		t.Fatal("expected error for invalid state ID")
	}
}
