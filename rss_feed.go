package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/BlessE1/gator/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch feed: %w", err)
	}
	defer res.Body.Close()

	// Parse the RSS feed
	feed, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse feed: %w", err)
	}

	// Convert the RSS feed to a struct
	var rssFeed RSSFeed
	err = xml.Unmarshal(feed, &rssFeed)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal feed: %w", err)
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)

	for idx, _ := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[idx].Title = html.UnescapeString(rssFeed.Channel.Item[idx].Title)
		rssFeed.Channel.Item[idx].Description = html.UnescapeString(rssFeed.Channel.Item[idx].Description)
	}

	return &rssFeed, nil
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		fmt.Println("Feed could not be fetched")
		os.Exit(1)
	}

	err = s.db.MarkFeedFetched(ctx, nextFeed.ID)
	if err != nil {
		fmt.Println("Feed could not be marked as fetched")
		os.Exit(1)
	}

	rssFeed, err := fetchFeed(ctx, nextFeed.Url)
	if err != nil {
		fmt.Println("Feed could not be fetched")
		os.Exit(1)
	}

	for _, item := range rssFeed.Channel.Item {
		publishedAt, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			fmt.Printf("Failed to parse publication date for item '%s': %v\n", item.Title, err)
		}

		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: sql.NullTime{Time: publishedAt, Valid: true},
			FeedID:      nextFeed.ID,
		}

		_, err = s.db.CreatePost(ctx, postParams)
		if err != nil {
			fmt.Printf("Failed to create post: %v\n", err)
			continue
		}
	}

	fmt.Printf("##### %s #####\n", rssFeed.Channel.Title)
	for _, item := range rssFeed.Channel.Item {
		fmt.Println("Title:", item.Title)
		fmt.Println("Link:", item.Link)
		fmt.Println("-------------------------")
	}

	return nil
}
