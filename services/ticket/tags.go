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

func (s *Service) GetTags(ctx context.Context) (*models.TagList, error) {
	var tags []models.Tag

	err := s.transport.Do(
		ctx,
		http.MethodGet,
		"/tag_list",
		nil,
		&tags,
	)
	if err != nil {
		return nil, err
	}

	return &models.TagList{
		Total: len(tags),
		Tags:  tags,
	}, nil
}

// GetTicketTags retrieves the tags assigned to a ticket.
//
// Zammad does not include tags in ticket responses, so GetTicket and
// GetTicketByNumber call this to populate Ticket.Tags. Search results do not
// include tags; call this per ticket when needed.
func (s *Service) GetTicketTags(ctx context.Context, ticketID int) ([]string, error) {
	if ticketID <= 0 {
		return nil, fmt.Errorf("zammad: invalid ticket ID: %d", ticketID)
	}

	values := url.Values{
		"object": {"Ticket"},
		"o_id":   {strconv.Itoa(ticketID)},
	}

	var result struct {
		Tags []string `json:"tags"`
	}

	if err := s.transport.Do(ctx, http.MethodGet, "/tags?"+values.Encode(), nil, &result); err != nil {
		return nil, err
	}

	if result.Tags == nil {
		return []string{}, nil
	}

	return result.Tags, nil
}

// withTags populates ticket.Tags from the tags endpoint.
func (s *Service) withTags(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error) {
	tags, err := s.GetTicketTags(ctx, ticket.ID)
	if err != nil {
		return nil, fmt.Errorf("get tags of ticket %d: %w", ticket.ID, err)
	}

	ticket.Tags = strings.Join(tags, ",")

	return ticket, nil
}
