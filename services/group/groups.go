package group

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

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

func (s *Service) List(ctx context.Context) (*models.GroupList, error) {
	var groups []models.Group

	err := s.transport.Do(
		ctx,
		http.MethodGet,
		"/groups",
		nil,
		&groups,
	)
	if err != nil {
		return nil, err
	}

	return &models.GroupList{
		Total:  len(groups),
		Groups: groups,
	}, nil
}

func (s *Service) Get(ctx context.Context, id int) (*models.Group, error) {
	if id <= 0 {
		return nil, fmt.Errorf("zammad: invalid group ID: %d", id)
	}

	var group models.Group

	err := s.transport.Do(
		ctx,
		http.MethodGet,
		"/groups/"+strconv.Itoa(id),
		nil,
		&group,
	)
	if err != nil {
		return nil, err
	}

	return &group, nil
}
