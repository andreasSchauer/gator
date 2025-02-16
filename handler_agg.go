package main

import (
	"context"
	"fmt"
)



func handlerAgg(_ *state, _ command) error {
	feedURL := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}

	printRSSFeed(feed)
	return nil
}


func printRSSFeed(feed *RSSFeed) {
	fmt.Println("==== RSS Feed ====")
	fmt.Printf("Channel: %s\n", feed.Channel.Title)
	fmt.Printf("Description: %s\n", feed.Channel.Description)
	fmt.Printf("Link: %s\n", feed.Channel.Link)
	fmt.Println("Posts:")

	for i := range feed.Channel.Item {
		item := &feed.Channel.Item[i]
		fmt.Println()
		fmt.Printf("Title: %s\n", item.Title)
		fmt.Printf("Link: %s\n", item.Link)
		fmt.Printf("Published: %s\n", item.PubDate)
		fmt.Printf("Description: %s\n", item.Description)
	}
}