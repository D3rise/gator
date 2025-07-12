package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"
)

type Feed struct {
	Channel struct {
		Title       string     `xml:"title"`
		Link        string     `xml:"link"`
		Description string     `xml:"description"`
		Item        []FeedItem `xml:"item"`
	} `xml:"channel"`
}

type FeedItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchRSSFeed(httpClient http.Client, url string) (*Feed, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", url, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "gator")

	res, err := httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading xml from response: %w", err)
	}

	feed := Feed{}
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling xml: %w", err)
	}

	unescapeFeed(&feed)
	return &feed, nil
}

func unescapeFeed(feed *Feed) {
	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
	}
}
