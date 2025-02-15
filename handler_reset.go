package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset the database: %v", err)
	}

	fmt.Println("Database reset successfully!")
	return nil
}