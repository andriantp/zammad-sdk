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

	// List groups
	groups, err := client.Groups.List(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total groups: %d\n", groups.Total)

	for _, group := range groups.Groups {
		fmt.Printf("ID: %d, Name: %s, Active: %t\n",
			group.ID,
			group.Name,
			group.Active,
		)
	}

	// Get group by ID
	if len(groups.Groups) == 0 {
		return
	}

	groupID := groups.Groups[0].ID

	group, err := client.Groups.Get(ctx, groupID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nGroup detail:\n")
	fmt.Printf("ID: %d\n", group.ID)
	fmt.Printf("Name: %s\n", group.Name)
	fmt.Printf("Active: %t\n", group.Active)
	fmt.Printf("User IDs: %v\n", group.UserIDs)
}
