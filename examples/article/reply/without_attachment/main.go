package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	zammad "github.com/andriantp/zammad-sdk"
	"github.com/andriantp/zammad-sdk/models"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run ./examples/article/reply/without_attachment/ <ticket_id>")
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

	input := models.ArticleInput{
		Subject:  "Re: Ticket update",
		Body:     "We have received your report and are investigating the issue.",
		Type:     "note",
		Sender:   "Agent",
		Internal: false,
	}

	article, err := client.Articles.Reply(context.Background(), ticketID, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Article created successfully\n")
	fmt.Printf("Article ID: %d\n", article.ID)
	fmt.Printf("Ticket ID: %d\n", article.TicketID)
	fmt.Printf("Subject: %s\n", article.Subject)
}
