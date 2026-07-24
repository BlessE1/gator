package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/BlessE1/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		fmt.Println("Time between fetches is required")
		os.Exit(1)
	}

	timeBetweenFetches, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		fmt.Println("Invalid time duration format")
		os.Exit(1)
	}

	fmt.Println("Collecting feeds every", timeBetweenFetches)
	ticker := time.NewTicker(timeBetweenFetches)
	for ; ; <-ticker.C {
		fmt.Println("----| UPDATE TO SCRAPED FEEDS  |----")
		err = scrapeFeeds(s)
		if err != nil {
			fmt.Println("Error scraping feeds:", err)
			os.Exit(1)
		}
	}
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) == 0 {
		fmt.Println("Feed url is required")
		os.Exit(1)
	}

	ctx := context.Background()

	feed, err := s.db.GetFeedByUrl(ctx, cmd.Args[0])
	if err != nil {
		fmt.Println("Feed could not be fetched")
		os.Exit(1)
	}

	followParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollowRow, err := s.db.CreateFeedFollow(ctx, followParams)
	if err != nil {
		fmt.Println("Feed could not be followed")
		os.Exit(1)
	}

	fmt.Println("Feed followed:", feedFollowRow[0].FeedName)
	fmt.Println("Feed followed by user:", feedFollowRow[0].UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()

	feedFollows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		fmt.Println("Feed follows could not be fetched")
		os.Exit(1)
	}

	fmt.Println("Feeds followed by user:", user.Name)
	for _, feedFollow := range feedFollows {
		fmt.Println("*", feedFollow.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) == 0 {
		fmt.Println("Feed url is required")
		os.Exit(1)
	}

	ctx := context.Background()
	feed, err := s.db.GetFeedByUrl(ctx, cmd.Args[0])
	if err != nil {
		fmt.Println("Feed could not be fetched")
		os.Exit(1)
	}

	removeFeedFollowParams := database.RemoveFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.db.RemoveFeedFollow(ctx, removeFeedFollowParams)
	if err != nil {
		fmt.Println("Feed could not be unfollowed")
		os.Exit(1)
	}

	fmt.Println("Feed unfollowed:", feed.Name)
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 2 {
		fmt.Println("Name and Feed URL are required")
		os.Exit(1)
	}

	ctx := context.Background()

	addFeedParams := database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	}

	insertedFeed, err := s.db.AddFeed(ctx, addFeedParams)
	if err != nil {
		fmt.Println("Feed could not be added")
		os.Exit(1)
	}

	cmd.Args = []string{insertedFeed.Url}
	err = handlerFollow(s, cmd, user)
	if err != nil {
		fmt.Println("Feed could not be followed")
		os.Exit(1)
	}

	fmt.Println("Feed added:", insertedFeed.Name, "with URL:", insertedFeed.Url)
	fmt.Println("Feed ID:", insertedFeed.ID)
	fmt.Println("Created at:", insertedFeed.CreatedAt)
	fmt.Println("Updated at:", insertedFeed.UpdatedAt)
	fmt.Println("User ID:", insertedFeed.UserID)
	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	ctx := context.Background()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		fmt.Println("Feeds could not be fetched")
		os.Exit(1)
	}

	for _, feed := range feeds {
		fmt.Println("Feed ID:", feed.ID)
		fmt.Println("Feed Name:", feed.Name)
		fmt.Println("Feed URL:", feed.Url)
		user, err := s.db.GetUserById(ctx, feed.UserID)
		if err != nil {
			fmt.Println("User could not be fetched")
			os.Exit(1)
		}

		fmt.Println("Feed Creator Name:", user.Name)
		fmt.Println("-------------------------")
	}

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.Args) == 0 {
		fmt.Println("Limit set to 2 by default")
		cmd.Args = []string{"2"}
	}

	limit, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		fmt.Println("Invalid limit format")
		os.Exit(1)
	}

	ctx := context.Background()
	getPostsForUserParams := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}

	posts, err := s.db.GetPostsForUser(ctx, getPostsForUserParams)
	if err != nil {
		fmt.Println("Posts could not be fetched")
		os.Exit(1)
	}
	fmt.Println("Posts for user:", user.Name)
	fmt.Println("-------------------------")

	for _, post := range posts {
		fmt.Println("Post Title:", post.Title)
		fmt.Println("Post Description:", post.Description.String)
		fmt.Println("Post Created At:", post.CreatedAt)
		fmt.Println("Post Updated At:", post.UpdatedAt)
		fmt.Println("-------------------------")
	}

	return nil
}
