package main

import (
	"context"
	"fmt"
	"time"

	"github.com/andreasSchauer/gator/internal/database"
	"github.com/google/uuid"
)


func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("usage: %v", cmd.name)
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	fmt.Printf("Found %d feeds:\n", len(feeds))
	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't get user: %w", err)
		}
		printFeed(feed, user)
		fmt.Println("=====================================")
	}

	return nil
}



func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: %v <feedName> <feedURL>", cmd.name)
	}
	
	feedName := cmd.args[0]
	feedURL := cmd.args[1]

	feedParams := database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().Local(),
		UpdatedAt: time.Now().Local(),
		Name: feedName,
		Url: feedURL,
		UserID: user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	fmt.Println("Feed created successfully!")
	printFeed(feed, user)
	fmt.Println()

	err = followFeed(s, user, feed)
	if err != nil {
		return err
	}

	return nil
}


func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:       	    %s\n", feed.ID)
	fmt.Printf("* Created:  	    %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:   	    %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:      	    %s\n", feed.Name)
	fmt.Printf("* URL:       	    %s\n", feed.Url)
	fmt.Printf("* User:        		%s\n", user.Name)
	fmt.Printf("* LastFetchedAt:	%v\n", feed.LastFetchedAt.Time)
}