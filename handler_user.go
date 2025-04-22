package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/firerockets/gator/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command, _ database.User) error {

	if len(cmd.args) == 0 {
		return fmt.Errorf("expected an argument: username")
	}

	userName := cmd.args[0]

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      userName,
	}

	ctx := context.Background()

	dbUser, err := s.db.GetUserByName(ctx, userName)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if dbUser.Name == userName {
		return fmt.Errorf("user '%s' already exists", userName)
	}

	user, err := s.db.CreateUser(ctx, userParams)

	if err != nil {
		return err
	}

	err = s.config.SetUser(user.Name)

	if err != nil {
		return err
	}

	fmt.Printf("User %s has been created!\n", user.Name)
	fmt.Printf("- ID: %v.\n", user.ID)
	fmt.Printf("- Created at: %v.\n", user.CreatedAt)
	fmt.Printf("- Updated at: %v.\n", user.UpdatedAt)

	return nil
}

func handlerLogin(s *state, cmd command, _ database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("expected an argument: username")
	}

	userName := cmd.args[0]

	ctx := context.Background()
	dbUser, err := s.db.GetUserByName(ctx, userName)

	if err != nil {
		return err
	}

	err = s.config.SetUser(dbUser.Name)

	if err != nil {
		return err
	}

	fmt.Printf("User %s has been set\n", dbUser.Name)

	return nil
}

func handlerUsers(s *state, _ command, _ database.User) error {

	ctx := context.Background()

	users, err := s.db.GetUsers(ctx)

	if err != nil {
		return err
	}

	for _, u := range users {
		currentStr := ""

		if s.config.CurrentUserName == u.Name {
			currentStr = " (current)"
		}

		fmt.Printf("* %s%s\n", u.Name, currentStr)
	}

	return nil
}

func handlerReset(s *state, _ command, _ database.User) error {

	ctx := context.Background()

	err := s.db.DeleteAllUsers(ctx)

	if err != nil {
		return err
	}

	fmt.Println("All users have been deleted.")

	return nil
}
