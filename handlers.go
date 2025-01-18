package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/acemouty/gator/internal/database"
	"github.com/acemouty/gator/internal/service"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"log"
	"time"
)

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
