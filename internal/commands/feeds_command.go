package commands

import (
	"context"
	"fmt"
	"github.com/D3rise/gator/internal/state"
)

func NewFeedsCommand() Command {
	return Command{
		Name:        "feeds",
		Handler:     feedsCommandHandler,
		Description: "List all added feeds",
	}
}

func feedsCommandHandler(state *state.State, _ ...string) error {
	feedList, err := state.Queries.GetFeedListSortedByCreation(context.Background())
	if err != nil {
		return fmt.Errorf("error whilst getting feed list: %w", err)
	}

	if len(feedList) == 0 {
		fmt.Println("There are currently no feeds added to the system")
		return nil
	}

	fmt.Print("These feeds are currently available:\n\n")

	for _, feed := range feedList {
		fmt.Printf(" * %s - %s (author: %s)\n", feed.Feed.Name, feed.Feed.Url, feed.User.Name)
	}

	return nil
}
