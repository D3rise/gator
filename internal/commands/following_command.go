package commands

import (
	"context"
	"fmt"
	"github.com/D3rise/gator/internal/state"
)

func NewFollowingCommand() Command {
	return Command{
		Name:                   "following",
		Handler:                followingCommandHandler,
		RequiresAuthentication: true,
		Description:            "List all feeds you are following right now",
	}
}

func followingCommandHandler(state *state.State, _ ...string) error {
	currentUser, err := state.Queries.GetUserByName(context.Background(), state.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting current user: %w", err)
	}

	feeds, err := state.Queries.GetFeedFollowListByUserId(context.Background(), currentUser.ID)
	if err != nil {
		return fmt.Errorf("error querying feed follow list: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("You are currently not following any feeds")
		return nil
	}

	fmt.Print("Here is a list of feeds you are following right now: \n\n")
	for _, feed := range feeds {
		fmt.Printf(" * %s - %s\n", feed.Feed.Name, feed.Feed.Url)
	}

	fmt.Println()
	return nil
}
