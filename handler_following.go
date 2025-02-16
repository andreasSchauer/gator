package main

import (
	"context"
	"fmt"
)

func handlerFollowing(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("usage: %v", cmd.name)
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
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



/*
type GetFeedFollowsForUserRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
	FeedName  string
	UserName  string
}

*/