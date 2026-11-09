package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"

	zammad "github.com/andriantp/zammad-sdk"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run ./examples/article/get/ <ticket_id>")
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

	result, err := client.Articles.GetArticles(context.Background(), ticketID)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Total articles: %d", result.Total)

	for _, article := range result.Articles {
		log.Println("--------------------------------")
		log.Printf("ID: %d", article.ID)
		log.Printf("Ticket ID: %d", article.TicketID)
		log.Printf("Subject: %s", article.Subject)
		log.Printf("Type: %s", article.Type)
		log.Printf("Sender: %s", article.Sender)
		log.Printf("Internal: %t", article.Internal)
		log.Printf("Body: %s", article.Body)

		if len(article.Attachments) > 0 {
			log.Println("Attachments:")
			if js, err := json.MarshalIndent(article.Attachments, "", "  "); err == nil {
				log.Println(string(js))
			}
		}
	}
}
