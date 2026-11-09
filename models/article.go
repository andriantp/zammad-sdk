package models

import (
	"time"
)

type ArticleList struct {
	Total    int       `json:"total"`
	Articles []Article `json:"articles"`
}

type (
	ArticleInput struct {
		Subject     string            `json:"subject"`
		Body        string            `json:"body"`
		Type        string            `json:"type"`
		Sender      string            `json:"sender"`
		Internal    bool              `json:"internal"`
		ContentType string            `json:"content_type,omitempty"`
		From        string            `json:"from,omitempty"`
		To          string            `json:"to,omitempty"`
		Cc          *string           `json:"cc,omitempty"`
		ReplyTo     *string           `json:"reply_to,omitempty"`
		Attachments []AttachmentInput `json:"attachments,omitempty"`
	}

	Article struct {
		ID          int          `json:"id"`
		TicketID    int          `json:"ticket_id"`
		TypeID      int          `json:"type_id"`
		SenderID    int          `json:"sender_id"`
		From        string       `json:"from"`
		To          string       `json:"to"`
		Cc          *string      `json:"cc"`
		ReplyTo     *string      `json:"reply_to"`
		Subject     string       `json:"subject"`
		Body        string       `json:"body"`
		ContentType string       `json:"content_type"`
		Internal    bool         `json:"internal"`
		OriginByID  *int         `json:"origin_by_id"`
		UpdatedByID int          `json:"updated_by_id"`
		CreatedByID int          `json:"created_by_id"`
		CreatedAt   time.Time    `json:"created_at"`
		UpdatedAt   time.Time    `json:"updated_at"`
		Type        string       `json:"type"`
		Sender      string       `json:"sender"`
		CreatedBy   string       `json:"created_by"`
		UpdatedBy   string       `json:"updated_by"`
		Attachments []Attachment `json:"attachments"`
	}
)
