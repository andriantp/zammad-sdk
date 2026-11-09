package article

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/andriantp/zammad-sdk/models"
)

type Transport interface {
	Do(ctx context.Context, method, path string, body, result any) error
	DoRaw(ctx context.Context, method, path string) ([]byte, error)
}

type Service struct {
	transport Transport
}

func New(transport Transport) *Service {
	return &Service{transport: transport}
}

func (s *Service) GetArticles(ctx context.Context, ticketID int) (*models.ArticleList, error) {
	if ticketID <= 0 {
		return nil, fmt.Errorf("zammad: invalid ticket ID: %d", ticketID)
	}

	path := "/ticket_articles/by_ticket/" + strconv.Itoa(ticketID)

	var articles []models.Article
	if err := s.transport.Do(ctx, http.MethodGet, path, nil, &articles); err != nil {
		return nil, err
	}

	return &models.ArticleList{
		Total:    len(articles),
		Articles: articles,
	}, nil
}

func (s *Service) Reply(ctx context.Context, ticketID int, input models.ArticleInput) (*models.Article, error) {
	if ticketID <= 0 {
		return nil, fmt.Errorf("zammad: invalid ticket ID: %d", ticketID)
	}

	if strings.TrimSpace(input.Body) == "" {
		return nil, fmt.Errorf("zammad: article body is required")
	}

	payload := struct {
		TicketID int `json:"ticket_id"`
		models.ArticleInput
	}{
		TicketID:     ticketID,
		ArticleInput: input,
	}

	var article models.Article
	if err := s.transport.Do(ctx, http.MethodPost, "/ticket_articles", payload, &article); err != nil {
		return nil, err
	}

	return &article, nil
}

func (s *Service) DownloadAttachment(ctx context.Context, ticketID, articleID, attachmentID int) ([]byte, error) {
	if ticketID <= 0 || articleID <= 0 || attachmentID <= 0 {
		return nil, fmt.Errorf("zammad: invalid attachment identifiers")
	}

	path := fmt.Sprintf(
		"/ticket_attachment/%d/%d/%d",
		ticketID,
		articleID,
		attachmentID,
	)

	return s.transport.DoRaw(ctx, http.MethodGet, path)
}
