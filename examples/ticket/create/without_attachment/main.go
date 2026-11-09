package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	input := models.TicketInput{
		Title:          "Cannot access dashboard",
		Customer:       "customer@example.com",
		CustomerID:     4,
		OwnerID:        1,
		Group:          "Users",
		GroupID:        1,
		PriorityID:     2,
		StateID:        1,
		OrganizationID: 2,
		Tags:           "sdk-zammad",
		Article: models.ArticleInput{
			Subject: "Dashboard access issue",
			Body:    "I cannot access my dashboard.",
			Type:    "note",
			Sender:  "Customer",
		},
	}

	ticket, err := client.Tickets.CreateTicket(context.Background(), input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ticket created successfully")
	fmt.Printf("ID: %d\n", ticket.ID)
	fmt.Printf("Number: %s\n", ticket.Number)
	fmt.Printf("Title: %s\n", ticket.Title)
}
