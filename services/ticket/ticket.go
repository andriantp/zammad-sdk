package ticket

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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
		"/tickets/"+strconv.Itoa(id),
		nil,
		&ticket,
	)
	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (s *Service) GetTicketByNumber(ctx context.Context, number string) (*models.Ticket, error) {
	number = strings.TrimSpace(number)
	if number == "" {
		return nil, fmt.Errorf("zammad: ticket number is required")
	}

	query := url.QueryEscape("number:" + number)

	var tickets []models.Ticket

	err := s.transport.Do(
		ctx,
		http.MethodGet,
		"/tickets/search?query="+query,
		nil,
		&tickets,
	)
	if err != nil {
		return nil, err
	}

	for _, ticket := range tickets {
		if ticket.Number == number {
			return &ticket, nil
		}
	}

	return nil, fmt.Errorf("zammad: ticket not found: %s", number)
}

func (s *Service) SearchTickets(ctx context.Context, query string) (*models.TicketList, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, fmt.Errorf("zammad: search query is required")
	}

	path := "/tickets/search?query=" + url.QueryEscape(query)

	var tickets []models.Ticket
	if err := s.transport.Do(ctx, http.MethodGet, path, nil, &tickets); err != nil {
		return nil, err
	}

	return &models.TicketList{
		Total:   len(tickets),
		Tickets: tickets,
	}, nil
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
