package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/D3rise/gator/internal/database"
	"github.com/D3rise/gator/internal/rss"
	"github.com/D3rise/gator/internal/state"
	"os"
	"os/signal"
	"sync"
	"time"
)

func NewAggCommand() Command {
	return Command{
		Name:        "agg",
		Handler:     aggCommandHandler,
		Description: "Aggregate RSS feed and save it to database",
	}
}

func aggCommandHandler(state *state.State, _ ...string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	feedCh := make(chan rss.Feed)
	go fetchFeeds(ctx, &wg, state, feedCh)

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)

	for {
		select {
		case feed := <-feedCh:
			fmt.Println(feed)
		case <-sigCh:
			cancel()

			fmt.Println("Waiting for remaining feeds to be fetched...")
			wg.Wait()
			return nil
		}
	}
}

func fetchFeeds(ctx context.Context, wg *sync.WaitGroup, state *state.State, out chan<- rss.Feed) {
	t := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-t.C:
			go fetchNextFeed(ctx, wg, state, out)
		case <-ctx.Done():
			fmt.Println("Received cancel signal, stopping ticker...")
			t.Stop()
			return
		}
	}
}

func fetchNextFeed(ctx context.Context, wg *sync.WaitGroup, state *state.State, out chan<- rss.Feed) {
	fmt.Println("Fetching next feed")
	time.Sleep(2 * time.Second)
	wg.Add(1)
	defer wg.Done()

	feedToFetch, err := getNextFeedToFetch(state.Queries)
	if err != nil {
		fmt.Printf("Error fetching next feed: %v\n", err)
		return
	}

	requestCtx, cancelRequest := context.WithDeadline(ctx, time.Now().Add(5*time.Second))
	defer cancelRequest()

	feed, err := rss.FetchRSSFeed(requestCtx, *state.HttpClient, feedToFetch.Url)
	if err != nil {
		fmt.Printf("Error fetching feed %s: %v", feedToFetch.Url, err)
		return
	}

	err = state.Queries.SetFeedFetchedAtToNowById(ctx, feedToFetch.ID)
	if err != nil {
		fmt.Printf("error whilst setting last fetched at value on a feed: %v", err)
		return
	}

	out <- *feed
}

func getNextFeedToFetch(queries *database.Queries) (*database.Feed, error) {
	feeds, err := queries.GetFeedListSortedByLastFetchedAt(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return new(database.Feed), nil
		}
		return new(database.Feed), fmt.Errorf("error whilst getting next feed to fetch: %w", err)
	}

	return &feeds[0], nil
}
