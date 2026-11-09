package zammad

import (
	"fmt"
	"net/url"

	"github.com/andriantp/zammad-sdk/internal/transport"
	"github.com/andriantp/zammad-sdk/services/article"
	"github.com/andriantp/zammad-sdk/services/group"
	"github.com/andriantp/zammad-sdk/services/organization"
	"github.com/andriantp/zammad-sdk/services/ticket"
	"github.com/andriantp/zammad-sdk/services/user"
)

type (
	AttachmentService struct{ client *Client }
)

type Client struct {
	config    *config
	transport *transport.Client

	Organizations *organization.Service
	Groups        *group.Service
	Users         *user.Service

	Tickets     *ticket.Service
	Articles    *article.Service
	Attachments *AttachmentService
}

func NewClient(opts ...Option) (*Client, error) {
	cfg := newConfig()

	for _, opt := range opts {
		if opt != nil {
			opt(cfg)
		}
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	httpTransport := transport.NewClient(transport.Config{
		BaseURL:    cfg.baseURL,
		Token:      cfg.token,
		HTTPClient: cfg.httpClient,
		Timeout:    cfg.timeout,
	})

	client := &Client{
		config:    cfg,
		transport: httpTransport,
	}

	client.Organizations = organization.New(client.transport)
	client.Groups = group.New(client.transport)
	client.Users = user.New(client.transport)

	client.Tickets = ticket.New(client.transport)
	client.Articles = article.New(client.transport)
	client.Attachments = &AttachmentService{client: client}

	return client, nil
}

func validateConfig(cfg *config) error {
	if cfg.baseURL == "" {
		return fmt.Errorf("zammad: base URL is required")
	}

	u, err := url.ParseRequestURI(cfg.baseURL)
	if err != nil || u.Host == "" {
		return fmt.Errorf("zammad: invalid base URL")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("zammad: base URL must use HTTP or HTTPS")
	}

	if cfg.token == "" {
		return fmt.Errorf("zammad: API token is required")
	}

	return nil
}
