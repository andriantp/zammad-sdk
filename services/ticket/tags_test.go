package ticket

import (
	"context"
	"net/http"
	"testing"
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
