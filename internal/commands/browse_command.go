package commands

import (
	"context"
	"fmt"
	"github.com/D3rise/gator/internal/database"
	"github.com/D3rise/gator/internal/state"
	"strconv"
	"time"
)

var (
	pageSize = 5
)

func NewBrowseCommand() Command {
	return Command{
		Name:        "browse",
		Args:        []string{"page"},
		Handler:     browseCommandHandler,
		Description: "Browse your followed feeds",
	}
}

func browseCommandHandler(state *state.State, args ...string) error {
	page, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	postCount, err := state.Queries.GetPostsCountByUserFeedFollows(context.Background(), state.Config.CurrentUserId)
	if err != nil {
		return fmt.Errorf("error whilst retrieving posts count: %w", err)
	}

	if postCount == 0 {
		fmt.Println("You don't have new posts to browse.")
		return nil
	}

	posts, err := state.Queries.GetPostsByUserFeedFollowsPaginated(
		context.Background(),
		database.GetPostsByUserFeedFollowsPaginatedParams{
			ID:     state.Config.CurrentUserId,
			Limit:  int32(pageSize),
			Offset: int32(pageSize*page - 1),
		},
	)

	for i, post := range posts {
		fmt.Printf("%d. [%s]: "+
			"\n\tTitle: %s"+
			"\n\tURL: %s"+
			"\n\tPublished at: %v"+
			"\n",
			i+1,
			post.Feed.Name,
			post.Post.Title,
			post.Post.Url,
			post.Post.PublishedAt.Format(time.RFC822),
		)
	}

	pageCount := postCount / 5

	fmt.Printf("\nPage %d of %d\n", page, pageCount)

	return nil
}
