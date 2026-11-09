package main

import (
	"context"
	"fmt"
	"log"
	"os"

	zammad "github.com/andriantp/zammad-sdk"
)

func main() {
	client, err := zammad.NewClient(
		zammad.WithBaseURL(os.Getenv("URL")),
		zammad.WithToken(os.Getenv("TOKEN")),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Get ticket tags
	tags, err := client.Tickets.GetTags(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total tags: %d\n", tags.Total)

	for _, tag := range tags.Tags {
		fmt.Printf(
			"ID: %d, Name: %s, Count: %d\n",
			tag.ID,
			tag.Name,
			tag.Count,
		)
	}
}
