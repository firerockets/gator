package main

import (
	"context"
	"fmt"
	"html"
	"strconv"
	"time"

	"github.com/firerockets/gator/internal/database"
	"github.com/firerockets/gator/internal/rss"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command, _ database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("expected argument time duration")
	}

	strDuration := cmd.args[0]

	duration, err := time.ParseDuration(strDuration)

	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %s.\n", strDuration)

	ticker := time.NewTicker(duration)

	for ; ; <-ticker.C {
		err = scrapeFeeds(s.db)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("missing arguments - required name and url")
	}

	feedName := cmd.args[0]
	feedUrl := cmd.args[1]

	if user == (database.User{}) {
		return fmt.Errorf("no user found - register a user first")
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    user.ID,
	})

	if err != nil {
		return err
	}

	s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	fmt.Printf("RSS Feed '%s' was added.\n", feed.Name)
	fmt.Printf("- URL: %s\n", feed.Url)
	fmt.Printf("- User ID: %v\n", feed.UserID)

	return nil
}

func handlerFeeds(s *state, _ command, _ database.User) error {

	feeds, err := s.db.GetFeeds(context.Background())

	if err != nil {
		return err
	}

	for _, f := range feeds {
		user, err := s.db.GetUserById(context.Background(), f.UserID)

		if err != nil {
			return err
		}

		fmt.Printf("Feed Name: %s\n", f.Name)
		fmt.Printf("Feed URL: %s\n", f.Url)
		fmt.Printf("Created by: %s\n", user.Name)
		fmt.Println("")
	}

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2

	var err error

	if len(cmd.args) == 1 {
		limit, err = strconv.Atoi(cmd.args[0])
	}

	if err != nil {
		return err
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})

	if err != nil {
		return err
	}

	fmt.Println("################")
	fmt.Printf("Printing posts for %s.\n", user.Name)

	for _, p := range posts {
		fmt.Println(p.Title)
		fmt.Println(p.Description)
	}

	return nil
}

func scrapeFeeds(db *database.Queries) error {
	feed, err := db.GetNextFeedToFetch(context.Background())

	if err != nil {
		return err
	}

	err = db.MarkFeedFetched(context.Background(), feed.ID)

	if err != nil {
		return err
	}

	rss, err := rss.FetchFeed(context.Background(), feed.Url)

	if err != nil {
		return err
	}

	fmt.Printf("Fetching rss feed for %s.\n", feed.Name)

	for i, item := range rss.Channel.Item {
		rss.Channel.Item[i].Title = html.UnescapeString(item.Title)
		rss.Channel.Item[i].Description = html.UnescapeString(item.Description)

		db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       rss.Channel.Item[i].Title,
			Url:         rss.Channel.Item[i].Link,
			Description: rss.Channel.Item[i].Description,
			FeedID:      feed.ID,
		})
	}

	return nil
}
