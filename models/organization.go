package models

import (
	"time"
)

type (
	OrganizationList struct {
		Total         int            `json:"total"`
		Organizations []Organization `json:"organizations"`
	}

	Organization struct {
		ID                 int       `json:"id"`
		Name               string    `json:"name"`
		Shared             bool      `json:"shared"`
		Domain             string    `json:"domain"`
		DomainAssignment   bool      `json:"domain_assignment"`
		Active             bool      `json:"active"`
		VIP                bool      `json:"vip"`
		Note               string    `json:"note"`
		UpdatedByID        int       `json:"updated_by_id"`
		CreatedByID        int       `json:"created_by_id"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		MemberIDs          []int     `json:"member_ids"`
		SecondaryMemberIDs []int     `json:"secondary_member_ids"`
	}
)
