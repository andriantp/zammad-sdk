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

	// List users
	users, err := client.Users.List(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total users: %d\n", users.Total)

	for _, user := range users.Users {
		fmt.Printf(
			"ID: %d, Name: %s %s, Email: %s, Active: %t\n",
			user.ID,
			user.Firstname,
			user.Lastname,
			user.Email,
			user.Active,
		)
	}

	// Get user by ID
	if len(users.Users) == 0 {
		return
	}

	userID := users.Users[0].ID

	user, err := client.Users.Get(ctx, userID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nUser detail:")
	fmt.Printf("ID: %d\n", user.ID)
	fmt.Printf("Login: %s\n", user.Login)
	fmt.Printf("Name: %s %s\n", user.Firstname, user.Lastname)
	fmt.Printf("Email: %s\n", user.Email)
	fmt.Printf("Phone: %s\n", user.Phone)
	fmt.Printf("Active: %t\n", user.Active)
	fmt.Printf("Role IDs: %v\n", user.RoleIDs)
	fmt.Printf("Group IDs: %v\n", user.GroupIDs)
}
