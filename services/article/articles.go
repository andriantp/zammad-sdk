package article

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path"
	"strconv"
	"strings"
	"unicode"

	"github.com/andriantp/zammad-sdk/models"
)

type Transport interface {
	Do(ctx context.Context, method, path string, body, result any) error
	DoRaw(ctx context.Context, method, path string) ([]byte, error)
}

// StreamTransport is implemented by transports that can return a response
// body unread. It is separate from Transport so existing Transport
// implementations keep compiling; OpenAttachment falls back to DoRaw for
// transports that do not implement it.
type StreamTransport interface {
	DoStream(ctx context.Context, method, path string) (*http.Response, error)
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

// DownloadAttachment downloads an attachment into memory. For large files,
// or to learn the original filename and content type, use OpenAttachment.
func (s *Service) DownloadAttachment(ctx context.Context, ticketID, articleID, attachmentID int) ([]byte, error) {
	p, err := attachmentPath(ticketID, articleID, attachmentID)
	if err != nil {
		return nil, err
	}

	return s.transport.DoRaw(ctx, http.MethodGet, p)
}

// OpenAttachment opens an attachment for streaming, along with its filename,
// content type and size. The caller must close the returned Body.
//
// The client timeout (WithTimeout, 30s by default) also covers reading Body;
// raise it or use a context deadline for large files on slow connections.
func (s *Service) OpenAttachment(ctx context.Context, ticketID, articleID, attachmentID int) (*models.AttachmentFile, error) {
	p, err := attachmentPath(ticketID, articleID, attachmentID)
	if err != nil {
		return nil, err
	}

	st, ok := s.transport.(StreamTransport)
	if !ok {
		data, err := s.transport.DoRaw(ctx, http.MethodGet, p)
		if err != nil {
			return nil, err
		}

		return &models.AttachmentFile{
			Body: io.NopCloser(bytes.NewReader(data)),
			Size: int64(len(data)),
		}, nil
	}

	resp, err := st.DoStream(ctx, http.MethodGet, p)
	if err != nil {
		return nil, err
	}

	return &models.AttachmentFile{
		Body:        resp.Body,
		Filename:    filenameFromHeader(resp.Header.Get("Content-Disposition")),
		ContentType: resp.Header.Get("Content-Type"),
		Size:        resp.ContentLength,
	}, nil
}

func attachmentPath(ticketID, articleID, attachmentID int) (string, error) {
	if ticketID <= 0 || articleID <= 0 || attachmentID <= 0 {
		return "", fmt.Errorf("zammad: invalid attachment identifiers")
	}

	return fmt.Sprintf("/ticket_attachment/%d/%d/%d", ticketID, articleID, attachmentID), nil
}

// filenameFromHeader extracts a safe base filename from a Content-Disposition
// header. Directory components are stripped so a hostile name such as
// "../../.bashrc" cannot escape the caller's target directory; names that are
// empty, dot-only or contain control characters yield "".
func filenameFromHeader(header string) string {
	_, params, err := mime.ParseMediaType(header)
	if err != nil {
		return ""
	}

	name := strings.ReplaceAll(params["filename"], "\\", "/")
	name = path.Base(name)

	if name == "." || name == ".." || name == "/" || strings.Trim(name, ".") == "" {
		return ""
	}

	for _, r := range name {
		if unicode.IsControl(r) {
			return ""
		}
	}

	return name
}
