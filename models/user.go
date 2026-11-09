package models

import "time"

type (
	UserList struct {
		Total int    `json:"total"`
		Users []User `json:"users"`
	}

	User struct {
		ID                       int                 `json:"id"`
		OrganizationID           *int                `json:"organization_id"`
		Login                    string              `json:"login"`
		Firstname                string              `json:"firstname"`
		Lastname                 string              `json:"lastname"`
		Email                    string              `json:"email"`
		Image                    *string             `json:"image"`
		ImageSource              *string             `json:"image_source"`
		Web                      string              `json:"web"`
		Phone                    string              `json:"phone"`
		Fax                      string              `json:"fax"`
		Mobile                   string              `json:"mobile"`
		Department               *string             `json:"department"`
		Street                   string              `json:"street"`
		Zip                      string              `json:"zip"`
		City                     string              `json:"city"`
		Country                  string              `json:"country"`
		Address                  *string             `json:"address"`
		VIP                      bool                `json:"vip"`
		Verified                 bool                `json:"verified"`
		Active                   bool                `json:"active"`
		Note                     string              `json:"note"`
		LastLogin                *time.Time          `json:"last_login"`
		Source                   *string             `json:"source"`
		LoginFailed              int                 `json:"login_failed"`
		OutOfOffice              bool                `json:"out_of_office"`
		OutOfOfficeStartAt       *time.Time          `json:"out_of_office_start_at"`
		OutOfOfficeEndAt         *time.Time          `json:"out_of_office_end_at"`
		OutOfOfficeReplacementID *int                `json:"out_of_office_replacement_id"`
		Preferences              map[string]any      `json:"preferences"`
		UpdatedByID              int                 `json:"updated_by_id"`
		CreatedByID              int                 `json:"created_by_id"`
		CreatedAt                time.Time           `json:"created_at"`
		UpdatedAt                time.Time           `json:"updated_at"`
		RoleIDs                  []int               `json:"role_ids"`
		TwoFactorPreferenceIDs   []int               `json:"two_factor_preference_ids"`
		OrganizationIDs          []int               `json:"organization_ids"`
		AuthorizationIDs         []int               `json:"authorization_ids"`
		OverviewSortingIDs       []int               `json:"overview_sorting_ids"`
		GroupIDs                 map[string][]string `json:"group_ids"`
	}
)
