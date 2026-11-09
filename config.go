package zammad

import (
	"net/http"
	"time"
)

const (
	defaultTimeout = 30 * time.Second
)

type config struct {
	baseURL    string
	token      string
	httpClient *http.Client
	timeout    time.Duration
}

func newConfig() *config {
	return &config{
		timeout: defaultTimeout,
	}
}
