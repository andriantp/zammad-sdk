package ticket

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/andriantp/zammad-sdk/internal/transport"
)

func TestService_GetTags(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
		}

		if r.URL.Path != "/api/v1/tag_list" {
			t.Errorf("expected path /api/v1/tag_list, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`[
			{
				"id": 1,
				"name": "americano",
				"count": 0
			},
			{
				"id": 2,
				"name": "complaint",
				"count": 5
			}
		]`))
	}))

	result, err := service.GetTags(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}

	if len(result.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(result.Tags))
	}

	if result.Tags[0].Name != "americano" {
		t.Errorf("expected tag name americano, got %s", result.Tags[0].Name)
	}

	if result.Tags[0].Count != 0 {
		t.Errorf("expected count 0, got %d", result.Tags[0].Count)
	}

	if result.Tags[1].Name != "complaint" {
		t.Errorf("expected tag name complaint, got %s", result.Tags[1].Name)
	}

	if result.Tags[1].Count != 5 {
		t.Errorf("expected count 5, got %d", result.Tags[1].Count)
	}
}

func TestService_GetTicketTags(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tags" {
			t.Errorf("expected path /api/v1/tags, got %s", r.URL.Path)
		}

		if q := r.URL.Query(); q.Get("object") != "Ticket" || q.Get("o_id") != "58" {
			t.Errorf("unexpected query %s", r.URL.RawQuery)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"tags": ["sdk-zammad", "urgent"]}`))
	}))

	tags, err := service.GetTicketTags(context.Background(), 58)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tags) != 2 || tags[0] != "sdk-zammad" || tags[1] != "urgent" {
		t.Errorf("unexpected tags %v", tags)
	}
}

func TestService_GetTicketTags_Empty(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"tags": []}`))
	}))

	tags, err := service.GetTicketTags(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tags == nil || len(tags) != 0 {
		t.Errorf("expected empty non-nil slice, got %#v", tags)
	}
}

func TestService_GetTicketTags_InvalidID(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected no HTTP request")
	}))

	if _, err := service.GetTicketTags(context.Background(), 0); err == nil {
		t.Fatal("expected error for invalid ticket ID")
	}
}

// A failing tags request must fail GetTicket instead of returning a ticket
// with silently missing tags.
func TestService_GetTicket_TagsError(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path == "/api/v1/tags" {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"error": "Not authorized"}`))

			return
		}

		_, _ = w.Write([]byte(`{"id": 123}`))
	}))

	ticket, err := service.GetTicket(context.Background(), 123)
	if err == nil {
		t.Fatalf("expected error, got ticket %+v", ticket)
	}

	var apiErr *transport.APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusForbidden {
		t.Errorf("expected wrapped 403 APIError, got %v", err)
	}
}
