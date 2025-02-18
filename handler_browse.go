package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/andreasSchauer/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 1 {
		return fmt.Errorf("usage: %v <postsLimit>", cmd.name)
	}
	
	var postsLimit int32
	if len(cmd.args) < 1 {
		postsLimit = 2
	} else {
		num, err := strconv.ParseInt(cmd.args[0], 10, 32)
		if err != nil {
			return err
		}
		postsLimit = int32(num)
	}

	if postsLimit <= 0 {
		return fmt.Errorf("postsLimit must be a positive number")
	}

	if postsLimit > 50 {
		return fmt.Errorf("postsLimit cannot exceed 50 posts")
	}

	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:	postsLimit,
	}

	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("couldn't get posts for user %s: %w", user.Name, err)
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	fmt.Println()
	
	for _, post := range posts {
		fmt.Printf("* Published At: %v from %s\n", post.PublishedAt.Time.Format("Mon 2006-01-02 15:04:05"), post.FeedName)
		fmt.Printf("* --- %s ---\n", post.Title)
		fmt.Printf("* Link: %s\n", post.Url)
		fmt.Printf("* 	%v\n", post.Description.String)
		fmt.Printf("\n\n")
	}

	return nil
}