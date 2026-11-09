package ticket

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/andriantp/zammad-sdk/models"
)

func (s *Service) GetStates(ctx context.Context) (*models.StateList, error) {
	var states []models.State

	err := s.transport.Do(
		ctx,
		http.MethodGet,
		"/ticket_states",
		nil,
		&states,
	)
	if err != nil {
		return nil, err
	}

	return &models.StateList{
		Total:  len(states),
		States: states,
	}, nil
}

func (s *Service) GetState(ctx context.Context, id int) (*models.State, error) {
	if id <= 0 {
		return nil, fmt.Errorf("zammad: invalid ticket state ID: %d", id)
	}

	var state models.State

	err := s.transport.Do(
		ctx,
		http.MethodGet,
		"/ticket_states/"+strconv.Itoa(id),
		nil,
		&state,
	)
	if err != nil {
		return nil, err
	}

	return &state, nil
}
