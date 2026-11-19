package ticket

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/andriantp/zammad-sdk/models"
)

func TestService_CreateTicket(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method %s, got %s", http.MethodPost, r.Method)
		}

		if r.URL.Path != "/api/v1/tickets" {
			t.Errorf("expected path /api/v1/tickets, got %s", r.URL.Path)
		}

		var input models.TicketInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if input.Title != "Cannot access dashboard" {
			t.Errorf("expected title Cannot access dashboard, got %s", input.Title)
		}

		if input.Customer != "customer@example.com" {
			t.Errorf("expected customer customer@example.com, got %s", input.Customer)
		}

		if input.Group != "Support" {
			t.Errorf("expected group Support, got %s", input.Group)
		}

		if input.Tags != "bug,urgent" {
			t.Errorf("expected tags bug,urgent, got %s", input.Tags)
		}

		if input.Article.Subject != "Dashboard access issue" {
			t.Errorf("unexpected article subject: %s", input.Article.Subject)
		}

		if input.Article.Body != "I cannot access my dashboard." {
			t.Errorf("unexpected article body: %s", input.Article.Body)
		}

		if input.Article.Type != "note" {
			t.Errorf("expected article type note, got %s", input.Article.Type)
		}

		if input.Article.Sender != "Customer" {
			t.Errorf("expected article sender Customer, got %s", input.Article.Sender)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_, _ = w.Write([]byte(`{
			"id": 123,
			"number": "1000123",
			"title": "Cannot access dashboard",
			"customer": "customer@example.com",
			"group": "Support",
			"tags": "bug,urgent",
			"article": {
				"id": 456,
				"ticket_id": 123,
				"subject": "Dashboard access issue",
				"body": "I cannot access my dashboard.",
				"type": "note",
				"sender": "Customer"
			}
		}`))
	}))

	input := models.TicketInput{
		Title:    "Cannot access dashboard",
		Customer: "customer@example.com",
		Group:    "Support",
		Tags:     "bug,urgent",
		Article: models.ArticleInput{
			Subject: "Dashboard access issue",
			Body:    "I cannot access my dashboard.",
			Type:    "note",
			Sender:  "Customer",
		},
	}

	result, err := service.CreateTicket(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != 123 {
		t.Errorf("expected ticket ID 123, got %d", result.ID)
	}

	if result.Number != "1000123" {
		t.Errorf("expected ticket number 1000123, got %s", result.Number)
	}

	if result.Title != input.Title {
		t.Errorf("expected title %s, got %s", input.Title, result.Title)
	}

	if result.Article.ID != 456 {
		t.Errorf("expected article ID 456, got %d", result.Article.ID)
	}

	if result.Article.Subject != input.Article.Subject {
		t.Errorf("expected article subject %s, got %s", input.Article.Subject, result.Article.Subject)
	}
}

