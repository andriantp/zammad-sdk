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

	// Get ticket states
	states, err := client.Tickets.GetStates(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total states: %d\n", states.Total)

	for _, state := range states.States {
		fmt.Printf(
			"ID: %d, Name: %s, Active: %t\n",
			state.ID,
			state.Name,
			state.Active,
		)
	}

	// Get state by ID
	if len(states.States) == 0 {
		return
	}

	stateID := states.States[0].ID

	state, err := client.Tickets.GetState(ctx, stateID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nState detail:")
	fmt.Printf("ID: %d\n", state.ID)
	fmt.Printf("Name: %s\n", state.Name)
	fmt.Printf("State Type ID: %d\n", state.StateTypeID)
	fmt.Printf("Default Create: %t\n", state.DefaultCreate)
	fmt.Printf("Default Follow Up: %t\n", state.DefaultFollowUp)
	fmt.Printf("Default Close: %t\n", state.DefaultClose)
	fmt.Printf("Ignore Escalation: %t\n", state.IgnoreEscalation)
	fmt.Printf("Active: %t\n", state.Active)
}
