package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	zammad "github.com/andriantp/zammad-sdk"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatal("usage: go run ./examples/article/download_attachment/ <ticket_id> <article_id> <attachment_id> [filename]")
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

	// OpenAttachment streams the file instead of loading it into memory and
	// reports the original filename and content type.
	file, err := client.Articles.OpenAttachment(
		context.Background(),
		ticketID,
		articleID,
		attachmentID,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Body.Close()

	// Use the filename given on the command line, else the original one.
	// file.Filename is already reduced to a base name, so it cannot point
	// outside the current directory.
	filename := file.Filename
	if len(os.Args) > 4 {
		filename = os.Args[4]
	}

	if filename == "" {
		filename = fmt.Sprintf("attachment-%d", attachmentID)
	}

	// O_EXCL refuses to overwrite an existing file.
	out, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		log.Fatal(err)
	}

	written, err := io.Copy(out, file.Body)
	if closeErr := out.Close(); err == nil {
		err = closeErr
	}

	if err != nil {
		_ = os.Remove(filename)
		log.Fatal(err)
	}

	fmt.Printf("Attachment downloaded successfully\n")
	fmt.Printf("Filename: %s\n", filename)
	fmt.Printf("Content-Type: %s\n", file.ContentType)
	fmt.Printf("Size: %d bytes\n", written)
}
