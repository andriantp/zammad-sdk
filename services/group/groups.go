package group

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

// ListAll retrieves every group by requesting pages until none remain.
// Total is the number of groups retrieved.
func (s *Service) ListAll(ctx context.Context) (*models.GroupList, error) {
	groups, err := pagination.All(ctx, s.fetchPage, func(item models.Group) int { return item.ID })
	if err != nil {
		return nil, err
	}

	return &models.GroupList{
		Total:  len(groups),
		Groups: groups,
	}, nil
}

// List retrieves every group.
//
// Deprecated: Use ListAll, or ListPage to fetch a single page.
func (s *Service) List(ctx context.Context) (*models.GroupList, error) {
	return s.ListAll(ctx)
}

// ListPage retrieves a single page of groups.
// Total is the number of groups in this page, not across all pages.
func (s *Service) ListPage(ctx context.Context, opts models.ListOptions) (*models.GroupList, error) {
	opts, err := pagination.Normalize(opts)
	if err != nil {
		return nil, err
	}

	groups, err := s.fetchPage(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &models.GroupList{
		Total:  len(groups),
		Groups: groups,
	}, nil
}

func (s *Service) fetchPage(ctx context.Context, opts models.ListOptions) ([]models.Group, error) {
	var groups []models.Group

	err := s.transport.Do(
		ctx,
		http.MethodGet,
		pagination.Path("/groups", opts, nil),
		nil,
		&groups,
	)
	if err != nil {
		return nil, err
	}

	return groups, nil
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
