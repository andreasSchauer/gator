package main

import (
	"context"
	"fmt"
	"time"

	"github.com/andreasSchauer/gator/internal/database"
	"github.com/google/uuid"
)


func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %v <feed_url>", cmd.name)
	}

	feedURL := cmd.args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("couldn't get feed: %w", err)
	}

	deleteFeedFollowParams := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.db.DeleteFeedFollow(context.Background(), deleteFeedFollowParams)
	if err != nil {
		return fmt.Errorf("user %s couldn't unfollow feed %s", user.Name, feed.Name)
	}

	fmt.Printf("User %s successfully unfollowed feed: %s\n", user.Name, feed.Name)
	fmt.Println()
	fmt.Println("=====================================")

	return nil
}



func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %v <feed_url>", cmd.name)
	}

	feedURL := cmd.args[0]
	
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



func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("usage: %v", cmd.name)
	}

	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %w", err)
	}

	if len(feedFollows) == 0 {
		fmt.Printf("user %s doesn't follow any feed yet\n", user.Name)
		return nil
	}

	fmt.Printf("user %s is currently following:\n", user.Name)
	fmt.Println()

	for _, feedFollow := range feedFollows {
		fmt.Printf("* %s\n", feedFollow.FeedName)
	}

	fmt.Println()
	fmt.Println("=====================================")
	return nil
}