func TestService_GetTicket(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/tags" {
			if r.URL.Query().Get("object") != "Ticket" || r.URL.Query().Get("o_id") != "123" {
				t.Errorf("unexpected tags query %s", r.URL.RawQuery)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"tags": ["sdk-zammad", "urgent"]}`))

			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
		}

		if r.URL.Path != "/api/v1/tickets/123" {
			t.Errorf("expected path /api/v1/tickets/123, got %s", r.URL.Path)
		}

		if got := r.URL.Query().Get("expand"); got != "true" {
			t.Errorf("expected expand=true, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`{
			"id": 123,
			"number": "1000123",
			"title": "Cannot access dashboard",
			"customer": "customer@example.com",
			"group": "Support",
			"state": "new",
			"priority": "2 normal",
			"owner": "-",
			"organization": "Acme",
			"created_by": "agent@example.com",
			"updated_by": "agent@example.com",
			"priority_id": 2,
			"state_id": 1,
			"article": {
				"id": 456,
				"ticket_id": 123,
				"subject": "Dashboard access issue",
				"body": "I cannot access my dashboard.",
				"type": "note",
				"sender": "Customer"
			}
		}`))
	}))

	result, err := service.GetTicket(context.Background(), 123)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Tags != "sdk-zammad,urgent" {
		t.Errorf("expected tags from tags endpoint, got %q", result.Tags)
	}

	if result.Group != "Support" || result.State != "new" || result.Priority != "2 normal" ||
		result.Owner != "-" || result.Organization != "Acme" ||
		result.CreatedBy != "agent@example.com" || result.UpdatedBy != "agent@example.com" {
		t.Errorf("expanded name fields not decoded: %+v", result)
	}

	if result.ID != 123 {
		t.Errorf("expected ID 123, got %d", result.ID)
	}

	if result.Number != "1000123" {
		t.Errorf("expected number 1000123, got %s", result.Number)
	}

	if result.Title != "Cannot access dashboard" {
		t.Errorf("unexpected title: %s", result.Title)
	}

	if result.Customer != "customer@example.com" {
		t.Errorf("unexpected customer: %s", result.Customer)
	}

	if result.PriorityID != 2 {
		t.Errorf("expected priority ID 2, got %d", result.PriorityID)
	}

	if result.Article.ID != 456 {
		t.Errorf("expected article ID 456, got %d", result.Article.ID)
	}
}

func TestService_GetTicket_InvalidID(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("HTTP request should not be executed for invalid ID")
	}))

	_, err := service.GetTicket(context.Background(), 0)
	if err == nil {
		t.Fatal("expected error for invalid ticket ID")
	}
}

func TestService_GetTicketByNumber(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/tags" {
			if r.URL.Query().Get("object") != "Ticket" || r.URL.Query().Get("o_id") != "123" {
				t.Errorf("unexpected tags query %s", r.URL.RawQuery)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"tags": ["sdk-zammad", "urgent"]}`))

			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
		}

		if r.URL.Path != "/api/v1/tickets/search" {
			t.Errorf("expected path /api/v1/tickets/search, got %s", r.URL.Path)
		}

		if query := r.URL.Query().Get("query"); query != "number:1000123" {
			t.Errorf("expected query number:1000123, got %s", query)
		}

		if got := r.URL.Query().Get("expand"); got != "true" {
			t.Errorf("expected expand=true, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`[
			{
				"id": 123,
				"number": "1000123",
				"title": "Cannot access dashboard",
				"customer": "customer@example.com",
				"group": "Support"
			}
		]`))
	}))

	result, err := service.GetTicketByNumber(context.Background(), "1000123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Tags != "sdk-zammad,urgent" {
		t.Errorf("expected tags from tags endpoint, got %q", result.Tags)
	}

	if result.ID != 123 {
		t.Errorf("expected ID 123, got %d", result.ID)
	}

	if result.Number != "1000123" {
		t.Errorf("expected number 1000123, got %s", result.Number)
	}

	if result.Title != "Cannot access dashboard" {
		t.Errorf("unexpected title: %s", result.Title)
	}
}

func TestService_GetTicketByNumber_NotFound(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`[]`))
	}))

	_, err := service.GetTicketByNumber(context.Background(), "9999999")
	if err == nil {
		t.Fatal("expected error for ticket not found")
	}
}

func TestService_GetTicketByNumber_EmptyNumber(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("HTTP request should not be executed for empty ticket number")
	}))

	_, err := service.GetTicketByNumber(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty ticket number")
	}
}

func TestService_GetTicketByNumber_WhitespaceNumber(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("HTTP request should not be executed for whitespace ticket number")
	}))

	_, err := service.GetTicketByNumber(context.Background(), "   ")
	if err == nil {
		t.Fatal("expected error for whitespace ticket number")
	}
}

func TestService_GetTicketByNumber_ExactMatch(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/tags" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"tags": []}`))

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`[
			{
				"id": 123,
				"number": "10001230",
				"title": "Another ticket"
			},
			{
				"id": 124,
				"number": "1000123",
				"title": "Target ticket"
			}
		]`))
	}))

	result, err := service.GetTicketByNumber(context.Background(), "1000123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != 124 {
		t.Errorf("expected ID 124, got %d", result.ID)
	}

	if result.Title != "Target ticket" {
		t.Errorf("expected title Target ticket, got %s", result.Title)
	}
}

