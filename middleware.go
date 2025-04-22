package main

import (
	"context"

	"github.com/firerockets/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, _ := s.db.GetUserByName(context.Background(), s.config.CurrentUserName)
		return handler(s, cmd, user)
	}
}
