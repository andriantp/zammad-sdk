package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	zammad "github.com/andriantp/zammad-sdk"
	"github.com/andriantp/zammad-sdk/models"
)

func main() {
	client, err := zammad.NewClient(
		zammad.WithBaseURL(os.Getenv("URL")),
		zammad.WithToken(os.Getenv("TOKEN")),
	)
	if err != nil {
		log.Fatal(err)
	}

	filter := models.TicketSearchFilter{
		OrganizationID: 2,
		Tag:            "sdk-zammad",
		From:           os.Getenv("FROM"),
		To:             os.Getenv("TO"),
		// PriorityID:     2,
		// StateID:        1,
	}

	query := buildSearchQuery(filter)
	if query == "" {
		log.Fatal("at least one search filter is required")
	}

	tickets, err := client.Tickets.SearchTickets(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total tickets: %d\n", tickets.Total)

	for _, ticket := range tickets.Tickets {
		fmt.Println("------------------------------")
		fmt.Printf("ID: %d\n", ticket.ID)
		fmt.Printf("Number: %s\n", ticket.Number)
		fmt.Printf("Title: %s\n", ticket.Title)
		fmt.Printf("Customer: %s\n", ticket.Customer)
		fmt.Printf("Group: %s\n", ticket.Group)
		fmt.Printf("Priority ID: %d\n", ticket.PriorityID)
		fmt.Printf("State ID: %d\n", ticket.StateID)
	}
}

func buildSearchQuery(filter models.TicketSearchFilter) string {
	var queries []string

	if filter.OrganizationID > 0 {
		queries = append(queries, "organization_id:"+strconv.Itoa(filter.OrganizationID))
	}

	if filter.Tag != "" {
		queries = append(queries, "tags:"+url.QueryEscape(filter.Tag))
	}

	if filter.From != "" && filter.To != "" {
		queries = append(
			queries,
			fmt.Sprintf("created_at:[%s TO %s]", filter.From, filter.To),
		)
	}

	if filter.PriorityID > 0 {
		queries = append(queries, "priority_id:"+strconv.Itoa(filter.PriorityID))
	}

	if filter.StateID > 0 {
		queries = append(queries, "state_id:"+strconv.Itoa(filter.StateID))
	}

	return strings.Join(queries, " ")
}