func TestService_SearchTicketsAll(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/tickets/search" {
			t.Fatalf("expected path /api/v1/tickets/search, got %s", r.URL.Path)
		}

		if got := r.URL.Query().Get("query"); got != "title:login" {
			t.Fatalf("expected query title:login, got %s", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
            {
                "id": 123,
                "number": "1000123",
                "title": "Login issue",
                "customer": "customer@example.com",
                "group": "Support",
                "priority_id": 2,
                "state_id": 1,
                "tags": "bug,urgent"
            },
            {
                "id": 124,
                "number": "1000124",
                "title": "Login failure",
                "customer": "user@example.com",
                "group": "Support",
                "priority_id": 3,
                "state_id": 2,
                "tags": "login"
            }
        ]`))
	}))

	result, err := service.SearchTicketsAll(context.Background(), "title:login")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Total != 2 {
		t.Fatalf("expected total 2, got %d", result.Total)
	}

	if len(result.Tickets) != 2 {
		t.Fatalf("expected 2 tickets, got %d", len(result.Tickets))
	}

	if result.Tickets[0].ID != 123 {
		t.Errorf("expected ticket ID 123, got %d", result.Tickets[0].ID)
	}

	if result.Tickets[0].Number != "1000123" {
		t.Errorf("expected ticket number 1000123, got %s", result.Tickets[0].Number)
	}

	if result.Tickets[1].ID != 124 {
		t.Errorf("expected ticket ID 124, got %d", result.Tickets[1].ID)
	}
}

func TestService_SearchTicketsAll_EmptyQuery(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected no HTTP request")
	}))

	result, err := service.SearchTicketsAll(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty query")
	}

	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestService_SearchTicketsAll_WhitespaceQuery(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected no HTTP request")
	}))

	result, err := service.SearchTicketsAll(context.Background(), "   ")
	if err == nil {
		t.Fatal("expected error for whitespace query")
	}

	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestService_UpdateTicket(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected method PUT, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/tickets/123" {
			t.Fatalf("expected path /api/v1/tickets/123, got %s", r.URL.Path)
		}

		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request body: %v", err)
		}

		if payload["title"] != "Updated ticket title" {
			t.Errorf("expected updated title, got %v", payload["title"])
		}

		if payload["priority_id"] != float64(3) {
			t.Errorf("expected priority_id 3, got %v", payload["priority_id"])
		}

		if payload["state_id"] != float64(2) {
			t.Errorf("expected state_id 2, got %v", payload["state_id"])
		}

		if payload["tags"] != "urgent,bug" {
			t.Errorf("expected tags urgent,bug, got %v", payload["tags"])
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))

	payload := map[string]any{
		"title":       "Updated ticket title",
		"priority_id": 3,
		"state_id":    2,
		"tags":        "urgent,bug",
	}

	err := service.UpdateTicket(context.Background(), 123, payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestService_UpdateTicket_InvalidID(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected no HTTP request")
	}))

	err := service.UpdateTicket(context.Background(), 0, map[string]any{
		"title": "Updated ticket title",
	})

	if err == nil {
		t.Fatal("expected error for invalid ticket ID")
	}
}

func TestService_UpdateTicket_EmptyPayload(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected no HTTP request")
	}))

	err := service.UpdateTicket(context.Background(), 123, map[string]any{})

	if err == nil {
		t.Fatal("expected error for empty payload")
	}
}

func TestService_UpdateTicket_APIError(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{
			"error": "Unprocessable Entity",
			"message": "Validation failed"
		}`))
	}))

	err := service.UpdateTicket(context.Background(), 123, map[string]any{
		"title": "Updated ticket title",
	})

	if err == nil {
		t.Fatal("expected API error")
	}
}

func TestService_DeleteTicket(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected method DELETE, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/tickets/123" {
			t.Fatalf("expected path /api/v1/tickets/123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))

	err := service.DeleteTicket(context.Background(), 123)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestService_DeleteTicket_InvalidID(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected no HTTP request")
	}))

	err := service.DeleteTicket(context.Background(), 0)
	if err == nil {
		t.Fatal("expected error for invalid ticket ID")
	}
}

