package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/acemouty/gator/internal/database"
	"github.com/acemouty/gator/internal/service"
	"github.com/google/uuid"
)

func scrapeFeeds(db *database.Queries, scrapeRound *int) {
	*scrapeRound++
	fmt.Printf("Scrape Round: %v\n\n", *scrapeRound)
	// Get the next feed to fetch.
	feed, err := db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Printf("Failed to get the next feed: %v\n", err)
		return
	}

	// Mark the feed as fetched.
	err = db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Failed to mark feed as fetched: %v\n", err)
		return
	}

	// Fetch the feed using its URL.
	rssFeed, err := service.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldnt collect feed %v: %v\n", feed.Name, err)
		return
	}

	// Iterate over the items and print their titles.
	fmt.Printf("Feed: %s\n", feed.Name)
	for _, item := range rssFeed.Channel.Item {
		//fmt.Printf("- %s\n", item)
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			FeedID:    feed.ID,
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Url:         item.Link,
			PublishedAt: publishedAt,
		})
		if err != nil {
			// TODO: use psql error codes
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}

	log.Printf("Feed %v collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
