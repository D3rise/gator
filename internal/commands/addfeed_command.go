package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/D3rise/gator/internal/database"
	"github.com/D3rise/gator/internal/state"
)

func NewAddFeedCommand() Command {
	return Command{
		Name:        "addfeed",
		Args:        []string{"feedName", "feedUrl"},
		Handler:     newAddFeedCommandHandler,
		Description: "Add new feed to database",
	}
}

func newAddFeedCommandHandler(state *state.State, args ...string) error {
	feedName, feedUrl := args[0], args[1]

	exists, err := state.Queries.CheckFeedExistenceByName(context.Background(), feedName)
	if err != nil {
		return fmt.Errorf("error checking feed existence: %w", err)
	}

	if exists {
		return fmt.Errorf("feed with this name already exists")
	}

	currentUser, err := state.Queries.GetUserByName(context.Background(), state.Config.CurrentUserName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("current username %s is not registered yet", state.Config.CurrentUserName)
		}

		return fmt.Errorf("error while retrieving current user: %w", err)
	}

	feed, err := state.Queries.CreateFeed(context.Background(), database.CreateFeedParams{
		UserID: currentUser.ID,
		Name:   feedName,
		Url:    feedUrl,
	})

	if err != nil {
		return fmt.Errorf("error creating feed: %w", err)
	}

	fmt.Println("Successfully created new feed with name", feed.Name, "and url", feed.Url)

	return nil
}