func TestService_DeleteTicket_APIError(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{
			"error": "Not Found",
			"message": "Ticket not found"
		}`))
	}))

	err := service.DeleteTicket(context.Background(), 123)
	if err == nil {
		t.Fatal("expected API error")
	}
}

func TestService_SearchTicketsAll_Paginates(t *testing.T) {
	const total = 103
	var pages []string

	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		pages = append(pages, q.Get("page"))

		if q.Get("query") != "state:open" {
			t.Errorf("expected query state:open, got %q", q.Get("query"))
		}

		page, _ := strconv.Atoi(q.Get("page"))
		perPage, _ := strconv.Atoi(q.Get("per_page"))

		tickets := []map[string]any{}
		for id := (page-1)*perPage + 1; id <= total && id <= page*perPage; id++ {
			tickets = append(tickets, map[string]any{"id": id})
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(tickets)
	}))

	result, err := service.SearchTicketsAll(context.Background(), "state:open")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Total != total || len(result.Tickets) != total {
		t.Fatalf("expected %d tickets, got total=%d len=%d", total, result.Total, len(result.Tickets))
	}

	if strings.Join(pages, ",") != "1,2,3" {
		t.Errorf("expected pages 1,2,3, got %v", pages)
	}
}

func TestService_SearchTicketsPage(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		if r.URL.Path != "/api/v1/tickets/search" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}

		if q.Get("query") != "title:login" || q.Get("page") != "2" || q.Get("per_page") != "10" {
			t.Errorf("unexpected query %s", r.URL.RawQuery)
		}

		if got := r.URL.Query().Get("expand"); got != "true" {
			t.Errorf("expected expand=true, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id": 11}, {"id": 12}]`))
	}))

	result, err := service.SearchTicketsPage(context.Background(), "title:login", models.ListOptions{Page: 2, PerPage: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Total != 2 || result.Tickets[0].ID != 11 {
		t.Errorf("unexpected result %+v", result)
	}
}

func TestService_SearchTicketsPage_InvalidInput(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected no HTTP request")
	}))

	if _, err := service.SearchTicketsPage(context.Background(), " ", models.ListOptions{}); err == nil {
		t.Error("expected error for empty query")
	}

	if _, err := service.SearchTicketsPage(context.Background(), "x", models.ListOptions{Page: -1}); err == nil {
		t.Error("expected error for negative page")
	}
}

// SearchTickets is a deprecated alias of SearchTicketsAll and must follow every page.
func TestService_SearchTicketsDeprecatedAlias(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Query().Get("page") {
		case "1":
			_, _ = w.Write([]byte(`[{"id": 1}]`))
		case "2":
			_, _ = w.Write([]byte(`[{"id": 2}]`))
		default:
			_, _ = w.Write([]byte(`[]`))
		}
	}))

	result, err := service.SearchTickets(context.Background(), "title:login")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Total != 2 {
		t.Errorf("expected total 2 across pages, got %d", result.Total)
	}
}

func TestService_CountTickets(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		if r.URL.Path != "/api/v1/tickets/search" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}

		if q.Get("query") != "state.name:open" || q.Get("full") != "true" ||
			q.Get("per_page") != "1" || q.Get("page") != "1" {
			t.Errorf("unexpected query %s", r.URL.RawQuery)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"record_ids": [58], "assets": {}, "total_count": 53}`))
	}))

	count, err := service.CountTickets(context.Background(), "  state.name:open ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if count != 53 {
		t.Errorf("expected 53, got %d", count)
	}
}

func TestService_CountTickets_Zero(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"record_ids": [], "assets": {}, "total_count": 0}`))
	}))

	count, err := service.CountTickets(context.Background(), "nomatch")
	if err != nil || count != 0 {
		t.Errorf("expected 0 without error, got %d (%v)", count, err)
	}
}

// A response without total_count (e.g. a Zammad version that ignores
// full=true) must be an error, not a silent count of 0.
func TestService_CountTickets_MissingTotal(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"record_ids": [], "assets": {}}`))
	}))

	if _, err := service.CountTickets(context.Background(), "x"); err == nil {
		t.Error("expected error when total_count is missing")
	}
}

func TestService_CountTickets_EmptyQuery(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected no HTTP request")
	}))

	if _, err := service.CountTickets(context.Background(), "  "); err == nil {
		t.Error("expected error for empty query")
	}
}
