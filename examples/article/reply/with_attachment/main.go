package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	zammad "github.com/andriantp/zammad-sdk"
	"github.com/andriantp/zammad-sdk/models"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("usage: go run ./examples/article/reply/with_attachment/ <ticket_id> <attachment...>")
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

	attachments, err := newAttachments(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}

	input := models.ArticleInput{
		Subject:     "Re: Ticket update",
		Body:        "Please find the attached files for further investigation.",
		Type:        "note",
		Sender:      "Agent",
		Internal:    false,
		Attachments: attachments,
	}

	article, err := client.Articles.Reply(context.Background(), ticketID, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Article created successfully\n")
	fmt.Printf("Article ID: %d\n", article.ID)
	fmt.Printf("Ticket ID: %d\n", article.TicketID)
	fmt.Printf("Subject: %s\n", article.Subject)
	fmt.Printf("Attachments: %d\n", len(attachments))
}

func newAttachments(paths []string) ([]models.AttachmentInput, error) {
	var attachments []models.AttachmentInput

	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read attachment %q: %w", path, err)
		}

		if len(content) == 0 {
			return nil, fmt.Errorf("attachment %q is empty", path)
		}

		mimeType := http.DetectContentType(content)

		attachments = append(attachments, models.AttachmentInput{
			Filename: filepath.Base(path),
			Data:     base64.StdEncoding.EncodeToString(content),
			MimeType: mimeType,
		})
	}

	return attachments, nil
}
