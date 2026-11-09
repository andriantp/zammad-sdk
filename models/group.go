package models

import "time"

type (
	GroupList struct {
		Total  int     `json:"total"`
		Groups []Group `json:"groups"`
	}

	Group struct {
		ID                 int       `json:"id"`
		SignatureID        int       `json:"signature_id"`
		EmailAddressID     *int      `json:"email_address_id"`
		Name               string    `json:"name"`
		NameLast           string    `json:"name_last"`
		ParentID           *int      `json:"parent_id"`
		AssignmentTimeout  *int      `json:"assignment_timeout"`
		FollowUpPossible   string    `json:"follow_up_possible"`
		ReopenTimeInDays   *int      `json:"reopen_time_in_days"`
		FollowUpAssignment bool      `json:"follow_up_assignment"`
		Active             bool      `json:"active"`
		SharedDrafts       bool      `json:"shared_drafts"`
		SummaryGeneration  string    `json:"summary_generation"`
		Note               string    `json:"note"`
		UpdatedByID        int       `json:"updated_by_id"`
		CreatedByID        int       `json:"created_by_id"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		UserIDs            []int     `json:"user_ids"`
	}
)
