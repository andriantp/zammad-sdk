package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	zammad "github.com/andriantp/zammad-sdk"
)

func main() {
	if len(os.Args) < 5 {
		log.Fatal("usage: go run ./examples/article/download_attachment/ <ticket_id> <article_id> <attachment_id> <filename>")
	}

	ticketID, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("invalid ticket ID: %v", err)
	}

	articleID, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("invalid article ID: %v", err)
	}

	attachmentID, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatalf("invalid attachment ID: %v", err)
	}

	client, err := zammad.NewClient(
		zammad.WithBaseURL(os.Getenv("URL")),
		zammad.WithToken(os.Getenv("TOKEN")),
	)
	if err != nil {
		log.Fatal(err)
	}

	data, err := client.Articles.DownloadAttachment(
		context.Background(),
		ticketID,
		articleID,
		attachmentID,
	)
	if err != nil {
		log.Fatal(err)
	}

	filename := os.Args[4]

	if err := os.WriteFile(filename, data, 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Attachment downloaded successfully\n")
	fmt.Printf("Filename: %s\n", filename)
	fmt.Printf("Size: %d bytes\n", len(data))
}
