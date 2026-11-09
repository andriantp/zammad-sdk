package main

import (
	"context"
	"fmt"
	"log"
	"os"

	zammad "github.com/andriantp/zammad-sdk"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run ./examples/ticket/get/by_number <ticket-number>")
	}

	ticketNumber := os.Args[1]

	client, err := zammad.NewClient(
		zammad.WithBaseURL(os.Getenv("URL")),
		zammad.WithToken(os.Getenv("TOKEN")),
	)
	if err != nil {
		log.Fatal(err)
	}

	ticket, err := client.Tickets.GetTicketByNumber(
		context.Background(),
		ticketNumber,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ticket detail:")
	fmt.Printf("ID: %d\n", ticket.ID)
	fmt.Printf("Number: %s\n", ticket.Number)
	fmt.Printf("Title: %s\n", ticket.Title)
	fmt.Printf("Customer ID: %d\n", ticket.CustomerID)
	fmt.Printf("Group ID: %d\n", ticket.GroupID)
	fmt.Printf("Priority ID: %d\n", ticket.PriorityID)
	fmt.Printf("State ID: %d\n", ticket.StateID)
}
