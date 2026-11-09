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

	// List ticket priorities
	priorities, err := client.Tickets.GetPriorities(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total priorities: %d\n", priorities.Total)

	for _, priority := range priorities.Priorities {
		fmt.Printf(
			"ID: %d, Name: %s, Active: %t, Default Create: %t\n",
			priority.ID,
			priority.Name,
			priority.Active,
			priority.DefaultCreate,
		)
	}

	// Get priority by ID
	if len(priorities.Priorities) == 0 {
		return
	}

	priorityID := priorities.Priorities[0].ID

	priority, err := client.Tickets.GetPriority(ctx, priorityID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nPriority detail:")
	fmt.Printf("ID: %d\n", priority.ID)
	fmt.Printf("Name: %s\n", priority.Name)
	fmt.Printf("Active: %t\n", priority.Active)
	fmt.Printf("Default Create: %t\n", priority.DefaultCreate)
	fmt.Printf("UI Icon: %v\n", priority.UIIcon)
	fmt.Printf("UI Color: %v\n", priority.UIColor)
}
