package main

import (
	"context"
	"log"
	"os"
	"strconv"

	zammad "github.com/andriantp/zammad-sdk"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run ./examples/ticket/delete/ <ticket_id>")
	}

	ticketID, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("invalid ticket ID: %v", err)
	}

	client, err := zammad.NewClient(
		zammad.WithBaseURL(os.Getenv("URL")),
		zammad.WithToken(os.Getenv("TOKEN")),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Tickets.DeleteTicket(context.Background(), ticketID)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Ticket %d deleted successfully", ticketID)
}
