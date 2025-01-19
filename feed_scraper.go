package main

import (
	"context"
	"fmt"
	"log"

	"github.com/acemouty/gator/internal/database"
	"github.com/acemouty/gator/internal/service"
)

func scrapeFeeds(db *database.Queries) {
	for {
		// Get the next feed to fetch.
		feed, err := db.GetNextFeedToFetch(context.Background())
		if err != nil {
			log.Printf("Failed to get the next feed: %v\n", err)
			break
		}

		// Mark the feed as fetched.
		err = db.MarkFeedFetched(context.Background(), feed.ID)
		if err != nil {
			log.Printf("Failed to mark feed as fetched: %v\n", err)
			break
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
			fmt.Printf("- %s\n", item)
		}

		log.Printf("Feed %v collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
	}
}
