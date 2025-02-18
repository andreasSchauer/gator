package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/andreasSchauer/gator/internal/database"
	"github.com/google/uuid"
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
		err := saveRssItemAsPost(s, item, nextFeed.ID)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Printf("Feed %s collected, %v posts found\n", nextFeed.Name, len(feed.Channel.Item))
}


func saveRssItemAsPost(s *state, item RSSItem, feedID uuid.UUID) error {
	publishedAt, err := evaluatePubDate(item)
	if err != nil {
		return err
	}
	
	postParams := database.CreatePostParams{
		ID:				uuid.New(),
		CreatedAt: 		time.Now().Local(),
		UpdatedAt: 		time.Now().Local(),
		Title:			item.Title,
		Url: 			item.Link,
		Description:	sql.NullString{
			String: 	item.Description,
			Valid:		item.Description != "",
		},
		PublishedAt: 	publishedAt,
		FeedID: 		feedID,
	}
	
	_ , err = s.db.CreatePost(context.Background(), postParams)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil
		} 
		return fmt.Errorf("couldn't create post with title %s: %v", item.Title, err)
	}

	return nil
}


func evaluatePubDate(item RSSItem) (sql.NullTime, error) {
	var postPublishDate time.Time
	var publishedAtValid bool
	
	if item.PubDate != "" {
		parsedTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return sql.NullTime{}, fmt.Errorf("couldn't parse publish date of post %s: %v", item.Title, err)
		}
		postPublishDate = parsedTime
		publishedAtValid = true
	}

	publishedAt := sql.NullTime{
		Time: 		postPublishDate,
		Valid:		publishedAtValid,
	}

	return publishedAt, nil
}