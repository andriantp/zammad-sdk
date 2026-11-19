package pagination

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/andriantp/zammad-sdk/models"
)

type item struct{ ID int }

func itemID(i item) int { return i.ID }

// pagedSource serves ids 1..total in pages, optionally capping the page size
// the way a server with a lower max per_page would.
func pagedSource(total, maxPerPage int, calls *int) FetchFunc[item] {
	return func(_ context.Context, opts models.ListOptions) ([]item, error) {
		*calls++

		perPage := opts.PerPage
		if maxPerPage > 0 && perPage > maxPerPage {
			perPage = maxPerPage
		}

		var out []item
		for id := (opts.Page-1)*perPage + 1; id <= total && id <= opts.Page*perPage; id++ {
			out = append(out, item{ID: id})
		}

		return out, nil
	}
}

func TestNormalize(t *testing.T) {
	got, err := Normalize(models.ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Page != 1 || got.PerPage != models.DefaultPerPage {
		t.Errorf("expected defaults {1 %d}, got %+v", models.DefaultPerPage, got)
	}

	got, err = Normalize(models.ListOptions{Page: 3, PerPage: 25})
	if err != nil || got.Page != 3 || got.PerPage != 25 {
		t.Errorf("expected {3 25}, got %+v (err %v)", got, err)
	}

	for _, opts := range []models.ListOptions{{Page: -1}, {PerPage: -1}} {
		if _, err := Normalize(opts); err == nil {
			t.Errorf("expected error for %+v", opts)
		}
	}
}

func TestPath(t *testing.T) {
	got := Path("/tickets/search", models.ListOptions{Page: 2, PerPage: 50}, url.Values{"query": {"title:a b&c"}})

	u, err := url.Parse(got)
	if err != nil {
		t.Fatalf("parse %q: %v", got, err)
	}

	q := u.Query()
	if u.Path != "/tickets/search" || q.Get("page") != "2" || q.Get("per_page") != "50" || q.Get("query") != "title:a b&c" {
		t.Errorf("unexpected path %q", got)
	}
}

func TestAll_MultiplePages(t *testing.T) {
	calls := 0

	got, err := All(context.Background(), pagedSource(253, 0, &calls), itemID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 253 {
		t.Fatalf("expected 253 items, got %d", len(got))
	}

	for i, it := range got {
		if it.ID != i+1 {
			t.Fatalf("item %d: expected ID %d, got %d", i, i+1, it.ID)
		}
	}

	// 3 pages with data + 1 empty page.
	if calls != 4 {
		t.Errorf("expected 4 requests, got %d", calls)
	}
}

func TestAll_Empty(t *testing.T) {
	calls := 0

	got, err := All(context.Background(), pagedSource(0, 0, &calls), itemID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got == nil || len(got) != 0 {
		t.Errorf("expected empty non-nil slice, got %#v", got)
	}

	if calls != 1 {
		t.Errorf("expected 1 request, got %d", calls)
	}
}

func TestAll_ServerCapsPerPage(t *testing.T) {
	calls := 0

	got, err := All(context.Background(), pagedSource(120, 50, &calls), itemID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 120 {
		t.Errorf("expected 120 items despite server cap, got %d", len(got))
	}
}

func TestAll_ServerIgnoresPage(t *testing.T) {
	calls := 0
	fetch := func(context.Context, models.ListOptions) ([]item, error) {
		calls++
		return []item{{1}, {2}}, nil
	}

	got, err := All(context.Background(), fetch, itemID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 2 || calls != 2 {
		t.Errorf("expected 2 items in 2 requests, got %d items in %d requests", len(got), calls)
	}
}

func TestAll_DeduplicatesShiftedPages(t *testing.T) {
	// A new record arriving between requests shifts page 2 by one,
	// so its first item repeats the last item of page 1.
	pages := map[int][]item{
		1: {{5}, {4}},
		2: {{4}, {3}},
		3: {{1}},
	}
	fetch := func(_ context.Context, opts models.ListOptions) ([]item, error) {
		return pages[opts.Page], nil
	}

	got, err := All(context.Background(), fetch, itemID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []int{5, 4, 3, 1}
	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}

	for i := range want {
		if got[i].ID != want[i] {
			t.Fatalf("expected %v, got %v", want, got)
		}
	}
}

func TestAll_ErrorIsWrapped(t *testing.T) {
	sentinel := errors.New("boom")
	fetch := func(_ context.Context, opts models.ListOptions) ([]item, error) {
		if opts.Page == 2 {
			return nil, sentinel
		}
		return []item{{opts.Page}}, nil
	}

	got, err := All(context.Background(), fetch, itemID)
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected wrapped sentinel error, got %v", err)
	}

	if got != nil {
		t.Errorf("expected nil result on error, got %v", got)
	}
}
