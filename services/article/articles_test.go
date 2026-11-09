package article

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/andriantp/zammad-sdk/internal/transport"
	"github.com/andriantp/zammad-sdk/models"
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

func TestService_GetArticles(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/ticket_articles/by_ticket/123" {
			t.Fatalf(
				"expected path /api/v1/ticket_articles/by_ticket/123, got %s",
				r.URL.Path,
			)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]map[string]any{
			{
				"id":         456,
				"ticket_id":  123,
				"subject":    "Login issue",
				"body":       "Unable to login",
				"type":       "note",
				"sender":     "Customer",
				"internal":   false,
				"created_by": "customer@example.com",
			},
			{
				"id":        457,
				"ticket_id": 123,
				"subject":   "Re: Login issue",
				"body":      "We are investigating the issue",
				"type":      "note",
				"sender":    "Agent",
				"internal":  false,
			},
		})
	}))

	result, err := service.GetArticles(context.Background(), 123)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}

	if len(result.Articles) != 2 {
		t.Fatalf("expected 2 articles, got %d", len(result.Articles))
	}

	if result.Articles[0].ID != 456 {
		t.Errorf("expected article ID 456, got %d", result.Articles[0].ID)
	}

	if result.Articles[0].Subject != "Login issue" {
		t.Errorf("expected subject Login issue, got %s", result.Articles[0].Subject)
	}

	if result.Articles[1].ID != 457 {
		t.Errorf("expected article ID 457, got %d", result.Articles[1].ID)
	}
}

func TestService_GetArticles_EmptyResult(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))

	result, err := service.GetArticles(context.Background(), 123)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Total != 0 {
		t.Errorf("expected total 0, got %d", result.Total)
	}

	if len(result.Articles) != 0 {
		t.Errorf("expected 0 articles, got %d", len(result.Articles))
	}
}

func TestService_GetArticles_InvalidTicketID(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected no HTTP request")
	}))

	result, err := service.GetArticles(context.Background(), 0)
	if err == nil {
		t.Fatal("expected error for invalid ticket ID")
	}

	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestService_GetArticles_APIError(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{
			"error": "Not Found",
			"message": "Ticket not found"
		}`))
	}))

	result, err := service.GetArticles(context.Background(), 123)
	if err == nil {
		t.Fatal("expected API error")
	}

	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestService_Reply(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/ticket_articles" {
			t.Fatalf("expected path /api/v1/ticket_articles, got %s", r.URL.Path)
		}

		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request body: %v", err)
		}

		if payload["ticket_id"] != float64(123) {
			t.Errorf("expected ticket_id 123, got %v", payload["ticket_id"])
		}

		if payload["subject"] != "Re: Login issue" {
			t.Errorf("expected subject Re: Login issue, got %v", payload["subject"])
		}

		if payload["body"] != "We are investigating the issue" {
			t.Errorf("unexpected body: %v", payload["body"])
		}

		if payload["type"] != "note" {
			t.Errorf("expected type note, got %v", payload["type"])
		}

		if payload["sender"] != "Agent" {
			t.Errorf("expected sender Agent, got %v", payload["sender"])
		}

		if payload["internal"] != false {
			t.Errorf("expected internal false, got %v", payload["internal"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"id": 456,
			"ticket_id": 123,
			"subject": "Re: Login issue",
			"body": "We are investigating the issue",
			"type": "note",
			"sender": "Agent",
			"internal": false
		}`))
	}))

	input := models.ArticleInput{
		Subject:  "Re: Login issue",
		Body:     "We are investigating the issue",
		Type:     "note",
		Sender:   "Agent",
		Internal: false,
	}

	result, err := service.Reply(context.Background(), 123, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != 456 {
		t.Errorf("expected article ID 456, got %d", result.ID)
	}

	if result.TicketID != 123 {
		t.Errorf("expected ticket ID 123, got %d", result.TicketID)
	}
}

func TestService_Reply_InvalidTicketID(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected no HTTP request")
	}))

	result, err := service.Reply(context.Background(), 0, models.ArticleInput{
		Body: "Test reply",
	})

	if err == nil {
		t.Fatal("expected error for invalid ticket ID")
	}

	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestService_Reply_EmptyBody(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected no HTTP request")
	}))

	result, err := service.Reply(context.Background(), 123, models.ArticleInput{
		Subject: "Test reply",
	})

	if err == nil {
		t.Fatal("expected error for empty article body")
	}

	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestService_Reply_APIError(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{
			"error": "Unprocessable Entity",
			"message": "Validation failed"
		}`))
	}))

	result, err := service.Reply(context.Background(), 123, models.ArticleInput{
		Subject: "Test reply",
		Body:    "Test body",
		Type:    "note",
		Sender:  "Agent",
	})

	if err == nil {
		t.Fatal("expected API error")
	}

	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestService_DownloadAttachment(t *testing.T) {
	expectedData := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected method GET, got %s", r.Method)
		}

		expectedPath := "/api/v1/ticket_attachment/123/456/789"
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(expectedData)
	}))

	result, err := service.DownloadAttachment(context.Background(), 123, 456, 789)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(result, expectedData) {
		t.Errorf("downloaded data does not match expected binary")
	}
}

func TestService_DownloadAttachment_InvalidTicketID(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected no HTTP request")
	}))

	result, err := service.DownloadAttachment(context.Background(), 0, 456, 789)
	if err == nil {
		t.Fatal("expected error for invalid ticket ID")
	}

	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestService_DownloadAttachment_InvalidArticleID(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected no HTTP request")
	}))

	result, err := service.DownloadAttachment(context.Background(), 123, 0, 789)
	if err == nil {
		t.Fatal("expected error for invalid article ID")
	}

	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestService_DownloadAttachment_InvalidAttachmentID(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected no HTTP request")
	}))

	result, err := service.DownloadAttachment(context.Background(), 123, 456, 0)
	if err == nil {
		t.Fatal("expected error for invalid attachment ID")
	}

	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestService_DownloadAttachment_APIError(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{
			"error": "Not Found",
			"message": "Attachment not found"
		}`))
	}))

	result, err := service.DownloadAttachment(context.Background(), 123, 456, 789)
	if err == nil {
		t.Fatal("expected API error")
	}

	if result != nil {
		t.Fatal("expected nil result")
	}
}
