package models

// DefaultPerPage is the page size used when ListOptions.PerPage is zero.
const DefaultPerPage = 100

// ListOptions selects a single page of a paginated Zammad endpoint.
// Zero values fall back to page 1 and DefaultPerPage.
type ListOptions struct {
	Page    int
	PerPage int
}
