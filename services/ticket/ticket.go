package ticket

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/andriantp/zammad-sdk/internal/pagination"
	"github.com/andriantp/zammad-sdk/models"
)

type Transport interface {
	Do(ctx context.Context, method, path string, body, result any) error
}

type Service struct {
	transport Transport
}

func New(transport Transport) *Service {
	return &Service{
		transport: transport,
	}
}

// CreateTicket creates a ticket with its first article.
//
// The ticket is requested with expand=true so name fields (Group, State,
// Customer, ...) are populated. Zammad omits tags from the response, so Tags
// is filled from input.Tags rather than with a second request: a failure
// after the ticket was created would otherwise invite a retry that creates
// a duplicate.
func (s *Service) CreateTicket(ctx context.Context, input models.TicketInput) (*models.Ticket, error) {
	var ticket models.Ticket

	err := s.transport.Do(
		ctx,
		http.MethodPost,
		"/tickets?expand=true",
		input,
		&ticket,
	)
	if err != nil {
		return nil, err
	}

	if ticket.Tags == "" {
		ticket.Tags = normalizeTags(input.Tags)
	}

	return &ticket, nil
}

// normalizeTags trims each comma-separated tag and drops empty ones.
func normalizeTags(tags string) string {
	var out []string

	for _, tag := range strings.Split(tags, ",") {
		if tag = strings.TrimSpace(tag); tag != "" {
			out = append(out, tag)
		}
	}

	return strings.Join(out, ",")
}

func (s *Service) GetTicket(ctx context.Context, id int) (*models.Ticket, error) {
	if id <= 0 {
		return nil, fmt.Errorf("zammad: invalid ticket ID: %d", id)
	}

	var ticket models.Ticket

	err := s.transport.Do(
		ctx,
		http.MethodGet,
		"/tickets/"+strconv.Itoa(id)+"?expand=true",
		nil,
		&ticket,
	)
	if err != nil {
		return nil, err
	}

	return s.withTags(ctx, &ticket)
}

func (s *Service) GetTicketByNumber(ctx context.Context, number string) (*models.Ticket, error) {
	number = strings.TrimSpace(number)
	if number == "" {
		return nil, fmt.Errorf("zammad: ticket number is required")
	}

	values := url.Values{
		"query":  {"number:" + number},
		"expand": {"true"},
	}

	var tickets []models.Ticket

	err := s.transport.Do(
		ctx,
		http.MethodGet,
		"/tickets/search?"+values.Encode(),
		nil,
		&tickets,
	)
	if err != nil {
		return nil, err
	}

	for _, ticket := range tickets {
		if ticket.Number == number {
			return s.withTags(ctx, &ticket)
		}
	}

	return nil, fmt.Errorf("zammad: ticket not found: %s", number)
}

// SearchTicketsAll retrieves every ticket matching query (Zammad search
// syntax) by requesting pages until none remain. Total is the number of
// tickets retrieved.
func (s *Service) SearchTicketsAll(ctx context.Context, query string) (*models.TicketList, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, fmt.Errorf("zammad: search query is required")
	}

	tickets, err := pagination.All(
		ctx,
		func(ctx context.Context, opts models.ListOptions) ([]models.Ticket, error) {
			return s.searchPage(ctx, query, opts)
		},
		func(t models.Ticket) int { return t.ID },
	)
	if err != nil {
		return nil, err
	}

	return &models.TicketList{
		Total:   len(tickets),
		Tickets: tickets,
	}, nil
}

// SearchTickets retrieves every ticket matching query.
//
// Deprecated: Use SearchTicketsAll, or SearchTicketsPage to fetch a single page.
func (s *Service) SearchTickets(ctx context.Context, query string) (*models.TicketList, error) {
	return s.SearchTicketsAll(ctx, query)
}

// SearchTicketsPage retrieves a single page of tickets matching query.
// Total is the number of tickets in this page, not across all pages.
func (s *Service) SearchTicketsPage(ctx context.Context, query string, opts models.ListOptions) (*models.TicketList, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, fmt.Errorf("zammad: search query is required")
	}

	opts, err := pagination.Normalize(opts)
	if err != nil {
		return nil, err
	}

	tickets, err := s.searchPage(ctx, query, opts)
	if err != nil {
		return nil, err
	}

	return &models.TicketList{
		Total:   len(tickets),
		Tickets: tickets,
	}, nil
}

// CountTickets returns the number of tickets matching query (Zammad search
// syntax) using a single request, without retrieving the tickets.
func (s *Service) CountTickets(ctx context.Context, query string) (int, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return 0, fmt.Errorf("zammad: search query is required")
	}

	// full=true makes Zammad report total_count; per_page=1 keeps the
	// response small since only the count is needed.
	path := pagination.Path("/tickets/search", models.ListOptions{Page: 1, PerPage: 1}, url.Values{
		"query": {query},
		"full":  {"true"},
	})

	var result struct {
		TotalCount *int `json:"total_count"`
	}

	if err := s.transport.Do(ctx, http.MethodGet, path, nil, &result); err != nil {
		return 0, err
	}

	if result.TotalCount == nil {
		return 0, fmt.Errorf("zammad: search response has no total_count")
	}

	return *result.TotalCount, nil
}

func (s *Service) searchPage(ctx context.Context, query string, opts models.ListOptions) ([]models.Ticket, error) {
	path := pagination.Path("/tickets/search", opts, url.Values{
		"query":  {query},
		"expand": {"true"},
	})

	var tickets []models.Ticket
	if err := s.transport.Do(ctx, http.MethodGet, path, nil, &tickets); err != nil {
		return nil, err
	}

	return tickets, nil
}

func (s *Service) UpdateTicket(ctx context.Context, ticketID int, payload map[string]any) error {
	if ticketID <= 0 {
		return fmt.Errorf("zammad: invalid ticket ID: %d", ticketID)
	}

	if len(payload) == 0 {
		return fmt.Errorf("zammad: update payload cannot be empty")
	}

	path := "/tickets/" + strconv.Itoa(ticketID)

	return s.transport.Do(ctx, http.MethodPut, path, payload, nil)
}

func (s *Service) DeleteTicket(ctx context.Context, ticketID int) error {
	if ticketID <= 0 {
		return fmt.Errorf("zammad: invalid ticket ID: %d", ticketID)
	}

	path := "/tickets/" + strconv.Itoa(ticketID)

	return s.transport.Do(ctx, http.MethodDelete, path, nil, nil)
}
