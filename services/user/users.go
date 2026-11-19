package user

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

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

// ListAll retrieves every user by requesting pages until none remain.
// Total is the number of users retrieved.
func (s *Service) ListAll(ctx context.Context) (*models.UserList, error) {
	users, err := pagination.All(ctx, s.fetchPage, func(item models.User) int { return item.ID })
	if err != nil {
		return nil, err
	}

	return &models.UserList{
		Total: len(users),
		Users: users,
	}, nil
}

// List retrieves every user.
//
// Deprecated: Use ListAll, or ListPage to fetch a single page.
func (s *Service) List(ctx context.Context) (*models.UserList, error) {
	return s.ListAll(ctx)
}

// ListPage retrieves a single page of users.
// Total is the number of users in this page, not across all pages.
func (s *Service) ListPage(ctx context.Context, opts models.ListOptions) (*models.UserList, error) {
	opts, err := pagination.Normalize(opts)
	if err != nil {
		return nil, err
	}

	users, err := s.fetchPage(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &models.UserList{
		Total: len(users),
		Users: users,
	}, nil
}

func (s *Service) fetchPage(ctx context.Context, opts models.ListOptions) ([]models.User, error) {
	var users []models.User

	err := s.transport.Do(
		ctx,
		http.MethodGet,
		pagination.Path("/users", opts, nil),
		nil,
		&users,
	)
	if err != nil {
		return nil, err
	}

	return users, nil
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
