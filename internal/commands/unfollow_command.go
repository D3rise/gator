package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/D3rise/gator/internal/database"
	"github.com/D3rise/gator/internal/state"
)

func NewUnfollowCommand() Command {
	return Command{
		Name:        "unfollow",
		Args:        []string{"feedUrl"},
		Handler:     unfollowCommandHandler,
		Description: "Unfollow a feed",
	}
}

func unfollowCommandHandler(state *state.State, args ...string) error {
	feedUrl := args[0]
	feed, err := state.Queries.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("feed with this url does not exist")
		}
		return fmt.Errorf("error while retrieving feed by url: %w", err)
	}

	isSubscribed, err := state.Queries.CheckFeedFollowExistence(context.Background(), database.CheckFeedFollowExistenceParams{
		UserID: state.Config.CurrentUserId,
		FeedID: feed.Feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error while retrieving feeds you follow: %w", err)
	}

	if !isSubscribed {
		return fmt.Errorf("you are not subscribed to this feed")
	}

	_, err = state.Queries.DeleteFeedFollowByUrlAndUserId(context.Background(), database.DeleteFeedFollowByUrlAndUserIdParams{
		Url:    feedUrl,
		UserID: state.Config.CurrentUserId,
	})

	if err != nil {
		return fmt.Errorf("error while deleting the follow: %w", err)
	}

	fmt.Printf("Successfully unfollowed the \"%s\" feed (%s)\n", feed.Feed.Name, feed.Feed.Url)

	return nil
}
