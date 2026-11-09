package ticket

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/andriantp/zammad-sdk/models"
)

func (s *Service) GetPriorities(ctx context.Context) (*models.Priorities, error) {
	var priorities []models.Priority

	err := s.transport.Do(
		ctx,
		http.MethodGet,
		"/ticket_priorities",
		nil,
		&priorities,
	)
	if err != nil {
		return nil, err
	}

	return &models.Priorities{
		Total:      len(priorities),
		Priorities: priorities,
	}, nil
}

func (s *Service) GetPriority(ctx context.Context, id int) (*models.Priority, error) {
	if id <= 0 {
		return nil, fmt.Errorf("zammad: invalid ticket priority ID: %d", id)
	}

	var priority models.Priority

	err := s.transport.Do(
		ctx,
		http.MethodGet,
		"/ticket_priorities/"+strconv.Itoa(id),
		nil,
		&priority,
	)
	if err != nil {
		return nil, err
	}

	return &priority, nil
}
