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

	// List organizations
	organizations, err := client.Organizations.List(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total organizations: %d\n", organizations.Total)

	for _, organization := range organizations.Organizations {
		fmt.Printf(
			"ID: %d, Name: %s, Active: %t\n",
			organization.ID,
			organization.Name,
			organization.Active,
		)
	}

	orgID := organizations.Organizations[0].ID

	// Get organization
	organization, err := client.Organizations.Get(ctx, orgID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nOrganization Detail:\n")
	fmt.Printf("ID: %d\n", organization.ID)
	fmt.Printf("Name: %s\n", organization.Name)
	fmt.Printf("Domain: %s\n", organization.Domain)
	fmt.Printf("Active: %t\n", organization.Active)
}
