package main

import (
	"context"
	"fmt"

	"github.com/BlessE1/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		ctx := context.Background()
		user, err := s.db.GetUser(ctx, s.config.CurrentUserName)
		if err != nil {
			fmt.Println("User could not be fetched")
			return err
		}
		return handler(s, cmd, user)
	}
}
