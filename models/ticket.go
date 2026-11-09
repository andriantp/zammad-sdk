package models

import "time"

type TicketSearchFilter struct {
	OrganizationID int
	Tag            string
	From           string
	To             string
	PriorityID     int
	StateID        int
}

type (
	TicketList struct {
		Total   int      `json:"total"`
		Tickets []Ticket `json:"tickets"`
	}

	TicketInput struct {
		Title          string       `json:"title"`
		Group          string       `json:"group,omitempty"`
		Customer       string       `json:"customer,omitempty"`
		CustomerID     int          `json:"customer_id,omitempty"`
		OwnerID        int          `json:"owner_id,omitempty"`
		GroupID        int          `json:"group_id,omitempty"`
		PriorityID     int          `json:"priority_id,omitempty"`
		StateID        int          `json:"state_id,omitempty"`
		State          string       `json:"state,omitempty"`
		OrganizationID int          `json:"organization_id,omitempty"`
		Tags           string       `json:"tags,omitempty"`
		Article        ArticleInput `json:"article"`
	}

	Ticket struct {
		ID                    int       `json:"id"`
		Title                 string    `json:"title"`
		Group                 string    `json:"group"`
		OwnerID               int       `json:"owner_id"`
		GroupID               int       `json:"group_id"`
		PriorityID            int       `json:"priority_id"`
		StateID               int       `json:"state_id"`
		State                 string    `json:"state"`
		OrganizationID        int       `json:"organization_id"`
		Number                string    `json:"number"`
		Customer              string    `json:"customer"`
		CustomerID            int       `json:"customer_id"`
		Tags                  string    `json:"tags"`
		LastContactAt         time.Time `json:"last_contact_at"`
		LastContactAgentAt    time.Time `json:"last_contact_agent_at"`
		LastContactCustomerAt time.Time `json:"last_contact_customer_at"`
		CreateArticleTypeID   int       `json:"create_article_type_id"`
		CreateArticleSenderID int       `json:"create_article_sender_id"`
		ArticleCount          int       `json:"article_count"`
		UpdatedByID           int       `json:"updated_by_id"`
		CreatedByID           int       `json:"created_by_id"`
		CreatedAt             time.Time `json:"created_at"`
		UpdatedAt             time.Time `json:"updated_at"`
		Article               Article   `json:"article"`
	}
)
