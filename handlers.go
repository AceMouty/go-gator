package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/acemouty/gator/internal/database"
	"github.com/acemouty/gator/internal/service"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("useage: %v <time_between_requests> ex time formats 1s, 1m or 1h", cmd.name)
	}

	time_between_requests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return errors.New("Unable to create interval and aggregate at this time")
	}

	ticker := time.NewTicker(time_between_requests)
	fmt.Printf("Collecting feeds every %v\n\n", cmd.args[0])

	/*
					  The agg command is a never-ending loop that fetches feeds and prints posts to the console.
				    The intended use case is to leave the agg command running in the background while
				    interacting with the program in another terminal.

		        Loop steps at a interval that is provided by the user when running that agg command
	*/
	scrapeRound := 0
	for ; ; <-ticker.C {
		scrapeFeeds(s.db, &scrapeRound)
	}
}

func handlerLogin(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return errors.New("expected a single argument, username")
	}

	userName := cmd.args[0]
	s.cfg.SetUser(userName)

	return nil
}

func handlerRegitser(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("expected a username to be provided")
	}

	userName := cmd.args[0]
	newId := uuid.New()
	now := time.Now()
	queryData := database.CreateUserParams{
		ID:        newId,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      userName,
	}

	_, err := s.db.CreateUser(context.Background(), queryData)
	if err != nil {
		log.Fatalf("unable to register user %v: reason; %v", userName, err)
	}

	log.Println("registered user successfully")
	s.cfg.SetUser(userName)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())

	if err != nil {
		log.Fatalf("ran into a issue deleting users: %v", err)
	}

	return nil
}

func handlerUsers(s *state, cmd command, user database.User) error {
	users, err := s.db.GetUsers(context.Background())

	if err != nil {
		log.Fatalf("ran into a issue getting users: %v", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("%v (current)\n", user.Name)
			continue
		}

		fmt.Println(user.Name)
	}

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("expected a <feed_name> and a <feed_url> to be provided")
	}

	currentUserName := s.cfg.CurrentUserName
	if currentUserName == "" {
		return errors.New("Must be signed in before you can add a feed")
	}

	ctx := context.Background()
	feedName := cmd.args[0]
	rssUrl := cmd.args[1]

	_, err := service.FetchFeed(context.Background(), rssUrl)
	if err != nil {
		return err
	}

	createFeedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      feedName,
		Url:       rssUrl,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdFeed, err := s.db.CreateFeed(ctx, createFeedParams)
	if err != nil {
		uniqeConstratintError := "23505"
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == pq.ErrorCode(uniqeConstratintError) {
			return errors.New("feed URL already exists")
		}
		return err
	}

	createFollowFeedParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    createdFeed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), createFollowFeedParams)
	if err != nil {
		return errors.New("Unable to follow feed at this time")
	}

	fmt.Println(createdFeed)
	fmt.Printf("\n%v now following %v\n", user.Name, createdFeed.Url)
	return nil
}

func handlerFeeds(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("Feed Name: %v | Feed Url: %v | Added By: %v\n", feed.Name, feed.Url, feed.Username)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return errors.New("Expected to be provided a url")
	}

	feedUrl := cmd.args[0]
	feedExists, err := s.db.FeedExists(context.Background(), feedUrl)
	if err != nil {
		return errors.New("Unable to validate feed url at this time")
	}

	if !feedExists {
		errMsg := fmt.Sprintf("Feed %v doesnt exist", feedUrl)
		return errors.New(errMsg)
	}

	feed, err := s.db.GetFeed(context.Background(), feedUrl)
	if err != nil {
		return errors.New("Unable to grab feed at this time")
	}

	queryData := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), queryData)
	if err != nil {
		return errors.New("Unable to follow feed at this time")
	}

	fmt.Printf("User %v now following %v", user.Name, feed.Url)
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return errors.New("Expected one argument <feed_url>")
	}

	feedUrl := cmd.args[0]
	deleteFeedFollowParams := database.DeleteFeedFollowParams{UserID: user.ID, Url: feedUrl}
	err := s.db.DeleteFeedFollow(context.Background(), deleteFeedFollowParams)
	if err != nil {
		return errors.New("Unable to delete feed at this time")
	}
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	followingFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		fmt.Println(err)
		return errors.New("Unable to get users following feeds")
	}

	for _, feed := range followingFeeds {
		fmt.Printf("Feed Name: %v\n", feed.Feedname)
	}

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 1

	if len(cmd.args) == 0 {
		return fmt.Errorf("useage: %v <limit>\nwhere limit between 1 and 100", cmd.name)
	}

	if providedLimit, err := strconv.Atoi(cmd.args[0]); err == nil {
		limit = providedLimit
	} else {
		return fmt.Errorf("invalid limit: %v", err)
	}

	getPostForUserParams := database.GetPostsForUserParams{UserID: user.ID, Limit: int32(limit)}
	posts, err := s.db.GetPostsForUser(context.Background(), getPostForUserParams)
	if err != nil {
		return fmt.Errorf("couldn't get posts for user: %v", err)
	}

	fmt.Printf("Found %v posts for user %v:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}
	return nil
}
