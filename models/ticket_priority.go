package models

import "time"

type (
	Priorities struct {
		Total      int        `json:"total"`
		Priorities []Priority `json:"priorities"`
	}

	Priority struct {
		ID            int       `json:"id"`
		Name          string    `json:"name"`
		DefaultCreate bool      `json:"default_create"`
		UIIcon        *string   `json:"ui_icon"`
		UIColor       *string   `json:"ui_color"`
		Note          *string   `json:"note"`
		Active        bool      `json:"active"`
		UpdatedByID   int       `json:"updated_by_id"`
		CreatedByID   int       `json:"created_by_id"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
	}
)
