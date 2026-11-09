package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	zammad "github.com/andriantp/zammad-sdk"
	"github.com/andriantp/zammad-sdk/models"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run ./examples/ticket/create/with_attachment [attachment...]")
	}

	client, err := zammad.NewClient(
		zammad.WithBaseURL(os.Getenv("URL")),
		zammad.WithToken(os.Getenv("TOKEN")),
	)
	if err != nil {
		log.Fatal(err)
	}

	attachments, err := newAttachments(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	input := models.TicketInput{
		Title:          "Dashboard image issue",
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
			Subject:     "Dashboard screenshot",
			Body:        "Please check the attached screenshot.",
			Type:        "note",
			Sender:      "Customer",
			Attachments: attachments,
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
