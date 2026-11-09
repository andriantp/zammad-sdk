package transport

import "fmt"

type APIError struct {
	StatusCode int
	Message    string
	Detail     string
}

func (e *APIError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf(
			"zammad API error: status %d: %s (%s)",
			e.StatusCode,
			e.Message,
			e.Detail,
		)
	}

	return fmt.Sprintf(
		"zammad API error: status %d: %s",
		e.StatusCode,
		e.Message,
	)
}
