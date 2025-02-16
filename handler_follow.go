package main

import (
	"context"
	"fmt"
	"time"

	"github.com/andreasSchauer/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %v <url>", cmd.name)
	}

	feedURL := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("couldn't find a feed at: %s", feedURL)
	}

	err = followFeed(s, user, feed)
	if err != nil {
		return err
	}

	return nil
}


func followFeed(s *state, user database.User, feed database.Feed) error {
	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().Local(),
		UpdatedAt: time.Now().Local(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return fmt.Errorf("couldn't follow feed: %w", err)
	}
	fmt.Printf("User %s successfully followed feed: %s\n", user.Name, feed.Name)
	fmt.Println()
	fmt.Println("=====================================")

	return nil
}