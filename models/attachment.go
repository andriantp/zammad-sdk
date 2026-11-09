package models

type (
	AttachmentInput struct {
		Filename string `json:"filename"`
		Data     string `json:"data"`
		MimeType string `json:"mime-type"`
	}

	Attachment struct {
		ID          int            `json:"id"`
		Filename    string         `json:"filename"`
		Size        string         `json:"size"`
		ContentType string         `json:"content_type,omitempty"`
		Preferences map[string]any `json:"preferences,omitempty"`
	}
)
