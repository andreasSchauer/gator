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
		return fmt.Errorf("couldn't log posts: %w", err)
	}

	fmt.Println("Listing Posts:")
	fmt.Println()
	
	for _, post := range posts {
		fmt.Printf("* Title: %s\n", post.Title)
		fmt.Printf("* Published At: %v\n", post.PublishedAt.Time)
		fmt.Printf("* URL: %s\n", post.Url)
		fmt.Printf("* Description: %v\n", post.Description.String)
		fmt.Println()
	}

	return nil
}