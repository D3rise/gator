package commands

import (
	"fmt"
	"github.com/D3rise/gator/internal/rss"
	"github.com/D3rise/gator/internal/state"
)

func NewAggCommand() Command {
	return Command{
		Name:        "agg",
		Args:        []string{"rssFeedUrl"},
		Handler:     aggCommandHandler,
		Description: "Aggregate RSS feed and save it to database",
	}
}

func aggCommandHandler(state *state.State, args ...string) error {
	httpClient := state.HttpClient
	rssFeedUrl := args[0]

	feed, err := rss.FetchRSSFeed(*httpClient, rssFeedUrl)
	if err != nil {
		return fmt.Errorf("error fetching rss feed: %w", err)
	}

	fmt.Println(feed)

	return nil
}
