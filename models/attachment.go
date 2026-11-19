package models

import "io"

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

// AttachmentFile is a downloaded attachment whose content is streamed.
// The caller must close Body.
type AttachmentFile struct {
	Body io.ReadCloser

	// Filename is the base name from the Content-Disposition header, or ""
	// if the server sent none or an unsafe one. It never contains a path.
	Filename string

	// ContentType is the Content-Type header, e.g. "image/png".
	ContentType string

	// Size is the Content-Length in bytes, or -1 if unknown.
	Size int64
}
