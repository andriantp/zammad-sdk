package user

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

func (s *Service) List(ctx context.Context) (*models.UserList, error) {
	var users []models.User

	err := s.transport.Do(
		ctx,
		http.MethodGet,
		"/users",
		nil,
		&users,
	)
	if err != nil {
		return nil, err
	}

	return &models.UserList{
		Total: len(users),
		Users: users,
	}, nil
}

func (s *Service) Get(ctx context.Context, id int) (*models.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("zammad: invalid user ID: %d", id)
	}

	var user models.User

	err := s.transport.Do(
		ctx,
		http.MethodGet,
		"/users/"+strconv.Itoa(id),
		nil,
		&user,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
