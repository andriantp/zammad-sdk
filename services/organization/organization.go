package organization

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

// ListAll retrieves every organization by requesting pages until none remain.
// Total is the number of organizations retrieved.
func (s *Service) ListAll(ctx context.Context) (*models.OrganizationList, error) {
	organizations, err := pagination.All(ctx, s.fetchPage, func(item models.Organization) int { return item.ID })
	if err != nil {
		return nil, err
	}

	return &models.OrganizationList{
		Total:         len(organizations),
		Organizations: organizations,
	}, nil
}

// List retrieves every organization.
//
// Deprecated: Use ListAll, or ListPage to fetch a single page.
func (s *Service) List(ctx context.Context) (*models.OrganizationList, error) {
	return s.ListAll(ctx)
}

// ListPage retrieves a single page of organizations.
// Total is the number of organizations in this page, not across all pages.
func (s *Service) ListPage(ctx context.Context, opts models.ListOptions) (*models.OrganizationList, error) {
	opts, err := pagination.Normalize(opts)
	if err != nil {
		return nil, err
	}

	organizations, err := s.fetchPage(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &models.OrganizationList{
		Total:         len(organizations),
		Organizations: organizations,
	}, nil
}

func (s *Service) fetchPage(ctx context.Context, opts models.ListOptions) ([]models.Organization, error) {
	var organizations []models.Organization

	err := s.transport.Do(
		ctx,
		http.MethodGet,
		pagination.Path("/organizations", opts, nil),
		nil,
		&organizations,
	)
	if err != nil {
		return nil, err
	}

	return organizations, nil
}

func (s *Service) Get(ctx context.Context, id int) (*models.Organization, error) {
	if id <= 0 {
		return nil, fmt.Errorf("zammad: invalid organization ID: %d", id)
	}

	var organization models.Organization

	err := s.transport.Do(
		ctx,
		http.MethodGet,
		"/organizations/"+strconv.Itoa(id),
		nil,
		&organization,
	)
	if err != nil {
		return nil, err
	}

	return &organization, nil
}
