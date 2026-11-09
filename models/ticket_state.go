package models

import "time"

type (
	StateList struct {
		Total  int     `json:"total"`
		States []State `json:"states"`
	}

	State struct {
		ID               int       `json:"id"`
		StateTypeID      int       `json:"state_type_id"`
		Name             string    `json:"name"`
		NextStateID      *int      `json:"next_state_id"`
		IgnoreEscalation bool      `json:"ignore_escalation"`
		DefaultCreate    bool      `json:"default_create"`
		DefaultFollowUp  bool      `json:"default_follow_up"`
		DefaultClose     bool      `json:"default_close"`
		Note             *string   `json:"note"`
		Active           bool      `json:"active"`
		UpdatedByID      int       `json:"updated_by_id"`
		CreatedByID      int       `json:"created_by_id"`
		CreatedAt        time.Time `json:"created_at"`
		UpdatedAt        time.Time `json:"updated_at"`
	}
)
