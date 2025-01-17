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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("expected a single argument, username")
	}

	userName := cmd.args[0]
	userExists, err := s.db.UserExists(context.Background(), userName)
	if err != nil {
		log.Fatalf("ran into a issue checking if user exists: %v", err)
	}

	if !userExists {
		log.Fatalf("user %v doesnt exist", userName)
	}

	s.cfg.CurrentUserName = userName
	s.cfg.SetUser(userName)

	return nil
}

func handlerRegitser(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("expected a username to be provided")
	}

	userName := cmd.args[0]
	userExists, err := s.db.UserExists(context.Background(), userName)
	if err != nil {
		log.Fatalf("ran into a issue checking if user exists: %v", err)
	}

	if userExists {
		log.Fatalf("User with name %v already registered", userName)
	}

	newId := uuid.New()
	now := time.Now()
	queryData := database.CreateUserParams{
		ID:        newId,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      userName,
	}

	_, err = s.db.CreateUser(context.Background(), queryData)
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

func handlerUsers(s *state, cmd command) error {
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

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return errors.New("expected a <feed_name> and a <feed_url> to be provided")
	}

	currentUserName := s.cfg.CurrentUserName
	if currentUserName == "" {
		return errors.New("Must be signed in before you can add a feed")
	}

	ctx := context.Background()
	userExists, err := s.db.UserExists(ctx, currentUserName)
	if err != nil {
		return err
	}

	if !userExists {
		errMsg := fmt.Sprintf("Username of %v does not exist", currentUserName)
		return errors.New(errMsg)
	}

	feedName := cmd.args[0]
	rssUrl := cmd.args[1]

	_, err = service.FetchFeed(context.Background(), rssUrl)
	if err != nil {
		return err
	}

	user, err := s.db.GetUser(ctx, currentUserName)
	if err != nil {
		return err
	}

	queryData := database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      feedName,
		Url:       rssUrl,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdFeed, err := s.db.CreateFeed(ctx, queryData)
	if err != nil {
		uniqeConstratintError := "23505"
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == pq.ErrorCode(uniqeConstratintError) {
			return errors.New("feed URL already exists")
		}
		return err
	}

	fmt.Println(createdFeed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("Feed Name: %v | Feed Url: %v | Added By: %v\n", feed.Name, feed.Url, feed.Username)
	}
	return nil
}
