package organization

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

func (s *Service) List(ctx context.Context) (*models.OrganizationList, error) {
	var organizations []models.Organization

	err := s.transport.Do(
		ctx,
		http.MethodGet,
		"/organizations",
		nil,
		&organizations,
	)
	if err != nil {
		return nil, err
	}

	return &models.OrganizationList{
		Total:         len(organizations),
		Organizations: organizations,
	}, nil
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
