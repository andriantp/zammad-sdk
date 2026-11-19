package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	zammad "github.com/andriantp/zammad-sdk"
)

const usage = `Count tickets matching a Zammad search query.

Usage:
  go run ./examples/ticket/count [flags] [query]

Flags build the query for you; any remaining arguments are added as a raw
query. All parts are combined with AND. Flags must come before the query.

Flags:
  -state  ticket state, comma-separated for any of several (new,open)
  -tag    ticket tag, comma-separated for any of several (bug,urgent)

Examples:
  go run ./examples/ticket/count -state open
  go run ./examples/ticket/count -tag sdk-zammad -state new
  go run ./examples/ticket/count -state "new,pending reminder"
  go run ./examples/ticket/count 'NOT state.name:closed'
  go run ./examples/ticket/count -tag bug 'created_at:>now-7d'
`

func main() {
	state := flag.String("state", "", "")
	tag := flag.String("tag", "", "")

	flag.Usage = func() { fmt.Fprint(os.Stderr, usage) }
	flag.Parse()

	query := buildQuery(*state, *tag, flag.Args())
	if query == "" {
		flag.Usage()
		os.Exit(2)
	}

	client, err := zammad.NewClient(
		zammad.WithBaseURL(os.Getenv("URL")),
		zammad.WithToken(os.Getenv("TOKEN")),
	)
	if err != nil {
		log.Fatal(err)
	}

	count, err := client.Tickets.CountTickets(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Query: %s\n", query)
	fmt.Printf("Total tickets: %d\n", count)
}

// buildQuery combines the flag filters and the raw query with AND.
func buildQuery(state, tag string, raw []string) string {
	var parts []string

	if q := anyOf("state.name", state); q != "" {
		parts = append(parts, q)
	}

	if q := anyOf("tags", tag); q != "" {
		parts = append(parts, q)
	}

	if q := strings.TrimSpace(strings.Join(raw, " ")); q != "" {
		// Parenthesize so an OR in the raw query cannot escape the AND.
		if len(parts) > 0 {
			q = "(" + q + ")"
		}

		parts = append(parts, q)
	}

	return strings.Join(parts, " AND ")
}

// anyOf turns "a,b" into `field:(a OR b)`, quoting values with spaces.
func anyOf(field, csv string) string {
	var values []string

	for _, v := range strings.Split(csv, ",") {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}

		if strings.ContainsAny(v, " \t") {
			v = strconv.Quote(v)
		}

		values = append(values, v)
	}

	switch len(values) {
	case 0:
		return ""
	case 1:
		return field + ":" + values[0]
	default:
		return field + ":(" + strings.Join(values, " OR ") + ")"
	}
}
