package commands

import (
	"context"
	"fmt"
	"github.com/D3rise/gator/internal/database"
	"github.com/D3rise/gator/internal/rss"
	"github.com/D3rise/gator/internal/state"
	"github.com/google/uuid"
	"os"
	"os/signal"
	"sync"
	"time"
)

func NewAggCommand() Command {
	return Command{
		Name:        "agg",
		Handler:     aggCommandHandler,
		Description: "Aggregate and display RSS feeds, updating their last fetched time",
	}
}

type feedWithRss struct {
	rss rss.Feed
	row database.Feed
}

func aggCommandHandler(state *state.State, _ ...string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	feedCh := make(chan feedWithRss)
	go fetchFeeds(ctx, &wg, state, feedCh)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

	fmt.Println("Starting RSS feed aggregation. Press Ctrl+C to stop.")

	for {
		select {
		case feed := <-feedCh:
			fmt.Printf("\n=== Fetched feed: %s ===\n", feed.rss.Channel.Title)
			fmt.Printf("Link: %s\n", feed.rss.Channel.Link)
			fmt.Printf("Description: %s\n", feed.rss.Channel.Description)
			fmt.Printf("Items: %d\n\n", len(feed.rss.Channel.Item))

			// Print details of each item
			for i, item := range feed.rss.Channel.Item {
				publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
				if err != nil {
					fmt.Printf("Can't parse publish date of item in %s feed: %v\n", feed.rss.Channel.Link, err)
					continue
				}

				err = state.Queries.UpsertPostOnUrl(ctx, database.UpsertPostOnUrlParams{
					Title:       item.Title,
					Url:         item.Link,
					Description: item.Description,
					PublishedAt: publishedAt,
					FeedID:      feed.row.ID,
				})

				if err != nil {
					fmt.Printf("Error while saving post to database: %v\n", err)
					continue
				}

				fmt.Printf(" Saved Item %d: %s\n", i+1, item.Title)
				fmt.Printf("    Link: %s\n", item.Link)
				fmt.Printf("    Published: %s\n", item.PubDate)
			}
		case <-sigCh:
			cancel()

			fmt.Println("\nWaiting for remaining feeds to be fetched...")
			wg.Wait()
			fmt.Println("Aggregation stopped.")
			return nil
		}
	}
}

func fetchFeeds(ctx context.Context, wg *sync.WaitGroup, state *state.State, out chan<- feedWithRss) {
	t := time.NewTicker(5 * time.Second)

	// Limit concurrent fetches to 5 at a time
	semaphore := make(chan struct{}, 5)

	for {
		select {
		case <-t.C:
			select {
			case semaphore <- struct{}{}:
				go func() {
					defer func() { <-semaphore }()
					fetchNextFeed(ctx, wg, state, out)
				}()
			default:
				fmt.Println("Too many concurrent fetches, skipping this tick")
			}
		case <-ctx.Done():
			fmt.Println("Received cancel signal, stopping ticker...")
			t.Stop()
			return
		}
	}
}

func fetchNextFeed(ctx context.Context, wg *sync.WaitGroup, state *state.State, out chan<- feedWithRss) {
	wg.Add(1)
	defer wg.Done()

	fmt.Println("Fetching next feed")

	feedToFetch, err := getNextFeedToFetch(state.Queries)
	if err != nil {
		fmt.Printf("Error fetching next feed: %v\n", err)
		return
	}

	if feedToFetch.ID == uuid.Nil {
		fmt.Println("No feeds to fetch")
		return
	}

	requestCtx, cancelRequest := context.WithDeadline(ctx, time.Now().Add(5*time.Second))
	defer cancelRequest()

	feed, err := rss.FetchRSSFeed(requestCtx, *state.HttpClient, feedToFetch.Url)
	if err != nil {
		fmt.Printf("Error fetching feed %s: %v\n", feedToFetch.Url, err)
		return
	}

	err = state.Queries.SetFeedFetchedAtToNowById(ctx, feedToFetch.ID)
	if err != nil {
		fmt.Printf("Error setting last fetched at value on feed: %v\n", err)
		return
	}

	out <- feedWithRss{
		row: *feedToFetch,
		rss: *feed,
	}
}

func getNextFeedToFetch(queries *database.Queries) (*database.Feed, error) {
	feeds, err := queries.GetFeedListSortedByLastFetchedAt(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting next feed to fetch: %w", err)
	}

	if len(feeds) == 0 {
		return new(database.Feed), nil
	}

	return &feeds[0], nil
}
