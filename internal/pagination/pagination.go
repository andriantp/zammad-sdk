// Package pagination implements Zammad's page/per_page pagination.
//
// Zammad does not report a total count, so the only way to retrieve every
// record is to request pages until one comes back empty.
package pagination

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/andriantp/zammad-sdk/models"
)

// Normalize validates opts and applies defaults.
func Normalize(opts models.ListOptions) (models.ListOptions, error) {
	if opts.Page < 0 {
		return opts, fmt.Errorf("zammad: invalid page: %d", opts.Page)
	}

	if opts.PerPage < 0 {
		return opts, fmt.Errorf("zammad: invalid per page: %d", opts.PerPage)
	}

	if opts.Page == 0 {
		opts.Page = 1
	}

	if opts.PerPage == 0 {
		opts.PerPage = models.DefaultPerPage
	}

	return opts, nil
}

// Path appends page and per_page (plus any extra parameters) to path.
func Path(path string, opts models.ListOptions, extra url.Values) string {
	values := url.Values{}
	for key, vals := range extra {
		values[key] = vals
	}

	values.Set("page", strconv.Itoa(opts.Page))
	values.Set("per_page", strconv.Itoa(opts.PerPage))

	return path + "?" + values.Encode()
}

// FetchFunc retrieves a single page.
type FetchFunc[T any] func(ctx context.Context, opts models.ListOptions) ([]T, error)

// All requests pages starting at 1 until a page is empty or contains no
// records that have not been seen already, and returns the collected records
// with duplicates (by id) removed.
//
// Stopping on an empty page rather than a short page keeps the result
// complete even if the server caps per_page below the requested size; the
// cost is one extra request. The "no new records" check guards against
// servers that ignore the page parameter, which would otherwise loop forever.
func All[T any](ctx context.Context, fetch FetchFunc[T], id func(T) int) ([]T, error) {
	seen := make(map[int]struct{})
	all := []T{}

	for page := 1; ; page++ {
		items, err := fetch(ctx, models.ListOptions{
			Page:    page,
			PerPage: models.DefaultPerPage,
		})
		if err != nil {
			return nil, fmt.Errorf("fetch page %d: %w", page, err)
		}

		added := 0

		for _, item := range items {
			key := id(item)
			if _, ok := seen[key]; ok {
				continue
			}

			seen[key] = struct{}{}
			all = append(all, item)
			added++
		}

		if added == 0 {
			return all, nil
		}
	}
}
