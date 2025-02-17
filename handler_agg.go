package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/andreasSchauer/gator/internal/database"
)



func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %v <time_between_requests>", cmd.name)
	}

	time_between_requests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't parse duration value: %w", err)
	}
	ticker := time.NewTicker(time_between_requests)

	fmt.Printf("Collecting feeds every %s\n", time_between_requests)
	fmt.Println()

	for ; ; <-ticker.C {
		fmt.Println("next feed incoming")
		scrapeFeeds(s)
	}
}


func scrapeFeeds(s *state) {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())

	if err != nil {
		fmt.Println("couldn't get next feeds to fetch", err)
		return
	}

	fmt.Println("Found a feed to fetch!")

	markFeedFetchedParams := database.MarkFeedFetchedParams{
		ID: nextFeed.ID,
		LastFetchedAt: sql.NullTime{
			Time: time.Now().Local(),
			Valid: true,
		},
	}

	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		fmt.Printf("couldn't collect feed %s: %v", nextFeed.Name, err)
		return
	}

	err = s.db.MarkFeedFetched(context.Background(), markFeedFetchedParams)
	if err != nil {
		fmt.Printf("couldn't mark feed %s as fetched: %v", nextFeed.Name, err)
		return
	}

	for _, item := range feed.Channel.Item {
		fmt.Printf("Found post: %s\n", item.Title)
	}

	fmt.Printf("Feed %s collected, %v posts found\n", nextFeed.Name, len(feed.Channel.Item))
}