package commands

import (
	"context"
	"fmt"
	"github.com/D3rise/gator/internal/database"
	"github.com/D3rise/gator/internal/state"
)

func NewFollowCommand() Command {
	return Command{
		Name:                   "follow",
		Args:                   []string{"feedUrl"},
		Handler:                followCommandHandler,
		RequiresAuthentication: true,
		Description:            "Follow an RSS feed",
	}
}

func followCommandHandler(state *state.State, args ...string) error {
	feedUrl := args[0]
	feed, err := state.Queries.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("error while trying to check feed existense: %w", err)
	}

	currentUser, err := state.Queries.GetUserByName(context.Background(), state.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting current user: %w", err)
	}

	alreadyFollowing, err := state.Queries.CheckFeedFollowExistence(
		context.Background(),
		database.CheckFeedFollowExistenceParams{
			UserID: currentUser.ID,
			FeedID: feed.Feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("error checking feed follow: %w", err)
	}

	if alreadyFollowing {
		return fmt.Errorf("you are already following this feed")
	}

	_, err = state.Queries.CreateNewFeedFollow(
		context.Background(),
		database.CreateNewFeedFollowParams{FeedID: feed.Feed.ID, UserID: currentUser.ID},
	)
	if err != nil {
		return fmt.Errorf("error creating new feed follow: %w", err)
	}

	fmt.Printf("Successfully following feed %s\n", feed.Feed.Name)

	return nil
}
