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

func (s *Service) CreateTicket(ctx context.Context, input models.TicketInput) (*models.Ticket, error) {
	var ticket models.Ticket

	err := s.transport.Do(
		ctx,
		http.MethodPost,
		"/tickets",
		input,
		&ticket,
	)
	if err != nil {
		return nil, err
	}

	return &ticket, nil
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
