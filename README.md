
# Zammad Go SDK

A Go SDK for interacting with the [Zammad](https://zammad.org/) REST API.

` sdk-zammad ` provides a simple and structured interface to manage tickets, articles, users, groups, and organizations.

## Features

- Configurable HTTP client and request timeout.
- API authentication using Zammad access tokens.
- Context-aware HTTP requests.
- Centralized HTTP transport and API error handling.
- Ticket management, including attachments.
- Article retrieval, replies, and attachment downloads.
- Resource-specific services and models.

## Installation

```bash
go get github.com/andriantp/zammad-sdk
```

## Configuration

Create a Zammad client using your instance URL and API token.

```go
package main

import (
    "log"

    zammad "github.com/andriantp/zammad-sdk"
)

func main() {
    client, err := zammad.NewClient(
        zammad.WithBaseURL("https://zammad.example.com"),
        zammad.WithToken("your-api-token"),
    )
    if err != nil {
        log.Fatal(err)
    }

    _ = client
}
```

### Available Options

| Option | Description |
|---|---|
| `WithBaseURL(url)` | Configure Zammad instance URL |
| `WithToken(token)` | Configure API access token |
| `WithHTTPClient(client)` | Use a custom HTTP client |
| `WithTimeout(duration)` | Configure the default HTTP client timeout |

## Usage

All service methods accept `context.Context` to support request cancellation and deadlines.

### Ticket

The ticket service supports creating, retrieving, searching, updating, and deleting tickets.

#### Create Ticket

```go
input := models.TicketInput{
    Title:       "Login issue",
    Customer:    "customer@example.com",
    Group:       "Support",
    PriorityID:  2,
    StateID:     1,
    Tags:        "bug,urgent",
    Article: models.ArticleInput{
        Subject: "Unable to login",
        Body:    "I cannot access my account.",
        Type:    "note",
        Sender:  "Customer",
    },
}

ticket, err := client.Tickets.CreateTicket(ctx, input)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Ticket created: %s\n", ticket.Number)
```

#### Get Ticket by ID

```go
ticket, err := client.Tickets.GetTicket(ctx, 123)
if err != nil {
    log.Fatal(err)
}

fmt.Println(ticket.Title)
```

#### Get Ticket by Number

```go
ticket, err := client.Tickets.GetTicketByNumber(ctx, "1000123")
if err != nil {
    log.Fatal(err)
}

fmt.Println(ticket.Title)
```

#### Search Tickets

Search tickets using the Zammad search query syntax.

```go
result, err := client.Tickets.SearchTickets(ctx, "title:login")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total tickets: %d\n", result.Total)

for _, ticket := range result.Tickets {
    fmt.Printf("%s: %s\n", ticket.Number, ticket.Title)
}
```

#### Update Ticket

Update ticket attributes using a flexible payload.

```go
payload := map[string]any{
    "title":       "Updated ticket title",
    "priority_id": 3,
    "state_id":    2,
    "tags":        "urgent,bug",
}

err := client.Tickets.UpdateTicket(ctx, 123, payload)
if err != nil {
    log.Fatal(err)
}
```

#### Delete Ticket

```go
err := client.Tickets.DeleteTicket(ctx, 123)
if err != nil {
    log.Fatal(err)
}
```

### Ticket Priorities

Retrieve available ticket priorities.

```go
result, err := client.Tickets.GetPriorities(ctx)
if err != nil {
    log.Fatal(err)
}

for _, priority := range result.Priorities {
    fmt.Printf("%d: %s\n", priority.ID, priority.Name)
}
```

Retrieve a priority by ID.

```go
priority, err := client.Tickets.GetPriority(ctx, 2)
if err != nil {
    log.Fatal(err)
}

fmt.Println(priority.Name)
```

### Ticket States

Retrieve available ticket states.

```go
result, err := client.Tickets.GetStates(ctx)
if err != nil {
    log.Fatal(err)
}

for _, state := range result.States {
    fmt.Printf("%d: %s\n", state.ID, state.Name)
}
```

Retrieve a state by ID.

```go
state, err := client.Tickets.GetState(ctx, 1)
if err != nil {
    log.Fatal(err)
}

fmt.Println(state.Name)
```

### Ticket Tags

Retrieve available ticket tags.

```go
result, err := client.Tickets.GetTags(ctx)
if err != nil {
    log.Fatal(err)
}

for _, tag := range result.Tags {
    fmt.Printf("%s: %d tickets\n", tag.Name, tag.Count)
}
```

### Articles

#### Get Articles

Retrieve articles associated with a ticket.

```go
result, err := client.Articles.GetArticles(ctx, 123)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total articles: %d\n", result.Total)

for _, article := range result.Articles {
    fmt.Printf("Article ID: %d\n", article.ID)
    fmt.Printf("Subject: %s\n", article.Subject)
    fmt.Printf("Body: %s\n", article.Body)
}
```

#### Reply to Ticket

Create a new article on an existing ticket.

```go
input := models.ArticleInput{
    Subject:  "Re: Login issue",
    Body:     "We are investigating the issue.",
    Type:     "note",
    Sender:   "Agent",
    Internal: false,
}

article, err := client.Articles.Reply(ctx, 123, input)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Article created: %d\n", article.ID)
```

#### Download Attachment

Download an article attachment as binary data.

```go
data, err := client.Articles.DownloadAttachment(
    ctx,
    123, // Ticket ID
    456, // Article ID
    789, // Attachment ID
)
if err != nil {
    log.Fatal(err)
}

err = os.WriteFile("attachment.png", data, 0644)
if err != nil {
    log.Fatal(err)
}
```

### Groups

Retrieve all groups.

```go
result, err := client.Groups.List(ctx)
if err != nil {
    log.Fatal(err)
}

for _, group := range result.Groups {
    fmt.Printf("%d: %s\n", group.ID, group.Name)
}
```

Retrieve a group by ID.

```go
group, err := client.Groups.Get(ctx, 1)
if err != nil {
    log.Fatal(err)
}

fmt.Println(group.Name)
```

### Organizations

Retrieve all organizations.

```go
result, err := client.Organizations.List(ctx)
if err != nil {
    log.Fatal(err)
}

for _, organization := range result.Organizations {
    fmt.Printf("%d: %s\n", organization.ID, organization.Name)
}
```

Retrieve an organization by ID.

```go
organization, err := client.Organizations.Get(ctx, 1)
if err != nil {
    log.Fatal(err)
}

fmt.Println(organization.Name)
```

### Users

Retrieve all users.

```go
result, err := client.Users.List(ctx)
if err != nil {
    log.Fatal(err)
}

for _, user := range result.Users {
    fmt.Printf("%d: %s %s\n", user.ID, user.Firstname, user.Lastname)
}
```

Retrieve a user by ID.

```go
user, err := client.Users.Get(ctx, 1)
if err != nil {
    log.Fatal(err)
}

fmt.Println(user.Email)
```

## Examples

Runnable examples are available in the `examples/` directory.

```text
examples/
├── article/
│   ├── download_attachment/
│   ├── get/
│   └── reply/
│       ├── with_attachment/
│       └── without_attachment/
├── group/
├── organization/
├── ticket/
│   ├── create/
│   │   ├── with_attachment/
│   │   └── without_attachment/
│   ├── delete/
│   ├── get/
│   │   ├── by_id/
│   │   └── by_number/
│   ├── priority/
│   ├── search/
│   ├── state/
│   ├── tag/
│   └── update/
└── user/
```

Set the following environment variables before running examples:

```bash
export URL="https://zammad.example.com"
export TOKEN="your-api-token"
```

Example:

```bash
go run ./examples/ticket/get/by_id/ 123
```

For more examples, explore the corresponding directories.

## Error Handling

The SDK exposes `APIError` as an alias of the internal transport API error.

```go
ticket, err := client.Tickets.GetTicket(ctx, 123)
if err != nil {
    var apiErr *zammad.APIError
    if errors.As(err, &apiErr) {
        log.Printf("Zammad API error: %v", apiErr)
        return
    }

    log.Fatal(err)
}
```

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests without using the Go test cache:

```bash
go test -count=1 ./...
```

Run static analysis:

```bash
go vet ./...
```

Run tests with the race detector:

```bash
go test -race ./...
```

## Project Structure

```text
.
├── client.go
├── config.go
├── errors.go
├── options.go
├── models/
├── services/
│   ├── article/
│   ├── group/
│   ├── organization/
│   ├── ticket/
│   └── user/
├── internal/
│   └── transport/
└── examples/
```

## License

This project is licensed under the MIT License.
```
