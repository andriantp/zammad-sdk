package ticket

import (
	"context"
	"net/http"

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
