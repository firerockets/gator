package main

import (
	"context"
	"fmt"
	"time"

	"github.com/firerockets/gator/internal/database"
	"github.com/google/uuid"
)

func handleFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("missing arguments - required url")
	}

	feedUrl := cmd.args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)

	if err != nil {
		return err
	}

	if user == (database.User{}) {
		return fmt.Errorf("no user found - register a user first")
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		return err
	}

	fmt.Printf("The user '%s' is now following the RSS feed '%s'.\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handleFollowing(s *state, _ command, user database.User) error {

	if user == (database.User{}) {
		return fmt.Errorf("no user found - register a user first")
	}

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)

	if err != nil {
		return err
	}

	if len(feeds) == 0 {
		fmt.Printf("The user '%s' is not following any feed yet.\n", user.Name)
		return nil
	}

	fmt.Printf("The user '%s' is following the feeds:\n", user.Name)

	for _, f := range feeds {
		fmt.Printf("- %s\n", f.FeedName)
	}

	return nil
}

func handleUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("missing arguments - required url")
	}

	feedUrl := cmd.args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)

	if err != nil {
		return err
	}

	if user == (database.User{}) {
		return fmt.Errorf("no user found - register a user first")
	}

	err = s.db.DeleteFeedFollowForUserIdAndFeedId(context.Background(), database.DeleteFeedFollowForUserIdAndFeedIdParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return err
	}

	fmt.Println("Successfully unfollowed.")

	return nil
}
