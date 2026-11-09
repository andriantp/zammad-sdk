package organization

import (
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

func TestOrganizationList(t *testing.T) {
	client := newTestService(t, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}

			if r.URL.Path != "/api/v1/organizations" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}

			response := []models.Organization{
				{
					ID:     1,
					Name:   "Acme Corporation",
					Active: true,
				},
				{
					ID:     2,
					Name:   "Example Inc.",
					Active: true,
				},
			}

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
		},
	))

	organizations, err := client.List(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if organizations == nil {
		t.Fatal("expected organizations, got nil")
	}

	if organizations.Total != 2 {
		t.Fatalf("expected total 2, got %d", organizations.Total)
	}

	if len(organizations.Organizations) != 2 {
		t.Fatalf(
			"expected 2 organizations, got %d",
			len(organizations.Organizations),
		)
	}

	if organizations.Organizations[0].Name != "Acme Corporation" {
		t.Errorf(
			"unexpected organization name: %s",
			organizations.Organizations[0].Name,
		)
	}

	if organizations.Organizations[1].ID != 2 {
		t.Errorf(
			"expected ID 2, got %d",
			organizations.Organizations[1].ID,
		)
	}
}

func TestOrganizationGet(t *testing.T) {
	client := newTestService(t, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}

			if r.URL.Path != "/api/v1/organizations/123" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}

			response := models.Organization{
				ID:     123,
				Name:   "Acme Corporation",
				Active: true,
			}

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
		},
	))

	organization, err := client.Get(
		context.Background(),
		123,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if organization.ID != 123 {
		t.Errorf("expected ID 123, got %d", organization.ID)
	}

	if organization.Name != "Acme Corporation" {
		t.Errorf("unexpected name: %s", organization.Name)
	}
}

func TestOrganizationGetInvalidID(t *testing.T) {
	client := newTestService(t, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			t.Error("HTTP request should not be executed")
		},
	))

	_, err := client.Get(context.Background(), 0)

	if err == nil {
		t.Fatal("expected error for invalid organization ID")
	}
}
