package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/andriantp/zammad-sdk"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run ./examples/ticket/update/ <ticket_id>")
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

	payload := map[string]any{
		"title":       "Updated ticket title",
		"priority_id": 3,
		"state_id":    2,
		// "tags":        "sdk-zammad",
	}

	err = client.Tickets.UpdateTicket(context.Background(), ticketID, payload)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Ticket %d updated successfully", ticketID)
}
